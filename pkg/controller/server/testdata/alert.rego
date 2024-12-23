package alert.test_service

alert contains msg if {
	msg := {
		"title": "test alert",
		"description": "test description",
		"attrs": {{
			"key": "test_attr",
			"value": "test_value",
		}},
	}
}
