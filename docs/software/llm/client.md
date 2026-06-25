
# LLM Client Interface and Provider Adapter Specification

This document details the software design of the language model (LLM) client architecture in the agents framework. 

The LLM abstraction is designed to decouple the application logic from any specific LLM provider API (such as Ollama, OpenAI, Anthropic, or Gemini), facilitating seamless drop-in replacements.

---

## 1. Architectural Overview

To ensure that the agentic application remains provider-agnostic, all LLM communication relies on a shared interface defined in the `llm` package. Concrete implementations (such as the `ollama` client adapter) are located in dedicated subdirectories and implement this interface.

```mermaid
graph TD
    subgraph Main Application
        App[Agentic Workflows / CLI]
    end

    subgraph LLM Package
        Interface[llm.Client Interface]
        Types[Data Types: ChatRequest, EmbedRequest, etc.]
    end

    subgraph Provider Implementation
        Ollama[llm/ollama.Client]
        OpenAI[llm/openai.Client (Future)]
        Gemini[llm/gemini.Client (Future)]
    end

    App -->|Uses| Interface
    App -->|Constructs| Types
    Ollama -.->|Implements| Interface
    OpenAI -.->|Implements| Interface
    Gemini -.->|Implements| Interface
```

---

## 2. Core Abstractions: [client.go](agents/llm/client.go)

The package [llm](agents/llm) defines the primary structs and interface that the application uses for text generation and vector embeddings.

### Data Types

- **`Message`**: Represents a single turn in a chat conversation, consisting of a `Role` (such as `"system"`, `"user"`, or `"assistant"`) and the text `Content`.
- **`ChatRequest`**: Standardizes the chat payload. It includes the `Model` name, a slice of conversational `Messages`, and a map of provider-specific optional parameters (`Options`).
- **`ChatResponse`**: Standardizes the LLM reply, wrapping a single generated `Message`.
- **`EmbedRequest`**: Standardizes vector embedding requests. It accepts a `Model`, a slice of inputs (`Input`), and provider-specific `Options`.
- **`EmbedResponse`**: Wraps the resulting slices of float32 vectors (`Embeddings`).

### The `Client` Interface

The `Client` interface defines three primary integration points:

```go
type Client interface {
	// Chat sends a full chat request and returns the final response.
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)

	// ChatStream sends a chat request and streams chunks of the response back via onChunk callback.
	// Returns the complete ChatResponse once the stream completes.
	ChatStream(ctx context.Context, req ChatRequest, onChunk func(string) error) (ChatResponse, error)

	// Embed generates vector embeddings for the given input strings.
	Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error)
}
```

---

## 3. Concrete Provider: [ollama.go](agents/llm/ollama/ollama.go)

The [ollama](agents/llm/ollama) package contains the concrete implementation of `llm.Client` using the official Ollama client library `github.com/ollama/ollama/api`.

### Key Design Aspects

1. **Client Initialization**:
   The client is initialized using `NewClient(baseURL string)`. If no base URL is provided, it defaults to the standard local address `http://localhost:11434`.
   ```go
   func NewClient(baseURL string) *Client {
       // Defaults to localhost:11434 if empty, ensures http/https scheme prefix
       // Initializes the internal Ollama api.Client
   }
   ```
2. **Translation Layer**:
   The methods (`Chat`, `ChatStream`, `Embed`) map the domain-agnostic `llm.ChatRequest`/`llm.EmbedRequest` structs into Ollama's native API parameters (`api.ChatRequest` / `api.EmbedRequest`), preventing external dependencies from leaking into the core application logic.
3. **Stream Handling**:
   The `ChatStream` implementation wraps Ollama's streaming callback, accumulating chunks into a string builder while simultaneously dispatching them via `onChunk(string)`. Once streaming completes, it returns a consolidated `llm.ChatResponse`.
4. **Configuration Options**:
   Preset parameters are defined to control the LLM's creativity and predictability:
   - **`ollamaOptions`**: Deterministic/highly predictable config (temperature `0.2`, top-p `0.9`, top-k `20`).
   - **`ollamaCreativeOptions`**: Balanced creativity config (temperature `0.7`, top-p `0.85`, top-k `35`).

---

## 4. Drop-in Extensibility (e.g. OpenAI / Gemini)

To replace or supplement Ollama with a different provider (e.g., OpenAI), one would only need to:
1. Create a new directory `/llm/openai`.
2. Define a client struct that wraps the OpenAI SDK (or uses custom HTTP requests).
3. Implement the `llm.Client` interface methods:
   ```go
   package openai

   import (
       "context"
       "agents/llm"
   )

   type Client struct {
       // OpenAI client config & SDK handles
   }

   func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) { ... }
   func (c *Client) ChatStream(ctx context.Context, req llm.ChatRequest, onChunk func(string) error) (llm.ChatResponse, error) { ... }
   func (c *Client) Embed(ctx context.Context, req llm.EmbedRequest) (llm.EmbedResponse, error) { ... }
   ```
4. Instantiate this `openai.Client` at the application entrypoint. Since it satisfies `llm.Client`, it can be passed directly to any workspace or workflow logic without modifying downstream functions.
