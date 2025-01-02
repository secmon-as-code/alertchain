package action

run contains job if {
	input.seq == 0
	job := {
		"id": "1st",
		"uses": "mock",
		"args": {"tick": 1},
	}
}

run contains job if {
	input.seq == 1
	job := {
		"id": "2nd",
		"uses": "mock",
		"args": {"tick": 2},
	}
}
