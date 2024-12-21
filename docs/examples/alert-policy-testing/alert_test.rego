package alert.aws_guardduty

# detect alert correctly
test_detect if {
	result := alert with input as data.test.aws_guardduty
	count(result) == 1
	result[_].title == "Trojan:EC2/DriveBySourceTraffic!DNS"
	result[_].source == "aws"
}

# ignore if severity is 7
test_ignore_severity if {
	result := alert with input as json.patch(
		data.test.aws_guardduty,
		[{
			"op": "replace",
			"path": "/Findings/0/Severity",
			"value": 7,
		}],
	)
	count(result) == 0
}

# ignore if prefix of Type does not match with "Trojan:"
test_ignore_type if {
	result := alert with input as json.patch(
		data.test.aws_guardduty,
		[{
			"op": "replace",
			"path": "/Findings/0/Type",
			"value": "Some alert",
		}],
	)
	count(result) == 0
}
