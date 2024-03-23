# incident.io

The `incident.io` action is used to create incidents in the incident.io platform.

## Pre-requisites

- Create an account in the incident.io platform
- `Team` or higher plan (for API access)
- Create Alert source and Alert route. See [the article](https://incident.io/changelog/automatically-create-incidents-with-alerts) for more information.
  - **API Token**: Generate an API token in the incident.io platform
  - **Alert Source Config ID**: Get the Alert Source Config ID from the Alert Source settings

## `incident_io.create_alert`

This action creates an alert in the incident.io platform. It calls [CreateHTTP Alert Event V2](https://api-docs.incident.io/tag/Alert-Events-V2#operation/Alert%20Events%20V2_CreateHTTP) API.

### Arguments

Example policy:

```rego
run[job] {
  job := {
    id: "your-action",
    uses: "incident_io.create_alert",
    args: {
      "secret_api_token": input.env.INCIDENT_IO_API_TOKEN,
      "alert_source_config_id": "B2BJGP4XSC4ZNWUY8FB5KIOV4SLGRS8N",
      "title": "Alert title",
      "description": "Alert description",
      "status": "firing",
      "deduplication_key": "dedup-key",
      "metadata": {
        "key1": "value1",
        "key2": "value2",
      },
    },
  },
}
```

- `secret_api_token` (string, required): Specifies the API token of the incident.io platform.
- `alert_source_config_id` (string, required): Specifies the Alert Source Config ID.
- `title` (string, optional): Specifies the title of the alert. Default is `title` of [Alert](../../../docs/policy.md#alert).
- `description` (string, optional): Specifies the description of the alert. Default is `description` of [Alert](../../../docs/policy.md#alert).
- `status` (`firing` or `resolved`, optional): Specifies the status of the alert. Default is `firing`.
- `deduplication_key` (string, optional): Specifies the deduplication key of the alert. Default is `id` of [Alert](../../../docs/policy.md#alert).
- `metadata` (object, optional): Specifies the metadata of the alert. Default is built from `attrs` of [Alert](../../../docs/policy.md#alert) and provided `metadata` will be merged into default metadata.
- `source_url`: A link to the alert in the upstream system

### Response

See the response part of [CreateHTTP Alert Event V2](https://api-docs.incident.io/tag/Alert-Events-V2#operation/Alert%20Events%20V2_CreateHTTP).
