package action

run[res] {
	input.alert.attrs[x].key == "status"
	input.alert.attrs[x].value == "done"
	res := {"abort": true}
}

run[job] {
	job := {
		"id": "my_job",
		"uses": "mock",
		"commit": [
			{
				"key": "status",
				"value": "done",
				"global": true,
			},
		],
	}
}
