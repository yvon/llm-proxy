# LLM Proxy

A reverse proxy server for OpenRouter.ai that adds prefill functionality and cache control for Anthropic models.

## Why Use This Proxy?

Many LLM clients don't support native prefill configuration or cache control. This proxy enables these advanced features from any message in the conversation - user messages, system prompts, or any other message type.

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

If any message contains a `prefill` block using `|prefill:` and `|` delimiters, the proxy will:

1. Remove the prefill block from the original message content
2. Add a new assistant message with the prefill content to the conversation

**Example:**

Suppose you send the following message:

```
Hello |prefill: I'm an AI assistant|
```

This would normally be sent as:

```json
{
  "role": "user",
  "content": "Hello |prefill: I'm an AI assistant|"
}
```

The proxy will transform it to:

```json
[
  {
    "role": "user",
    "content": "Hello "
  },
  {
    "role": "assistant",
    "content": "I'm an AI assistant"
  }
]
```

### Cache Control (Anthropic Models)

If any message contains the `|cache|` command, the proxy will add cache control to the last user message for Anthropic models. This enables [prompt caching functionality](https://docs.anthropic.com/en/docs/build-with-claude/prompt-caching) to reduce costs for repeated prompts.

**Example:**

```
Your message content |cache|
```

The proxy will transform it to:

```json
{
  "role": "user",
  "content": "Your message content",
  "cache_control": {
    "type": "ephemeral"
  }
}
```
