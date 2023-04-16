# AlertChain

AlertChain is a simple SOAR (Security Orchestration, Automation and Response) framework with [OPA](https://github.com/open-policy-agent/opa) (Open Policy Agent).

![](https://user-images.githubusercontent.com/605953/232273906-a3df56fb-3201-4336-b897-e327e8d49981.jpg)

## Motivation

Security Orchestration, Automation, and Response (SOAR) is a platform for automating the detection, analysis, and response to security events. To enable the automated analysis of events and rapid response in SOAR systems, automated security procedures and policies must be executed.

Using OPA and Rego, a SOAR system can flexibly apply a set of user-defined policies to maintain the security of applications and systems. This makes it easier to update or change security policies and enables more accurate policy application. Additionally, the Rego language is flexible and expressive, making it easy to add or modify policies.

## Install

```bash
$ go install github.com/m-mizutani/alertchain@latest
```

See [document](./docs) for configuration and alert/action policy details.

## An Example

A user need to configure 3 files at least.

- configuration file (config.jsonnet)
- alert policy (alert.rego)
- action policy (action.rego)

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

`github-issuer` is a name of action. This action is a creator of GitHub issue as alert ticket by GitHub Apps.

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

The example alert policy is for [AWS GuardDuty](https://docs.aws.amazon.com/cli/latest/reference/guardduty/get-findings.html#examples). The alert evaluates GuardDuty event data as following:

- finding type has "Trojan:" prefix,
- and severity is greater than 7,
- then, creating a new alert

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

The action policy invokes `my_create_github_issue` action defined in the configuration file if alert source is `aws`.

After preparing the files, AlertChain can start with following command.

```bash
$ alertchain -c config.json serve
```

And, let's create alert by AWS GuardDuty event data (guardduty.json).

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

Sending the event data to AlertChain API endpoint.

```bash
$ curl -XPOST http://127.0.0.1:8080/alert/aws_guardduty -d @guardduty.json
```

Then, AlertChain works as following:

1. Evaluates the event data with alert policy, and creates a new alert
2. Evaluates the `action.main` policy with the new alert, and chooses `my_create_github_issue` action.
3. Invokes `my_create_github_issue` that uses `github-issuer` and an GitHub issue will be created.

## License

Apache License 2.0
