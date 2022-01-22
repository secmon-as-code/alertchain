# AlertChain

AlertChain is a simple SOAR (Security Orchestration, Automation and Response) framework  with [OPA](https://github.com/open-policy-agent/opa) (Open Policy Agent).

## Concept

SOAR is a platform to automate security alert handling (investigating alert, evaluation risk, remediation, etc). There are a lot of existing SOAR product and most of them provide rich GUI configuration feature. That is familiar with beginners, however complexity is increased according to rule/policy size and become difficult to change rule/policy. A concept of [Policy as Code](https://docs.hashicorp.com/sentinel/concepts/policy-as-code) helps the problem with test automation, deploy automation, version control with review and so on.

[OPA](https://github.com/open-policy-agent/opa) (Open Policy Agent) is a generic policy engine. OPA is open technology and can be also useful to manage SOAR rule/policy.

![AlertChain architecture overview](https://user-images.githubusercontent.com/605953/147866269-48a0df6f-181d-4fc1-ac90-3b1650fd0dff.jpg)

AlertChain is a simple SOAR framework that is fully integrated with OPA. AlertChain has 2 type of components.

- `source`: A security alert receiver. AlertChain does not generate/detect security alert by own. AlertChain can import various security alerts from `source` instead.
- `action`: A security alert handling workflow. There are 2 major categories of actions.
    - **Enrich**: Actions to append additional information to a security alert. A security alert has multiple attributes. E.g. Source IP address of an attacker, Cloud instance ID as victim. Internal (e.g. cloud service used by own) and external (e.g. open threat intelligence database) data source might have related information of the attribute. For example, an information if source IP address is in blacklist helps to determine risk of the security alert. An enrich workflow inquiry to internal data sources and adds the result(s) to a security alert. The subsequent workflows will refer to it.
    - **Response**: Actions to handle security alert.

A set of actions is called `job`. In the figure, actions are in either one of Enrich job and Response job. Actions in a job are executed in parallel. When all actions are completed, security alert is passed to a next job.

### How OPA works in AlertChain

- For source:
    - **Update alert**: Update parameters of security alert, e.g. severity, status and meta data. AlertChain queries to OPA server or local policy and update parameters with returned results.
- For action (job)
    - **Update alert**: Same with one for source
    - **Control actions**: Decide which action should be executed in a next job. Especially actions in response category might have affect to production environment. Therefore, not every action should be executed in every case, and the decisions need to be complex and precise. A user can write policy to choose invoking actions with attributes and parameters of the security alert.

A following figure describes more detail of relationship with source, job and OPA server (or local policy).

![Workflow with OPA inquiry](https://user-images.githubusercontent.com/605953/147866890-3dd8b613-e249-4fce-9121-ad2581460524.jpg)

Between source/job and job, both of "Update alert" and "Control actions" procedures are happened. After last job, only "Update alert" is happened. This architecture achieves followings:

- **Avoid massive inquiry to high cost data source** in Enrich job(s). Some data source requires cost per inquiry. For example, search query from massive data (e.g. BigQuery), or external paid data source (e.g. threat information database by vendor). In an other case, some data API service has rate limit. Therefore, inquiry traffic control is needed.
- **Achieves Precise remediation action** to avoid destroy of important resource by false positive alert. Response actions includes destructive procedure, e.g. terminate exploited instance or quarantine affected endpoint. These actions are valuable for true positive alert, however not good for false positive alert. Therefore, remediate action is required more tighter controls. Users can create policies to control remediation action more carefully.

## Usage

### Getting Started

Create a configuration file.

```jsonnet:config.jsonnet
{
    "policy": {
        "type": "local",
        "path": "alert.rego",
    },
    "database: {
        "type": "postgres",
        "config": "host=127.0.0.1 port=5432 user=postgres dbname=alertchain password=xxxxxxxx",
    },
    "sources": {
        "api": {

        }
    },
    "actions": [
        {
            "id": "vt",
            "use": "inspect-virustotal"
            "config": {
                "api_key": "xxxxxx",
            },
        },
        {
            "id": "notify-slack",
            "use": "notify-slack"
            "config": {
                "webhook_url": "https://slack.com/xxxxxxxx",
            },
        },
    ],
    "jobs": [
        {
            "name": "enrich",
            "actions": ["vt"],
        },
        {
            "name": "response",
            "actions": ["notify-slack"],
        },
    ],
}
```

And run alertchain.

```bash
$ alertchain -c config.jsonnet
```

## How to write policy

### For update alert

#### Specification

- package name: `alertchain.alert`
- input:
    - `input.alert`: Security alert in handling
- output:
    - `severity` (string): e.g. `low`, `medium`, `high`
    - `status` (string): e.g. `new`, `closed`

If `severity` and/or `status` is returned, update them.

#### Example

```rego
package alertchain.alert

# m-mizutani is always detected as suspicious user.
status = result {
    input.alert.title == "suspicious activity"
    attr := input.alert.attributes[_]
    attr.key == "account"
    attr.value == "m-mizutani"

    result := "closed"
    # Then, alert status will be updated to "closed"
}
```

### For control action

#### Specification

- package name: `alertchain.action`
- input:
    - `input.alert`: Security alert in handling
    - `input.job`: Job information
        - `input.job.name`: Job name
    - `input.action`:
        - `input.action.name`: Action name
- output:
    - `cancel` (boolean): If true, action will be cancelled
    - `args` (array of `alert.attributes[_]`): Target attribute(s) of action

#### Example

```rego
package alertchain.action

default execute = true

execute = result {
    input.job.name == "enrich"
    input.alert.title == "misconfiguration in cloud service"

    # This alert is just human error, no additional investigation required
    result := false
}
```

### Alert structure

```json
{
    "title": "Alert title",
    "detector": "who issues the alert",
    "description": "alert descriptions",
    "attributes": [
        {
            "key": "name of attribute, e.g. `src-ip-addr`",
            "value": "value of attribute, e.g. `192.168.0.1`",
            "context": ["context of the value", "e.g. `local` or `remote`"],
            "type": "type of attribute, e.g. `ipaddr`",
            "annotations": [
                {

                }
            ],
        }
    ]
}
```

#### Enumeration

- `attribute.type`: Choose one from followings
    - `ipaddr`: IP address (IPv4, IPv6)
    - `port`: TCP/UDP port number
    - `domain`: Domain name
    - `user_id`: User ID
    - `email`: Email
    - `sha256`: SHA 256 hash value
    - `filepath`: File path
    - `url`: URL
- `attribute.context`: Choose one ore more from followings
    - `local`: Means local or internal environment
    - `remote`: Means remote or external environment
    - `client`: Means source of the activity
    - `server`: Means target of the activity

## License

MIT License
