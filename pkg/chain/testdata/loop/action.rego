package action

run[res] {
	p := input.alert.attrs[_]
	p.key == "c"
	p.value < 10

	res := {
		"uses": "mock",
		"commit": [
			{
				"id": p.id,
				"key": "c",
				"value": p.value + 1,
			},
		],
	}
}
