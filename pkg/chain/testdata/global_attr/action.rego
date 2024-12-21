package action

run contains res if {
	input.alert.attrs[x].key == "status"
	input.alert.attrs[x].value == "done"
	res := {"abort": true}
}

run contains job if {
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
