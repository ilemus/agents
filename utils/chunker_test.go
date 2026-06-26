package utils

import (
	"context"
	"math/rand"
	"testing"

	"github.com/ilemus/agents/llm"
)

type mockClient struct{}

func (m *mockClient) Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	return llm.ChatResponse{}, nil
}

func (m *mockClient) ChatStream(ctx context.Context, req llm.ChatRequest, onChunk func(string) error) (llm.ChatResponse, error) {
	return llm.ChatResponse{}, nil
}

func (m *mockClient) Embed(ctx context.Context, req llm.EmbedRequest) (llm.EmbedResponse, error) {
	embeddings := make([][]float32, len(req.Input))
	for i := range req.Input {
		emb := make([]float32, 768)
		for j := 0; j < 768; j++ {
			// Generate random float between -1 and 1 to ensure low cosine similarity
			emb[j] = rand.Float32()*2 - 1.0
		}
		embeddings[i] = emb
	}
	return llm.EmbedResponse{Embeddings: embeddings}, nil
}

func TestSemanticChunker(t *testing.T) {
	client := &mockClient{}

	chunker := &SemanticChunker{
		ModelName:           "mock_model",
		BreakpointThreshold: 0.85,
		MinSentences:        1, // As requested
		MaxChars:            1500,
		WindowSize:          3,
		Client:              client,
	}

	sampleText := "This is the first sentence. Here is the second sentence. And the third one! What about a fourth? This is the fifth sentence."

	ctx := context.Background()
	chunks, err := chunker.ChunkText(ctx, sampleText)
	if err != nil {
		t.Fatalf("ChunkText failed: %v", err)
	}

	if len(chunks) <= 1 {
		t.Errorf("Expected multiple chunks due to random embeddings, got %d", len(chunks))
	}

	for i, chunk := range chunks {
		t.Logf("Chunk %d: %q (Chars: %d, Sentences: %d)", i, chunk.Text, chunk.CharLength, len(chunk.Sentences))
	}
}
