# Opsgenie

Actions for [Opsgenie](https://www.atlassian.com/software/opsgenie)

## Prerequisites

- **Opsgenie account**: You need to have an Opsgenie account and the necessary permissions to create alerts.
- **API key**: You need to create an integration API key. You can find management console at `https://<your-opsgenie-tenant>/settings/integrations/`

## `opsgenie.create_alert`

This action creates an alert in Opsgenie.

### Arguments

Example policy:

```rego
run contains job if {
  job := {
    id: "your-action",
    uses: "opsgenie.create_alert",
    args: {
      "secret_api_key": input.env.OPSGENIE_API_KEY,
      "responders": [
        {
          "id": "3f68caf0-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
          "type": "team",
        },
      ],
    }
  }
}
```

- `secret_api_key` (string, required): Specifies the API key for Opsgenie.
- `responders` (array of structure, optional): Specifies the responders of the alert. See https://docs.opsgenie.com/docs/alert-api#section-create-alert for details.
  - `id` (string): Specifies the ID of the responder.
  - `name` (string): Specifies the name of the responder.
  - `username` (string): Specifies the username of the responder.
  - `type` (string, required): Specifies the type of the responder. Possible values are `team`, `user`, `escalation` and `schedule`.
