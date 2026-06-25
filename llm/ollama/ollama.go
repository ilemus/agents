package ollama

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"agents/llm"

	"github.com/ollama/ollama/api"
)

// This is an example default ollama option for a highly predictable and consistent response.
var ollamaOptions = map[string]interface{}{
	"temperature": 0.2, // Low temperature for more deterministic responses
	"top_p":       0.9, // Nucleus sampling
	"top_k":       20,  // Limit pool of top tokens
	// "num_predict": 1000,  // Max tokens to generate
	// Note: "reasoning_effort" or similar parameters depend on the specific model support (e.g., DeepSeek R1)
}

// This is an example default ollama option for a predictable but creative response.
var ollamaCreativeOptions = map[string]interface{}{
	"temperature": 0.7,  // Balanced creativity
	"top_p":       0.85, // High nucleus sampling for diverse token choices
	"top_k":       35,   // Wider pool of tokens to choose from
	"num_predict": 2000, // Max tokens to generate, try to limit unlimited creativity.
}

// Client implements the llm.Client interface using the official Ollama API client.
type Client struct {
	apiClient *api.Client
}

// NewClient creates a new Ollama API client.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		u, _ = url.Parse("http://localhost:11434")
	}
	return &Client{
		apiClient: api.NewClient(u, http.DefaultClient),
	}
}

// Chat sends a full chat completion request to the Ollama endpoint.
func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	apiMessages := make([]api.Message, len(req.Messages))
	for i, msg := range req.Messages {
		apiMessages[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	stream := false
	apiReq := &api.ChatRequest{
		Model:    req.Model,
		Messages: apiMessages,
		Stream:   &stream,
		Options:  req.Options,
	}

	var finalResp llm.ChatResponse

	fn := func(resp api.ChatResponse) error {
		finalResp.Message = llm.Message{
			Role:    resp.Message.Role,
			Content: resp.Message.Content,
		}
		return nil
	}

	err := c.apiClient.Chat(ctx, apiReq, fn)
	if err != nil {
		return llm.ChatResponse{}, err
	}

	return finalResp, nil
}

// ChatStream sends a chat request and streams output chunks back to the client.
func (c *Client) ChatStream(ctx context.Context, req llm.ChatRequest, onChunk func(string) error) (llm.ChatResponse, error) {
	apiMessages := make([]api.Message, len(req.Messages))
	for i, msg := range req.Messages {
		apiMessages[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	stream := true
	apiReq := &api.ChatRequest{
		Model:    req.Model,
		Messages: apiMessages,
		Stream:   &stream,
		Options:  req.Options,
	}

	var sb strings.Builder
	var role string

	fn := func(resp api.ChatResponse) error {
		if role == "" && resp.Message.Role != "" {
			role = resp.Message.Role
		}
		if resp.Message.Content != "" {
			sb.WriteString(resp.Message.Content)
			if onChunk != nil {
				if err := onChunk(resp.Message.Content); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err := c.apiClient.Chat(ctx, apiReq, fn)
	if err != nil {
		return llm.ChatResponse{}, err
	}

	if role == "" {
		role = "assistant"
	}

	return llm.ChatResponse{
		Message: llm.Message{
			Role:    role,
			Content: sb.String(),
		},
	}, nil
}

// Embed generates vector embeddings for the given input strings.
func (c *Client) Embed(ctx context.Context, req llm.EmbedRequest) (llm.EmbedResponse, error) {
	apiReq := &api.EmbedRequest{
		Model:   req.Model,
		Input:   req.Input,
		Options: req.Options,
	}

	resp, err := c.apiClient.Embed(ctx, apiReq)
	if err != nil {
		return llm.EmbedResponse{}, err
	}

	return llm.EmbedResponse{
		Embeddings: resp.Embeddings,
	}, nil
}
