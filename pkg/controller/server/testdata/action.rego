package action

run contains job if {
	job := {
		"id": "test",
		"uses": "http.fetch",
		"args": {
			"method": "GET",
			"url": "https://emhkq5vqrco2fpr6zqlctbjale0eyygt.lambda-url.ap-northeast-1.on.aws",
		},
		"commit": [
			{
                "key": "added_attr",
                "value": "swirls",
			},
		],
	}
}
