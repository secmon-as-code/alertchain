# Policy

AlertChain has two types of policies: "Alert Policy" and "Action Policy". Both policies are written in the Rego language. This document describes the input and output schemas for these policies. See [Open Policy Agent official document](https://www.openpolicyagent.org/docs/latest/policy-language/) for more detail of Rego language.

The **Alert Policy** is responsible for determining whether the incoming event data from external sources should be treated as an alert or not. For example, when receiving notifications from external services, you may want to handle only alerts related to specific categories, or you may want to exclude events that meet certain conditions (such as specific users or hosts). The Alert Policy can be used to achieve these goals by excluding certain events or including only specific events as alerts.

On the other hand, the **Action Policy** determines the appropriate response for detected alerts. For example, when an issue is detected on a cloud instance, the response may differ depending on the type of alert or the elements involved in the alert, such as stopping the instance, restricting the instance's communication, or notifying an administrator. You may also want to retrieve reputation information from external services and adjust the response accordingly. The Action Policy is responsible for defining and controlling these response procedures.

## Alert Policy

### Package

The package name for Alert Policy must follow the naming convention below:

```rego
package alert.{schema}
```

Here, `{schema}` must match the `{schema}` specified when receiving event data. For example, if the endpoint path for receiving data via Pub/Sub is `/alert/pubsub/my_alert`, the policy `package alert.my_alert` will be called.

### Input

The input for Alert Policy will be the structured data (mainly JSON) received from the previous phase. For example, if the following message is input via Google Cloud Pub/Sub:

```
{
    "message": {
        "data": "eyJuYW1lIjoiaG9nZSJ9Cg==",
    },
}
```

From the Pub/Sub schema, `message.data` is extracted, and `eyJuYW1lIjoiaG9nZSJ9Cg==` is Base64 decoded to:

```json
{
    "name": "hoge"
}
```

This data is stored in Rego's `input`. The policy will determine whether this data will be treated as an alert or not based on this data.

### Output

Once the alert determination is made, store the data with the schema below in the `alert` rule. The stored data will be treated as an alert. Output schema is according to Alert structure.

The `attrs` field (Attribute) serves not only to extract event data fields but also to accommodate user-defined values. For instance, users can add their own `severity` key Attribute to determine the appropriate action. Attributes bind the alert and can be added or replaced by the action policy. (Refer to the Action Policy section for more details)

### Example

```rego
package alert.my_alert

alert[res] {
    input.name == "hoge"
    res := {
        "title": "detected hoge",
        "attrs": [
            {
                "key": "color",
                "value": "blue",
            },
        ],
    }
}
```

In this example, the policy checks if the input contains the name "hoge". If it does, an alert will be created with the title "detected hoge" and a Attribute of "color" set to "blue".

## Action Policy

An Action Policy is responsible for defining the following:

- `run` rule
  - Choose the next action to execute
  - Provide arguments for the next action
- `exit` rule
  - Create new or replace Attributes for the next action
  - Abort all action processes if necessary

The relationship between the `run` and `exit` rules in the Action Policy and the execution order of actions is illustrated in the diagram below.

![AlertChain - Frame 1](https://user-images.githubusercontent.com/605953/236360762-af2675db-9adc-47a0-bf00-030196e8ec9a.jpg)

When an alert is detected by the Alert Policy, the `run` and `exit` rules within the Action Policy are called alternately. The `run` rule can specify the execution of multiple actions. Also, the `exit` rule is called after each action is completed. If no actions are selected by the `run` rule, all operations will terminate.

```rego
package action

run[res] {
    input.seq == 0
    res := {
        "uses": "chatgpt.query",
        "args": {
            "secret_api_key": input.env.CHATGPT_API_KEY,
        },
    }
}
```

In this example, an action called `chatgpt.query` is launched to query the alert content to ChatGPT. The action to be launched is specified by `uses`, and the required arguments are specified by `args`. The `input.seq` value increments by 1 each time the `run` rule is called. Therefore, when the `run` rule is called for the second time, the result of `input.seq == 0` will be false, making the rule invalid, and no subsequent actions will be specified. If no actions are specified, the entire process will stop.

The `exit` rule primarily handles the transfer of Attributes obtained from the results of actions. Let's modify the ChatGPT calling rule slightly.

```rego
package action

run[res] {
    input.seq == 0
    res := {
        "id": "ask-chatgpt",
        "uses": "chatgpt.query",
        "args": {
            "secret_api_key": input.env.CHATGPT_API_KEY,
        },
    }
}

exit[res] {
    input.action.id == "ask-chatgpt"
    res := {
        "attrs": [
            {
                "name": "ChatGPT's comment",
                "value": input.action.result.choices[0].message.content,
            }
        ]
    }
}
```

We added the `id` value `ask-chatgpt` to the `run` rule, and then checked for it in the `exit` rule with `input.action.id == "ask-chatgpt"`. This ensures that the `exit` rule is only valid after the first `run` rule has been executed. In this `exit` rule, we extract the response message stored in the result of the ChatGPT action (https://platform.openai.com/docs/guides/chat/response-format) and store it as a Attribute. The stored Attribute will then be available for use in subsequent `run` and `exit` rules.

After the `exit` rule is called, the `run` rule is called again. For example, by adding the following `run` rule, we can ensure that the `run` rule is only valid after the `ask-chatgpt` action has been executed.

```
package action

run[res] {
    input.called[_].id == "ask-chatgpt"
    res := {
        "id": "notify-slack",
        "uses": "slack.post",
        "args": {
            "secret_webhook_url": input.env.SLACK_INCOMING_WEBHOOK,
        },
    }
}
```

### `run` rule

#### Input

An Action Policy accepts the following input:

- `input.alert`: [Alert](#alert)
- `input.env`: Map of (string, string): Map of environment variables of the AlertChain process.
- `input.seq` (number): Sequence number of actions, starting from 0.
- `input.called`: Array of [Action](#action): Actions that have already been called.

Using this input, the action policy can process the alert data and determine the most appropriate action to perform next, along with the necessary arguments and Attributes.

#### Output

After evaluating the action policy, if the next action is required, set the `action` field according to the schema of [Action](#action):

### `exit` rule

#### Input

- `input.alert`: [Alert](#alert)
- `input.seq` (number): Sequence number of actions, starting from 0.
- `input.called`: Array of [Action](#action): Actions that have already been called.
- `input.action`: [Action](#action). The previously executed action is stored here. The response from that action is included in `input.action.result`.

#### Output

- `attrs`: Array of [Attribute](#Attribute).

## Basic data structures

### Alert

- `title` (string, required): Title of the alert
- `description` (string, optional): Human-readable explanation about the alert
- `source` (string, optional): Data source
- `attrs` (array, optional): Array of [Attribute](#Attribute)

### Attribute

- `name` (string, required): Name of the Attribute
- `value` (any, required): Value of the Attribute
- `id` (string, optional): ID of the Attribute. If not set, it will be assigned automatically.
- `type` (string, optional): Type of the Attribute

In a single alert, the `name` of a Attribute can be duplicated, but the `id` must be unique. If duplicate `id`s are passed, the later-specified Attribute will overwrite the earlier one. Keep in mind that the execution order of actions within the same sequence is not guaranteed, so be careful of duplication when specifying IDs. If you want to modify a Attribute, you can intentionally overwrite it by specifying its ID.

### Action

- `id` (string, optional): A unique ID for the action within the alert. If not specified, it will be assigned automatically. An ID should only be executed once, so do not specify an ID for actions that need to be executed multiple times. Conversely, by explicitly specifying an `id`, you can prevent an action from being executed multiple times.
- `uses` (string, required): Specify the name of the action to be launched.
- `args`: Specify the arguments for each action in a key-value format.
- `result`: When called in the `exit` rule, the result of the action is stored.

NOTE: Arguments with the `secret_` prefix in `args` have a special meaning. This indicates that the value is confidential (e.g., API keys) and will not be output in logs or similar records.