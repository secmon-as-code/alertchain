package action

run[job] {
	job := {
		"id": "my_job",
		"uses": "mock",
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

exit[res] {
	res := {"attrs": [{
		"id": attr.id,
		"key": "counter",
		"value": attr.value + 1,
		"global": true,
	}]}
}
