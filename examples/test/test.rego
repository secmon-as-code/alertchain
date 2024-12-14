package test

test_play_result {
	# Only one alert should be detected
	count(data.output.scenario1.results) == 1
	result := data.output.scenario1.results[0]

	# Alert should be of type "Trojan:EC2/DropPoint!DNS"
	result.alert.title == "Trojan:EC2/DropPoint!DNS"

	# The alert should trigger two actions
	count(result.actions) == 2

	# test first action
	result.actions[0].id == "ask-gpt"
	result.actions[0].args.secret_api_key == "test_api_key_xxxxxxxxxx"

	# test second action
	result.actions[1].id == "notify-slack"
	result.actions[1].args.secret_url == "https://hooks.slack.com/services/xxxxxxxxx"
}
