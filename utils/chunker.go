package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"sync"

	"github.com/ilemus/agents/llm"
	"github.com/ilemus/agents/llm/ollama"

	"gonum.org/v1/gonum/floats"
)

var (
	BREAKPOINT_THRESHOLD = 0.85 // Cosine similarity BELOW this triggers a split
	MIN_SENTENCES        = 2    // Prevents trivially small chunks
	MAX_CHARS            = 1500 // Safety valve: forces a split if a chunk gets too long
	WINDOW_SIZE          = 3    // Number of sentences to average for smoothing
)

// Chunk represents a semantically cohesive piece of text.
// It retains metadata (indices and raw sentences) for downstream RAG citation.
type Chunk struct {
	Text       string
	Sentences  []string
	StartIdx   int
	EndIdx     int
	CharLength int
}

// SemanticChunker encapsulates the configuration and state for semantic text splitting.
type SemanticChunker struct {
	ModelName           string
	BreakpointThreshold float64 // Cosine similarity BELOW this triggers a split
	MinSentences        int     // Prevents trivially small chunks
	MaxChars            int     // Safety valve: forces a split if a chunk gets too long
	WindowSize          int     // Number of sentences to average for smoothing
	Client              llm.Client
}

// NewSemanticChunker initializes the chunker with sensible defaults.
func NewSemanticChunker(modelName string) *SemanticChunker {
	client := ollama.NewClient("")

	return &SemanticChunker{
		ModelName:           modelName,
		BreakpointThreshold: BREAKPOINT_THRESHOLD,
		MinSentences:        MIN_SENTENCES,
		MaxChars:            MAX_CHARS,
		WindowSize:          WINDOW_SIZE,
		Client:              client,
	}
}

// splitSentences divides raw text into sentences.
// Note: Go's RE2 regex engine does not support lookarounds (e.g., (?<=...)).
// We use a replacement strategy to preserve punctuation while splitting cleanly.
func splitSentences(text string) []string {
	// 1. Find sentence-ending punctuation followed by whitespace.
	// 2. Replace it with the punctuation + a newline delimiter.
	re := regexp.MustCompile(`([.?!])\s+`)
	normalized := re.ReplaceAllString(text, "$1\n")

	// 3. Split by the newline delimiter to get clean sentences with punctuation intact.
	rawParts := strings.Split(normalized, "\n")

	var sentences []string
	for _, s := range rawParts {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			sentences = append(sentences, trimmed)
		}
	}
	return sentences
}

// cosineSimilarity calculates the cosine similarity between two float64 vectors.
// It leverages gonum/floats for highly optimized, SIMD-accelerated linear algebra.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	dotProduct := floats.Dot(a, b)
	normA := floats.Norm(a, 2) // L2 norm
	normB := floats.Norm(b, 2)

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (normA * normB)
}

// meanPool calculates the element-wise average of a slice of vectors.
// This is used to smooth out noise by averaging the embeddings of the sliding window.
func meanPool(vectors [][]float64) []float64 {
	if len(vectors) == 0 {
		return nil
	}

	dim := len(vectors[0])
	mean := make([]float64, dim)

	// Sum all vectors element-wise
	for _, v := range vectors {
		floats.Add(mean, v)
	}

	// Divide by the count to get the average
	count := float64(len(vectors))
	floats.Scale(1.0/count, mean)

	return mean
}

