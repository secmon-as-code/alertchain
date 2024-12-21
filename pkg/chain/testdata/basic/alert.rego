package alert.scc

alert contains msg if {
	msg := {
		"title": input.finding.category,
		"attrs": [
			{
				"key": "db_name",
				"value": input.finding.database.displayName,
			},
			{
				"key": "db_user",
				"value": input.finding.database.userName,
			},
			{
				"key": "db_query",
				"value": input.finding.database.query,
			},
		],
	}
}
