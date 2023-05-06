# ChatGPT

## `chatgpt.comment_alert`

This action provides a summary and suggested response to a security alert using ChatGPT.

## Prerequisite

1. Create an account on OpenAI.
2. Generate an API key at https://platform.openai.com/account/api-keys.

### Arguments

Example policy:

```rego
run[res] {
  res := {
    "id": "your-action",
    "uses": "chatgpt.comment_alert",
    "args": {
      "secret_api_key": input.env.CHATGPT_API_KEY,
    },
  },
}
```

- `secret_api_key` (required, string): The API key for OpenAI.
- `prompt` (optional, string): The ChatGPT prompt. The default is:

```
Please analyze and summarize the given JSON-formatted security alert data, and suggest appropriate actions for the security administrator to respond to the alert:
```

## Response

The response format follows the OpenAI Chat API guidelines, which can be found at:
https://platform.openai.com/docs/guides/chat/response-format