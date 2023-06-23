package action

run[res] {
	p := input.alert.attrs[_]
	p.key == "c"
	p.value < 3

	res := {"uses": "mock"}
}

exit[res] {
	p := input.alert.attrs[_]
	p.key == "c"

	res := {"attrs": [{
		"id": p.id,
		"key": "c",
		"value": p.value + 1,
	}]}
}
