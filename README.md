# AlertChain

AlertChain is a simple SOAR (Security Orchestration, Automation, and Response) framework that leverages [OPA](https://github.com/open-policy-agent/opa) (Open Policy Agent) to enhance security management.

![AlertChain Diagram](https://user-images.githubusercontent.com/605953/232273906-a3df56fb-3201-4336-b897-e327e8d49981.jpg)

## Motivation

Security Orchestration, Automation, and Response (SOAR) is a platform designed for automating the detection, analysis, and response of security events. In order to enable automated event analysis and rapid response in SOAR systems, it is essential to execute automated security procedures and policies.

By utilizing OPA and Rego, a SOAR system can flexibly apply a set of user-defined policies to maintain the security of applications and systems. This approach simplifies the process of updating or changing security policies and ensures a more accurate policy application. Moreover, the Rego language is flexible and expressive, making it easy to add or modify policies.

## Concept

AlertChain is a versatile software that accepts structured event data through HTTP or other means, and then determines its actions based on policies written in Rego.

### Action

Actions, the basic units of operation, are primarily implemented within AlertChain using Go code. For example, there is an action called [chatgpt.comment_alert](pkg/action/chatgpt/README.md#chatgptcomment_alert) which creates an issue on GitHub. Users can define any number of actions in a configuration written in Rego, each of which needs a unique ID.

### Policy

There are two main types of policies in AlertChain, Alert Policy and Action Policy.

1. **Alert Policy**: Responsible for determining whether the incoming event data from external sources should be treated as an alert or not. For example, when receiving notifications from external services, you may want to handle only alerts related to specific categories, or you may want to exclude events that meet certain conditions (such as specific users or hosts). The Alert Policy can be used to achieve these goals by excluding certain events or including only specific events as alerts.
2. **Action Policy**: Determines the appropriate response for detected alerts. For example, when an issue is detected on a cloud instance, the response may differ depending on the type of alert or the elements involved in the alert, such as stopping the instance, restricting the instance's communication, or notifying an administrator. You may also want to retrieve reputation information from external services and adjust the response accordingly. The Action Policy is responsible for defining and controlling these response procedures.

Overall, AlertChain provides a flexible and powerful framework for handling structured event data and determining appropriate actions based on user-defined policies.

### Test

AlertChain is an advanced tool that not only allows you to detect alerts through Alert Policies but also enables you to intentionally execute actions using Action Policies. For more information on how to test these features, please refer to the [Test](./docs/test.md) documentation.

## Usage

To install AlertChain, run the following command:

```bash
$ go install github.com/secmon-lab/alertchain@latest
```

To get started with AlertChain, please refer to the [Getting Started](./docs/getting_started.md) documentation.

Other more documentations is here.

- [Policy](docs/policy.md)
- [Actions](./pkg/action/README.md)
- [Test](docs/test.md)
- [Deployment](docs/deployment.md)
- [Authorization](docs/authz.md)

## Example

In this example, we will demonstrate how AlertChain operates using an event detected by AWS GuardDuty. The policies and data used in this example can be found in the [examples](./examples/basic) directory.

### 1. Write Alert Policy

First, prepare an Alert Policy to detect alerts from the input event data.

**policy/alert.rego**
```rego
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
```

This example alert policy is designed for [AWS GuardDuty](https://docs.aws.amazon.com/cli/latest/reference/guardduty/get-findings.html#examples). The alert evaluates GuardDuty event data based on the following criteria:

- The finding type has a "Trojan:" prefix,
- The severity is greater than 7, and
- If these conditions are met, a new alert is created

Additionally, this policy stores the detected instance's ID as a Attribute, allowing it to be used in a subsequent Action.

### 2. Write Action Policy

Next, prepare an Action Policy. In this example, the action requests a summary and recommended response for the alert from [ChatGPT](https://platform.openai.com/docs/guides/chat), and posts the result to a Slack channel.

**policy/action.rego**
```rego
package action

run contains res if {
	input.alert.source == "aws"
	res := {
		"id": "ask-gpt",
		"uses": "chatgpt.comment_alert",
		"args": {"secret_api_key": input.env.CHATGPT_API_KEY},
	}
}

run contains res if {
	gtp := input.called[_]
	gtp.id == "ask-gpt"

	res := {
		"id": "notify-slack",
		"uses": "slack.post",
		"args": {
			"secret_url": input.env.SLACK_WEBHOOK_URL,
			"channel": "alert",
			"body": gtp.result.choices[0].message.content,
		},
	}
}
```

Action policies are triggered by writing `run` rules. In this case, the first rule is triggered when the `source` of the alert is set to `aws` by the Alert Policy. The `uses` field specifies the Action Name to be executed. The `chatgpt.comment_alert` action requires a `secret_api_key` argument to access ChatGPT via API. The API key is retrieved from the `input.env` environment variables, and the action is executed to make a query to ChatGPT.

The second rule is triggered only if an action with the ID `ask-gpt` has already been executed. The `called` field contains not only information about the executed action but also its result. The result of the query to ChatGPT is retrieved and set as the `body` field, and a message is posted to Slack.

### 3. Run AlertChain as server

After preparing these files, you can start AlertChain using the following command:

```bash
$ alertchain -d policy serve
```

Now, let's create an alert using AWS GuardDuty event data (guardduty.json):

**guardduty.json**
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

To send the event data to the AlertChain API endpoint, use this command:

```bash
$ curl -XPOST http://127.0.0.1:8080/alert/aws_guardduty -d @guardduty.json
```

Upon receiving the data, AlertChain performs the following actions:

1. Evaluates the event data using the alert policy and creates a new alert
2. Evaluates the `action` policy with the new alert, executes `chatgpt.comment_alert`.
3. Evaluate the `action` policy again with not only the alert but also results of executed action, and executes `slack.post` next
4. Evaluate the `action` policy again and no action is triggered. Then stop workflow for the alert

Finally, we can find a Slack message as shown below:

<img width="680" src="https://user-images.githubusercontent.com/605953/236592991-f2411b46-501d-4a4f-9a0d-ff7cf2defc84.png">

## License

Apache License 2.0