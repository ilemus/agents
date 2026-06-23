package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agents/llm"
	"github.com/ollama/ollama/api"
)

func TestChat(t *testing.T) {
	// Start a local HTTP server to mock the Ollama response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("expected path /api/chat, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Decode the request body using official API ChatRequest
		var req api.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if req.Model != "test-model" {
			t.Errorf("expected model 'test-model', got '%s'", req.Model)
		}

		// Send mock response using official API ChatResponse
		resp := api.ChatResponse{
			Model: "test-model",
			Message: api.Message{
				Role:    "assistant",
				Content: "Hello, this is a mock response!",
			},
			Done: true,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req := llm.ChatRequest{
		Model: "test-model",
		Messages: []llm.Message{
			{Role: "user", Content: "Hello!"},
		},
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if resp.Message.Content != "Hello, this is a mock response!" {
		t.Errorf("expected content 'Hello, this is a mock response!', got '%s'", resp.Message.Content)
	}
	if resp.Message.Role != "assistant" {
		t.Errorf("expected role 'assistant', got '%s'", resp.Message.Role)
	}
}

func TestChatStream(t *testing.T) {
	// Start a local HTTP server to mock the Ollama streaming response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Write chunks using official API ChatResponse
		chunks := []api.ChatResponse{
			{
				Model: "test-model",
				Message: api.Message{
					Role:    "assistant",
					Content: "Hello ",
				},
				Done: false,
			},
			{
				Model: "test-model",
				Message: api.Message{
					Role:    "assistant",
					Content: "world!",
				},
				Done: true,
			},
		}

		encoder := json.NewEncoder(w)
		for _, chunk := range chunks {
			if err := encoder.Encode(chunk); err != nil {
				t.Errorf("failed to encode chunk: %v", err)
			}
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req := llm.ChatRequest{
		Model: "test-model",
		Messages: []llm.Message{
			{Role: "user", Content: "Hello!"},
		},
	}

	var streamOutput []string
	onChunk := func(chunk string) error {
		streamOutput = append(streamOutput, chunk)
		return nil
	}

	resp, err := client.ChatStream(context.Background(), req, onChunk)
	if err != nil {
		t.Fatalf("ChatStream failed: %v", err)
	}

	if resp.Message.Content != "Hello world!" {
		t.Errorf("expected content 'Hello world!', got '%s'", resp.Message.Content)
	}

	if len(streamOutput) != 2 {
		t.Errorf("expected 2 stream chunks, got %d", len(streamOutput))
	} else {
		if streamOutput[0] != "Hello " {
			t.Errorf("expected first chunk 'Hello ', got '%s'", streamOutput[0])
		}
		if streamOutput[1] != "world!" {
			t.Errorf("expected second chunk 'world!', got '%s'", streamOutput[1])
		}
	}
}

func TestEmbed(t *testing.T) {
	// Start a local HTTP server to mock the Ollama response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embed" {
			t.Errorf("expected path /api/embed, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Decode the request body using official API EmbedRequest
		var req api.EmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if req.Model != "test-embed-model" {
			t.Errorf("expected model 'test-embed-model', got '%s'", req.Model)
		}

		// Verify input matches expected strings
		inputs, ok := req.Input.([]interface{})
		if !ok || len(inputs) != 2 || inputs[0] != "hello" || inputs[1] != "world" {
			t.Errorf("expected inputs ['hello', 'world'], got %v", req.Input)
		}

		// Send mock response using official API EmbedResponse
		resp := api.EmbedResponse{
			Model: "test-embed-model",
			Embeddings: [][]float32{
				{0.1, 0.2, 0.3},
				{0.4, 0.5, 0.6},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req := llm.EmbedRequest{
		Model: "test-embed-model",
		Input: []string{"hello", "world"},
	}

	resp, err := client.Embed(context.Background(), req)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(resp.Embeddings) != 2 {
		t.Fatalf("expected 2 embeddings, got %d", len(resp.Embeddings))
	}

	if resp.Embeddings[0][0] != 0.1 || resp.Embeddings[0][1] != 0.2 || resp.Embeddings[0][2] != 0.3 {
		t.Errorf("expected first embedding [0.1, 0.2, 0.3], got %v", resp.Embeddings[0])
	}
	if resp.Embeddings[1][0] != 0.4 || resp.Embeddings[1][1] != 0.5 || resp.Embeddings[1][2] != 0.6 {
		t.Errorf("expected second embedding [0.4, 0.5, 0.6], got %v", resp.Embeddings[1])
	}
}
