package action

run contains {
	"id": "my_job",
	"uses": "mock",
	"commit": [
		{
			"key": "status",
			"value": "done",
			"persist": true,
		},
	],
} if {
	not is_done
}

is_done if {
	input.alert.attrs[x].key == "status"
	input.alert.attrs[x].value == "done"
}
