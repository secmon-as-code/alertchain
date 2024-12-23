package action

run contains res if {
	input.alert.attrs[k2].key == "k2"
	res := {
		"id": "run_mock",
		"uses": "mock",
		"commit": [
			{
				"id": input.alert.attrs[k2].id,
				"key": "k2",
				"value": "v2a",
			},
			{
				"key": "k3",
				"value": "v3",
			},
		],
	}
}

run contains res if {
	input.called[_].id == "run_mock"
	res := {
		"id": "run2",
		"uses": "mock.after",
	}
}
