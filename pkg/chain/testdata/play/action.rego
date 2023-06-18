package action

run[res] {
	p := input.alert.attrs[_]
	p.name == "c"
	p.value < 3

	res := {"uses": "mock"}
}

exit[res] {
	p := input.alert.attrs[_]
	p.name == "c"

	res := {"attrs": [{
		"id": p.id,
		"name": "c",
		"value": p.value + 1,
	}]}
}