// ChunkText is the main pipeline. It fetches embeddings, calculates similarities
// concurrently, and builds chunks based on semantic thresholds and length limits.
func (sc *SemanticChunker) ChunkText(ctx context.Context, text string) ([]Chunk, error) {
	sentences := splitSentences(text)

	// Edge case: Text is too short to chunk meaningfully
	if len(sentences) <= 1 {
		return []Chunk{{
			Text: text, Sentences: sentences,
			StartIdx: 0, EndIdx: 0, CharLength: len(text),
		}}, nil
	}

	// ---------------------------------------------------------
	// 1. PERFORMANCE: Batch Embedding (Fixing Gemini's sequential flaw)
	// ---------------------------------------------------------
	// We pass the ENTIRE slice of sentences to the `Input` field.
	// This results in a SINGLE API request, drastically reducing network overhead.
	req := llm.EmbedRequest{
		Model: sc.ModelName,
		Input: sentences,
	}

	resp, err := sc.Client.Embed(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get embeddings: %w", err)
	}

	// Ollama returns [][]float32. We must convert to [][]float64 for gonum/floats.
	embeddings := make([][]float64, len(resp.Embeddings))
	for i, emb32 := range resp.Embeddings {
		emb64 := make([]float64, len(emb32))
		for j, val := range emb32 {
			emb64[j] = float64(val)
		}
		embeddings[i] = emb64
	}

	// ---------------------------------------------------------
	// 2. ACCURACY & CONCURRENCY: Sliding Window Similarities
	// ---------------------------------------------------------
	numSimilarities := len(sentences) - 1
	similarities := make([]float64, numSimilarities)
	var wg sync.WaitGroup

	// We spawn a goroutine for EACH similarity calculation.
	// Because this is pure CPU-bound math (no I/O), this perfectly utilizes all available CPU cores.
	for i := range numSimilarities {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Calculate the boundaries for the left and right sliding windows
			leftStart := int(math.Max(0, float64(idx-sc.WindowSize+1)))
			rightEnd := int(math.Min(float64(len(sentences)), float64(idx+sc.WindowSize+1)))

			// Extract the vectors for the current window
			leftVecs := embeddings[leftStart : idx+1]
			rightVecs := embeddings[idx+1 : rightEnd]

			// Mean-pool the vectors to smooth out anomalous single sentences
			leftMean := meanPool(leftVecs)
			rightMean := meanPool(rightVecs)

			// Compute and store the similarity score safely (each goroutine writes to a unique index)
			similarities[idx] = cosineSimilarity(leftMean, rightMean)
		}(i)
	}

	// Wait for all CPU-bound goroutines to finish before proceeding
	wg.Wait()

	// ---------------------------------------------------------
	// 3. LOGIC: Reconstruct Chunks (Sequential State Machine)
	// ---------------------------------------------------------
	var chunks []Chunk
	currentSentences := []string{sentences[0]}
	prevIdx := 0

	for i, sim := range similarities {
		nextSentence := sentences[i+1]

		// Simulate adding the next sentence to the current chunk
		potentialSentences := append([]string{}, currentSentences...)
		potentialSentences = append(potentialSentences, nextSentence)
		potentialText := strings.Join(potentialSentences, " ")

		// SAFETY VALVE (from Qwen): Force split if we exceed the character limit
		forceSplitByLength := len(potentialText) > sc.MaxChars

		// SEMANTIC SPLIT: Similarity dropped below the cohesion threshold
		semanticSplit := sim < sc.BreakpointThreshold

		if forceSplitByLength || semanticSplit {
			// Only finalize the chunk if it meets the minimum sentence requirement
			if len(currentSentences) >= sc.MinSentences {
				chunkText := strings.Join(currentSentences, " ")
				chunks = append(chunks, Chunk{
					Text:       chunkText,
					Sentences:  currentSentences,
					StartIdx:   prevIdx,
					EndIdx:     i,
					CharLength: len(chunkText),
				})

				// Reset state for the next chunk
				currentSentences = []string{nextSentence}
				prevIdx = i + 1
			} else {
				// If it's too small, absorb the next sentence anyway to satisfy MinSentences
				currentSentences = potentialSentences
			}
		} else {
			// Similarity is high and length is fine; keep building the current chunk
			currentSentences = potentialSentences
		}
	}

	// ---------------------------------------------------------
	// 4. FINALIZE: Flush remaining sentences
	// ---------------------------------------------------------
	if len(currentSentences) > 0 {
		chunkText := strings.Join(currentSentences, " ")
		chunks = append(chunks, Chunk{
			Text:       chunkText,
			Sentences:  currentSentences,
			StartIdx:   prevIdx,
			EndIdx:     len(sentences) - 1,
			CharLength: len(chunkText),
		})
	}

	return chunks, nil
}

func Main() {
	sampleText := "PostgreSQL is a powerful, open-source object-relational database system. " +
		"It uses and extends the SQL language combined with many features that safely store and scale data. " +
		"Its performance and reliability have earned it a strong reputation among engineers. " +
		"In contrast, car camping offers a unique way to explore nature without sacrificing comfort. " +
		"When packing for a weekend trip, bringing a lightweight sleeping pad and a reliable hatchet is critical. " +
		"Cooking fresh salmon over an open fire or a small camp stove makes the outdoor experience memorable."

	chunker := NewSemanticChunker("embeddinggemma:latest")

	// Lowering MaxChars to 250 to explicitly demonstrate the Qwen safety valve in action
	chunker.MaxChars = 250

	fmt.Println("Processing semantic chunking concurrently in Go...")

	ctx := context.Background()
	chunks, err := chunker.ChunkText(ctx, sampleText)
	if err != nil {
		log.Fatalf("Error chunking text: %v", err)
	}

	for idx, chunk := range chunks {
		fmt.Printf("--- Chunk %d (Chars: %d, Sentences: %d, Idx: %d-%d) ---\n",
			idx+1, chunk.CharLength, len(chunk.Sentences), chunk.StartIdx, chunk.EndIdx)
		fmt.Println(chunk.Text)
		fmt.Println()
	}
}
