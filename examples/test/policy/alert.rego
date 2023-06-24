package alert.aws_guardduty

alert[res] {
	f := input.Findings[_]
	startswith(f.Type, "Trojan:")
	f.Severity > 7

	res := {
		"title": f.Type,
		"source": "aws",
		"description": f.Description,
		"attrs": [{
			"key": "instance ID",
			"value": f.Resource.InstanceDetails.InstanceId,
		}],
	}
}
