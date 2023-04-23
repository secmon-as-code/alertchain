# AlertChain

AlertChain is a simple SOAR (Security Orchestration, Automation, and Response) framework that leverages [OPA](https://github.com/open-policy-agent/opa) (Open Policy Agent) to enhance security management.

![AlertChain Diagram](https://user-images.githubusercontent.com/605953/232273906-a3df56fb-3201-4336-b897-e327e8d49981.jpg)

## Motivation

Security Orchestration, Automation, and Response (SOAR) is a platform designed for automating the detection, analysis, and response of security events. In order to enable automated event analysis and rapid response in SOAR systems, it is essential to execute automated security procedures and policies.

By utilizing OPA and Rego, a SOAR system can flexibly apply a set of user-defined policies to maintain the security of applications and systems. This approach simplifies the process of updating or changing security policies and ensures a more accurate policy application. Moreover, the Rego language is flexible and expressive, making it easy to add or modify policies.

## Concept

AlertChain is a versatile software that accepts structured event data through HTTP or other means, and then determines its actions based on policies written in Rego.

### Action

Actions, the basic units of operation, are primarily implemented within AlertChain using Go code. For example, there is an action called [github-issuer](docs/action/github-issuer.md) which creates an issue on GitHub. Users can define any number of actions in a configuration written in Jsonnet, each of which needs a unique ID. Basic settings, such as the username and API key required to execute the action, are defined within [the Jsonnet configuration](docs/config.md). Additionally, runtime adjustments, such as specifying labels when creating a GitHub issue, can be made within the policy itself.

### Policy

There are two main types of policies in AlertChain:

1. Alert Policy: This policy evaluates the input structured data and determines whether the event should trigger an alert. Users can add any parameters as metadata during the evaluation process. The policy is invoked only once.

2. Action Policy: This policy decides how to respond if an event is determined to be an alert. The `action.main` policy is always invoked first to determine the initial action to be taken. Actions are specified by their unique ID defined in the Jsonnet configuration, and if the action accepts arguments, they can be specified by the policy. After an action is executed, the `action.<action ID>` policy is invoked, allowing users to specify further actions or end the process. Parameters can be added or overwritten for calling new actions, allowing for alert state maintenance and enabling conditional branching and repetition if necessary.

Overall, AlertChain provides a flexible and powerful framework for handling structured event data and determining appropriate actions based on user-defined policies.

## Installation

To install AlertChain, run the following command:

```bash
$ go install github.com/m-mizutani/alertchain@latest
```

Refer to the [documentation](./docs) for details on configuration and alert/action policies.

## Example

To configure AlertChain, users need to create at least three files:

1. Configuration file (config.jsonnet)
2. Alert policy (alert.rego)
3. Action policy (action.rego)

**config.jsonnet**
```jsonnet
{
  policy: {
    path: './policy',
  },
  actions: [
    {
      id: 'my_create_github_issue',
      uses: 'github-issuer',
      config: {
        app_id: 12345,
        install_id: 67890,
        private_key: std.extVar('GITHUB_PRIVATE_KEY'),
        owner: 'm-mizutani',
        repo: 'security-alert',
      },
    },
  ],
}
```

`github-issuer` is an action that creates a GitHub issue as an alert ticket using GitHub Apps.

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

This example alert policy is designed for [AWS GuardDuty](https://docs.aws.amazon.com/cli/latest/reference/guardduty/get-findings.html#examples). The alert evaluates GuardDuty event data based on the following criteria:

- The finding type has a "Trojan:" prefix,
- The severity is greater than 7, and
- If these conditions are met, a new alert is created

**policy/action.rego**
```rego
package action.main

action[res] {
    input.alert.source == "aws"
    res := {
        "id": "my_create_github_issue",
    }
}
```

The action policy triggers the `my_create_github_issue` action defined in the configuration file if the alert source is from `aws`.

After preparing these files, you can start AlertChain using the following command:

```bash
$ alertchain -c config.json serve
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
2. Evaluates the `action.main` policy with the new alert, selecting the `my_create_github_issue` action
3. Executes the `my_create_github_issue` action that uses `github-issuer` to create a new GitHub issue

## License

Apache License 2.0