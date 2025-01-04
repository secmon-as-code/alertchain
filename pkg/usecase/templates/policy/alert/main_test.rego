package alert.your_schema

test_detect_your_schema if {
	resp := alert with input as data.alert.testdata.your_schema
	count(resp) > 0
}

test_not_detect_your_schema if {
	resp := alert with input as json.patch(data.alert.testdata.your_schema, [{"op": "replace", "path": "/severity", "value": ["LOW"]}])
	count(resp) == 0
}
