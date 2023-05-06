# Slack Integration

## `slack.post`

This action posts a message to a Slack channel.

## Prerequisites

To use this action, you need to create a Slack App and obtain an incoming webhook URL. You can find instructions for setting up a Slack App and generating a webhook URL [here](https://api.slack.com/messaging/webhooks).

### Arguments

Here's an example policy using the `slack.post` action:

```rego
run[res] {
  res := {
    "id": "your-action",
    "uses": "slack.post",
    "args": {
      "secret_url": input.env.SLACK_WEBHOOK_URL,
      "channel": "alert",
    },
  },
}
```

- `secret_url` (required, string): The Slack webhook URL used to post messages to your Slack channel.
- `channel` (required, string): The name of the Slack channel where the message will be posted. The `#` symbol is not required.
- `text` (optional, string): The title of the Slack message. The default value is `Notification from AlertChain`.
- `body` (optional, string): The body of the Slack message. The default value is the alert title and description.
- `color` (optional, string): The color of the Slack message's banner. The default value is `#2EB67D`.

## Response

This action does not return a response.
