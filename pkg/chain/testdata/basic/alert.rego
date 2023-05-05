package alert.scc

alert[msg] {
    msg := {
        "title": input.finding.category,
        "params": [
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
        ]
    }
}
