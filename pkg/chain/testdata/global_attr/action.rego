package action

init[res] {
	input.alert.attrs[x].key == "status"
	input.alert.attrs[x].value == "done"
	res := {"abort": true}
}

run[job] {
	job := {
		"id": "my_job",
		"uses": "mock",
	}
}

exit[res] {
	res := {"attrs": [{"key": "status", "value": "done", "global": true}]}
}
