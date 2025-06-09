# LLM Proxy

A reverse proxy server for OpenRouter.ai that adds prefill functionality.

## What is Prefill?

Prefill allows you to start the LLM's response with predefined text, controlling the format and style of the AI's output.

## Why Use This Proxy?

Many LLM clients don't support native prefill configuration. This proxy enables prefill functionality from any message in the conversation - user messages, system prompts, or any other message type.

## Compilation

```
go build
```

## Usage

1. **Start the server:**
   ```
   llm_proxy.exe
   ```
   The server will start on port 8080.

2. **Configure your OpenRouter API key:**
   The proxy provides OpenAI-compatible endpoints. Include your OpenRouter API key in the `Authorization` header when making requests to `http://localhost:8080`.

## Features

### Prefill

If any message contains a `prefill` block using `[[` and `]]` delimiters, the proxy will:

1. Remove the prefill block from the original message content
2. Add a new assistant message with the prefill content to the conversation

**Example:**

Input message:
```json
{
  "role": "user", 
  "content": "Hello [[ define \"prefill\" ]]I'm an AI assistant[[ end ]]"
}
```

The proxy will:
- Modify the user message content to: "Hello "
- Add an assistant message: `{"role": "assistant", "content": "I'm an AI assistant"}`