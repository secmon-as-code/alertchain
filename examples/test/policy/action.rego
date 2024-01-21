package action

run[res] {
	input.alert.source == "aws"
	res := {
		"id": "ask-gpt",
		"uses": "chatgpt.query",
		"args": {"secret_api_key": input.env.CHATGPT_API_KEY},
	}
}

run[res] {
	gtp := input.called[_]
	gtp.id == "ask-gpt"

	res := {
		"id": "notify-slack",
		"uses": "slack.post",
		"args": {
			"secret_url": input.env.SLACK_WEBHOOK_URL,
			"channel": "alert",
			"body": gtp.result.choices[0].message.content,
		},
	}
}
