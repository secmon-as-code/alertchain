{
  scenarios: [
    {
      id: 'scenario1',
      title: 'Test 1',
      event: import 'event/guardduty.json',
      schema: 'aws_guardduty',
      results: {
        'chatgpt.comment_alert': [
          import 'results/chatgpt.json',
        ],
      },
    },
  ],
  env: {
    CHATGPT_API_KEY: 'test_api_key_xxxxxxxxxx',
    SLACK_WEBHOOK_URL: 'https://hooks.slack.com/services/xxxxxxxxx',
  },
}
