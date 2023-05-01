package alert.aws_guardduty

alert[res] {
	startswith(input.Findings[x].Type, "Trojan:")
	input.Findings[_].Severity > 7
	res := {
		"title": input.Findings[x].Type,
		"source": "aws",
	}
}
