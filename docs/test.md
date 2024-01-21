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

Action Policy is a policy that controls the behavior of actions. As such, testing its behavior requires interactions with external services. However, using responses from external services directly in tests can be inconvenient due to constraints such as inconsistent responses or difficulty in preparing expected answers. To address this, AlertChain has implemented a "play" mode. In play mode, you can pre-define a **Scenario**, which describes workflow specifying how actions should respond.

The play mode itself is not for verifying the behavior of the policy; it only logs the execution results. However, by testing these logs using OPA/Rego, you can verify how the Action Policy behaved based on the responses obtained from each action. This achieves the "Automatic test for orchestration and automated response," which is one of the challenges in SOAR implementation.

### Playbook

Here is an example of a Scenario jsonnet file:

```jsonnet
{
  id: 'scenario1',
  title: 'Test 1',
  events: [
    {
      input: import 'event/guardduty.json',
      schema: 'aws_guardduty',
      actions: {
        'chatgpt.comment_alert': [
          import 'results/chatgpt.json',
        ],
      },
    },
  ],
  env: {
    CHATGPT_API_KEY: 'test_api_key_xxxxxxxxxx',
    SLACK_WEBHOOK_URL: 'https://hooks.slack.com/services/xxxxxxxxx',
  },
}
```

A scenario is composed of the following fields:


- `id`: Specify any string, ensuring it is unique within the playbook. This serves as a key to identify the scenario when writing tests using Rego.
- `title`: Specify any string. This is used to describe the scenario for human readability.
- `events`: This describes scenarios for each event.
  - `input`: This field specifies the event data to be used for the scenario.
  - `schema`: This field specifies the schema to be used for the scenario.
  - `actions`: This field contains the expected results for each action involved in the scenario. The results are defined as key-value pairs, where the key represents the action Name and the value is an array of expected responses for that action.
- `env`: Environment variables that will be used in play mode.

By defining multiple scenarios within the playbook, you can effectively test various use cases and ensure that your Action Policy behaves as expected under different circumstances. This allows for comprehensive testing and validation of your SOAR implementation, leading to more robust and reliable automated response systems.

### Testing logs with OPA/Rego

To execute the play mode with the policy, refer to the [examples](../examples/test) directory. Run the following command:

```bash
$ alertchain -d ./policy/ play -b playbook.jsonnet
```

This will generate a file named `output/scenario1/data.json`. A sample file would look like this:

```json
{
  "id": "scenario1",
  "title": "Test 1",
  "results": [
    {
      "alert": {
        "title": "Trojan:EC2/DropPoint!DNS",
        "description": "EC2 instance i-99999999 is querying a domain name of a remote host that is known to hold credentials and other stolen data captured by malware.",
        "source": "aws",
        "attrs": [
          {
            "id": "e6fc6cbc-dd90-47b8-a73f-b392a53addcd",
            "key": "instance ID",
            "value": "i-99999999",
            "type": ""
          }
        ],
        "id": "9d64a4b4-15c2-4e64-a4ff-6af253b80b95",
        "schema": "aws_guardduty",
        "namespace": "",
        "created_at": "2023-05-07T12:17:56.929635+09:00"
      },
      "actions": [
        {
          "seq": 0,
          "init": [
            {
              "attrs": [
                {
                  "key": "test",
                  "value": "this is a test"
                }
              ]
            }
          ]
          "run": [
            {
              "id": "ask-gpt",
              "uses": "chatgpt.query",
              "args": {
                "secret_api_key": "test_api_key_xxxxxxxxxx"
              }
            },
          ],
        },
        {
          "seq": 1,
          "run": [
            {
              "id": "notify-slack",
              "uses": "slack.post",
              "args": {
                "body": "This is a test message.",
                "channel": "alert",
                "secret_url": "https://hooks.slack.com/services/xxxxxxxxx"
              }
            }
          ],
          "exit": [
            {
              "attrs": [
                {
                  "key": "counter",
                  "value": 1,
                  "global": true
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

One JSON file will be generated for each scenario.

#### Schema of playbook result

##### root object

- `id`: Provided from playbook
- `title`: provided from playbook
- `results`: Array of [Result](#result) for each event

##### `Result`

- `alert`: [Alert](policy.md#alert) object. This is the alert object that was generated by the Alert Policy and does not have any additional attributes.
- `actions`: Array of [Action](#action) objects. This is the action object that was generated by the Action Policy and does not have any additional attributes. Array order is the same as the order of action policy evaluation sequence.

##### `Action`

- `init`: Array of the result of `init` rule evaluation.
- `run`: Array of the result of `run` rule evaluation.
- `exit`: Array of the result of `exit` rule evaluation.
- `seq`: Sequence number of the action. It starts from 0 and is incremented by 1 for each action.

#### Testing with Rego

You can test if the actions were triggered as expected by examining these JSON files. You can use any language or framework for testing, but in this case, we will use Rego.

**test.rego**
```rego
package test

test_play_result {
    # Only one alert should be detected
    count(data.output.scenario1.results) == 1
    result := data.output.scenario1.results[0]

    # Alert should be of type "Trojan:EC2/DropPoint!DNS"
    result.alert.title == "Trojan:EC2/DropPoint!DNS"

    # The alert should trigger two actions
    count(result.actions) == 2

    # test first action
    first := result.actions[_]
    first.seq == 0
    first.run[r1].id == "ask-gpt"
    first.run[r1].args.secret_api_key == "test_api_key_xxxxxxxxxx"

    # test second action
    second := result.actions[_]
    second.seq == 1
    second.run[r2].id == "notify-slack"
    second.run[r2].args.secret_url == "https://hooks.slack.com/services/xxxxxxxxx"
}
```

Once you have prepared this file, you can run the test using the OPA command as follows:

```bash
$ opa test -v .
test.rego:
data.test.test_play_result: PASS (323.125µs)
--------------------------------------------------------------------------------
PASS: 1/1
```

Using this approach, you can continuously and automatically inspect whether the entire workflow is functioning correctly.