package action

run[job] {
	job := {
		"id": "my_job",
		"uses": "mock",
		"commit": [
			{
				"id": attr.id,
				"key": "counter",
				"value": attr.value + 1,
				"global": true,
			},
		],
	}
}

attr := input.alert.attrs[x] {
	input.alert.attrs[x].key == "counter"
} else := init {
	init := {
		"id": null,
		"key": "counter",
		"value": 0,
	}
}
