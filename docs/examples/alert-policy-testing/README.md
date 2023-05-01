# Example of alert policy testing

This is an example of AWS GuardDuty alert policy and testing it.

## Files

- [alert.rego](alert.rego): Alert policy
- [alert_test.rego](alert_test.rego): Testing policy
- [test/aws_guardduty/data.json](test/aws_guardduty/data.json): Testing data

## How to test

```bash
$ opa test -v -b .
alert_test.rego:
data.alert.aws_guardduty.test_detect: PASS (407.709µs)
data.alert.aws_guardduty.test_ignore_severity: PASS (235.25µs)
data.alert.aws_guardduty.test_ignore_type: PASS (187.459µs)
--------------------------------------------------------------------------------
PASS: 3/3
```
