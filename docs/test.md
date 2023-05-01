# Test

AlertChain's basic concept is based on Policy as Code. Orchestration and automation of AlertChain can be described by policy language Rego. Testing is a crucial component of Policy as Code for the following reasons:

1. Improved reliability: By conducting tests, you can ensure that policies are accurately coded and function as intended. This reduces the risk of policy violations, increasing the overall reliability of the system.

2. Continuous improvement: Repeated testing helps identify weaknesses and areas for improvement in policies, enabling continuous refinement. This continuously enhances the organization's security posture.

3. Automation and scalability: Automating tests allows for the rapid and efficient application and auditing of policies. This enables security measures that scale with the organization's growth.

4. Documentation function: Test cases also serve as documentation, clearly outlining policy requirements and expected behavior. This makes it easier for other teams and members within the organization to understand and appropriately respond to policies.

In summary, testing is essential for the reliability, effectiveness, continuous improvement, and improvement of the overall security culture within an organization in the context of Policy as Code.

## Testing Alert Policy

Testing an Alert Policy can be written in the same way as for a normal OPA/Rego test: the value of "input" when evaluating an Alert Policy is the event data input from the outside, so you can override this "input" and test the policy you have written. Override this "input" to test the policy you wrote.

### Preparing policy and data files

It is recommended to prepare the test by creating an `alert.rego` file containing the Alert Policy, as well as a JSON file containing the event data. In this example, we will explain using an AWS GuardDuty detection event.

**alert.rego**
```rego
package alert.aws_guardduty

alert[res] {
    startswith(input.Findings[x].Type, "Trojan:")
    input.Findings[_].Severity > 7
    res := {
        "title": input.Findings[x].Type,
        "source": "aws",
    }
}
```

Obtain a sample AWS GuardDuty event, either from the AWS documentation or by triggering a test event in your AWS environment. And save the event data in a JSON file as `test/aws_guardduty/data.json`. Please note OPA can read only 1 json file in a directory.

**test/aws_guardduty/data.json**
```json
{
    "Findings": [
        {
            "Type": "Trojan:EC2/DriveBySourceTraffic!DNS",
            "Region": "us-east-1",
            "Severity": 8,
            (snip)
        }
    ]
}
```

For a detailed explanation of how to write tests, please refer to the [OPA official documentation](https://www.openpolicyagent.org/docs/latest/policy-testing/). In the example below, we use sample data (test/aws_guardduty/data.json) that correctly detects violations. We test not only for successful detection cases but also for ignore cases by modifying parts of the data using the `json.patch()` built-in function. Of course, it is also acceptable to prepare and test data for each individual case.

```rego
package alert.aws_guardduty

# detect alert correctly
test_detect {
	result := alert with input as data.test.aws_guardduty
	count(result) == 1
	result[_].title == "Trojan:EC2/DriveBySourceTraffic!DNS"
	result[_].source == "aws"
}

# ignore if severity is 7
test_ignore_severity {
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
test_ignore_type {
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
```

## Testing Action Policy

Action Policy is a policy that controls the behavior of actions. As such, testing its behavior requires interactions with external services. However, using responses from external services directly in tests can be inconvenient due to constraints such as inconsistent responses or difficulty in preparing expected answers. To address this, AlertChain has implemented a "play" mode. In play mode, you can pre-define a playbook, which describes scenarios specifying how actions should respond.

The play mode itself is not for verifying the behavior of the policy; it only logs the execution results. However, by testing these logs using OPA/Rego, you can verify how the Action Policy behaved based on the responses obtained from each action. This achieves the "Automatic test for orchestration and automated response," which is one of the challenges in SOAR implementation.

### Playbook

Here is an example of a Playbook jsonnet file:

```jsonnet
{
  scenarios: [
    {
      id: 'aws_guardduty_test1',
      alert: import 'aws_guardduty/data.json',
      schema: 'aws_guardduty',
      results: {
        query_chatgpt: [
          {
            choices: [
              {
                message: {
                  role: 'assistant',
                  content: 'this is a test message',
                },
              },
            ],
          },
        ],
      },
    },
  ],
}
```

A scenario is composed of the following fields:

- `id`: Specify any string, ensuring it is unique within the playbook. This serves as a key to identify the scenario when writing tests using Rego.
- `alert`: This field is used to import the alert data required for the scenario. The data is usually stored in a separate JSON file, which is imported using the `import` keyword. In this example, the alert data is imported from the 'aws_guardduty/data.json' file.
- `schema`: This field specifies the schema to be used for the scenario. The schema specifies a policy to evaluate the alert. In the given example, the schema is set to 'aws_guardduty'.
- `results`: This field contains the expected results for each action involved in the scenario. The results are defined as key-value pairs, where the key represents the action ID and the value is an array of expected responses for that action. In the example provided, the `query_chatgpt` action has an expected response containing a message with the role 'assistant' and the content 'this is a test message'.

By defining multiple scenarios within the playbook, you can effectively test various use cases and ensure that your Action Policy behaves as expected under different circumstances. This allows for comprehensive testing and validation of your SOAR implementation, leading to more robust and reliable automated response systems.

### Testing logs with OPA/Rego

```rego
package action.main

action[msg] {
    msg := {
        "id": "query_chatgpt",
    }
}
```

```rego
package action.query_chatgpt

action[msg] {
    msg := {
        "id": "create_github_issue",
        "params": [
            {
                "key": "analyst comment",
                "value": input.result.choices[0].message.content,
            }
        ]
    }
}
```


```json
{
  "id": "aws_guardduty_test1",
  "title": "",
  "alerts": [
    {
      "alert": {
        "title": "Trojan:EC2/DriveBySourceTraffic!DNS",
        "description": "",
        "source": "aws",
        "params": [],
        "id": "4671c703-5c83-4cfe-a956-f1581be6406c",
        "schema": "aws_guardduty",
        "created_at": "2023-05-01T09:22:20.02218+09:00"
      },
      "created_at": 22416000,
      "actions": [
        {
          "action": {
            "id": "inquery_chatgpt",
            "args": null,
            "params": []
          },
          "next": [
            {
              "id": "create_github_issue",
              "args": null,
              "params": [
                {
                  "key": "analyst comment",
                  "value": "this is a test message",
                  "type": ""
                }
              ]
            }
          ],
          "started_at": 22549000,
          "ended_at": 22781000
        },
        {
          "action": {
            "id": "create_github_issue",
            "args": null,
            "params": [],
          },
          "next": null,
          "started_at": 22783000,
          "ended_at": 22910000
        }
      ]
    }
  ]
}
```