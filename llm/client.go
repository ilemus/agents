package llm

import "context"

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`    // "system", "user", or "assistant"
	Content string `json:"content"` // The text content of the message
}

// ChatRequest represents the standard request payload for a chat completion.
type ChatRequest struct {
	Model    string                 `json:"model"`             // The name of the LLM model to use
	Messages []Message              `json:"messages"`          // The conversational context
	Options  map[string]interface{} `json:"options,omitempty"` // Provider-specific optional parameters
}

// ChatResponse represents the standard response payload from a chat completion.
type ChatResponse struct {
	Message Message `json:"message"` // The generated message from the assistant
}

// EmbedRequest represents the payload for generating text embeddings.
type EmbedRequest struct {
	Model   string                 `json:"model"`             // The name of the LLM model to use
	Input   []string               `json:"input"`             // The input texts to embed
	Options map[string]interface{} `json:"options,omitempty"` // Provider-specific optional parameters
}

// EmbedResponse represents the generated text embeddings.
type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"` // List of embedding vectors, one for each input text
}

// Client defines the interface for interacting with LLM providers.
// This interface allows the application to remain decoupled from any specific API (e.g. Ollama, OpenAI).
type Client interface {
	// Chat sends a full chat request and returns the final response.
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)

	// ChatStream sends a chat request and stream chunks of the response back via onChunk callback.
	// Returns the complete ChatResponse once the stream completes.
	ChatStream(ctx context.Context, req ChatRequest, onChunk func(string) error) (ChatResponse, error)

	// Embed generates vector embeddings for the given input strings.
	Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error)
}
