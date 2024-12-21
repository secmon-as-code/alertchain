package alert.aws_guardduty

alert contains res if {
	startswith(input.Findings[x].Type, "Trojan:")
	input.Findings[_].Severity > 7
	res := {
		"title": input.Findings[x].Type,
		"source": "aws",
	}
}
