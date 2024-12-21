package action

run contains job if {
	job := {
		"id": "my_job",
		"uses": "mock",
		"commit": [
			{
				"id": attr.id,
				"key": "counter",
				"value": attr.value + 1,
				"persist": true,
			},
		],
	}
}

attr := input.alert.attrs[x] if {
	input.alert.attrs[x].key == "counter"
} else := init if {
	init := {
		"id": null,
		"key": "counter",
		"value": 0,
	}
}
