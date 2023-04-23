# ChatGPT Analyst

## Prerequisite

1. Create your account of OpenAI
2. Create a API key at https://platform.openai.com/account/api-keys

## Config

Example
```jsonnet
    {
      id: 'inquiry_chatgpt',
      use: 'chatgpt-analyst',
      config: {
        api_key: std.extVar('CHATGPT_API_KEY'),
      },
    },
```

### `api_key` (required, string)

The API key of OpenAI.

### `prompt` (optional, string)

Prompt of chat GTP. The default is following:

```
Summarize the following json formatted data of security alert and propose security administrator's action:
```

## Arguments

N/A

## Response

https://platform.openai.com/docs/guides/chat/response-format
