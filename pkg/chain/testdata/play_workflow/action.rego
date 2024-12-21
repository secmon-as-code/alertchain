package action

# init[res] {
# 	input.seq == 0
# 	res := {"attrs": [{
# 		"key": "color",
# 		"value": "blue",
# 	}]}
# }

run contains job if {
	input.seq == 0
	job := {
		"id": "1st",
		"uses": "mock",
		"args": {"tick": 1},
	}
}

# exit[job] {
# 	input.action.id == "1st"
# 	print(input.action.result)

# 	job := {"attrs": [{
# 		"key": "index1",
# 		"value": input.action.result.index,
# 	}]}
# }

run contains job if {
	input.seq == 1
	job := {
		"id": "2nd",
		"uses": "mock",
		"args": {"tick": 2},
	}
}
