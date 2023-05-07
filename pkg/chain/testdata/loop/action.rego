package action

run[res] {
	p := input.alert.params[_]
	p.name == "c"
	p.value < 10

	res := {"uses": "mock"}
}

exit[res] {
	p := input.alert.params[_]
	p.name == "c"

	res := {"params": {{
		"id": p.id,
		"name": "c",
		"value": p.value + 1,
	}}}
}
