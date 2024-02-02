package action

run[job] {
	input.seq == 0
	job := {
		"id": "force_continue",
		"uses": "mock",
		"args": {"step": 1},
        "force": true,
	}
}

run[job] {
	input.seq == 1
	job := {
		"id": "stop_by_error",
		"uses": "mock",
		"args": {"step": 2},
        "force": false,
	}
}

run[job] {
	input.seq == 1
	job := {
		"id": "not_run",
		"uses": "mock",
		"args": {"step": 3},
        "force": false,
	}
}
