# Getting Started

This is a guide to help you get started with AlertChain.

## Prerequisites

### Install tools

- [AlertChain](../README.md#usage)
- [OPA](https://www.openpolicyagent.org/docs/latest/#running-opa)
- Container runtime. E.g. [Docker](https://docs.docker.com/get-docker/)

### Integration services
AlertChain requires alert event sender. Currently, AlertChain supports only HTTP endpoint as an alert event receiver. You can use the following services to send alert events to AlertChain:

- [Amazon SNS](https://aws.amazon.com/sns) with e.g. [GuardDuty](https://docs.aws.amazon.com/guardduty/latest/ug/guardduty_sns.html)
- [Google Cloud Pub/Sub](https://cloud.google.com/pubsub/docs/overview) with e.g. [Cloud Security Command Center](https://cloud.google.com/security-command-center/docs/how-to-notifications)
- [GitHub Webhook](https://docs.github.com/en/developers/webhooks-and-events/webhooks)
- [CrowdStrike Falcon Webhook](https://marketplace.crowdstrike.com/listings/webhook)

## Create a new policy repository

You can create a new policy repository by `alertchain new` command.

```shell
$ alertchain new
11:19:40.442 INFO Copy file path=".gitignore"
11:19:40.442 INFO Copy file path="Dockerfile"
11:19:40.443 INFO Copy file path="Makefile"
11:19:40.443 INFO Copy file path="policy/action/main.rego"
11:19:40.443 INFO Copy file path="policy/alert/main.rego"
11:19:40.443 INFO Copy file path="policy/alert/main_test.rego"
11:19:40.444 INFO Copy file path="policy/alert/testdata/your_schema/event.json"
11:19:40.444 INFO Copy file path="policy/authz/http.rego"
11:19:40.444 INFO Copy file path="policy/play/test.rego"
11:19:40.444 INFO Copy file path="scenario/data/event.json"
11:19:40.444 INFO Copy file path="scenario/env.libsonnet"
11:19:40.444 INFO Copy file path="scenario/my_first_scenario.jsonnet"
```

This command creates directories and new sample files.

## Customize the setting files and policies

### Alert Policy

Alert policies is for determination if the input event is acceptable alert or not. You can customize the alert policy by editing [policy/alert/main.rego](../pkg/usecase/templates/policy/alert/main.rego).

```rego
package alert.your_schema

alert contains {
	"title": input.name,
	"description": "Your description here",
	"source": "your_source",
	"namespace": input.key,
} if {
	input.severity == ["HIGH", "CRITICAL"][_]
}
```

This is an example of alert policy. It assumes the event data is like following:

```json
{
  "name": "my_event",
  "key": "value",
  "severity": "HIGH"
}
```

- `your_schema` is for identification of alert data schema. When AlertChain receives event via `/alert/raw/your_schema` (raw event) or `/alert/pubsub/your_schema` (Google Cloud Pub/Sub), the policy is triggered.
- `input` is a structured data that is same with the input event data.
- `alert` is a rule to determine if the input event is acceptable alert or not. If the `alert` "contains"

Please see [Alert Policy document](./policy.md#alert-policy) for more details.

### Action Policy

Action policy is a rule of workflow for detected alerts. You can customize the action policy by editing [policy/action/main.rego](../pkg/usecase/templates/policy/action/main.rego).

The sample action policy is for creating an issue in GitHub repository.

1. Create a GitHub App and install it to the target repository.
2. Set the GitHub App ID, installation ID, owner, repository name, and private key to the action policy.

You can see more details in [Action Policy document](./policy.md#action-policy).

## End-to-end test with `play` command

You can test the alert and action policies with the `play` command. `play` command simulates the alert event and executes the action policy and output the result into JSON files.

```bash
% alertchain play -d ./policy -s ./scenario -o ./policy/play/output
13:13:15.141 INFO loading policy package="alert" path="./policy"
13:13:15.143 INFO loading policy package="action" path="./policy"
13:13:15.146 INFO starting alertchain with play mode scenario dir="./scenario" output dir="./policy/play/output" targets=[]
13:13:15.147 INFO Done scenario id=my_first_scenario
% tree ./policy/play/output/
./policy/play/output/
├── my_first_scenario
│   └── data.json
└── result.json

2 directories, 2 files
```

Please see [test document](./test.md#testing-action-policy) for more details.

By the way, the default template has `Makefile` to test both of the policies.
