# Jira

Actions for jira.com

## Prerequisites

- **Jira account**: You need to have a Jira account and the necessary permissions to create and comment on tickets in the specified project.
- **Account ID**: You need to know your account ID. You can find account ID from URL of your profile page. e.g. `https://your-domain.atlassian.net/jira/people/5f6c3a1c1a2b3c4d5e6f7a8b` and `5f6c3a1c1a2b3c4d5e6f7a8b` is your account ID.
- **API token**: You need to create an API token. You can find instructions on how to do so [here](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/).

## `jira.create_issue`

This action creates a ticket in the specified Jira project to serve as an alert handling ticket.

### Arguments

Example policy:

```rego
run[job] {
  job := {
    id: "your-action",
    uses: "jira.create_issue",
    args: {
      "account_id": "5f6c3a1c1a2b3c4d5e6f7a8b",
      "user": "mizutani@hey.com",
      "secret_token": input.env.JIRA_API_TOKEN,
      "base_url": "https://your-domain.atlassian.net",
      "project": "SEC",
      "issue_type": "Task",
      "labels": ["alert"],
    },
  },
}
```

- `account_id` (string, required): Specifies the account ID of the Jira user.
- `user` (string, required): Specifies the Jira user name as email address.
- `secret_token` (string, required): Specifies the Jira API token.
- `base_url` (string, required): Specifies the base URL of the Jira instance. e.g. `https://your-domain.atlassian.net`
- `project` (string, required): Specifies the project key of the Jira project where the ticket will be created.
- `issue_type` (string, required): Specifies the issue type of the Jira project where the ticket will be created. e.g. `Bug`
- `labels` (array of string, optional): Specifies the labels of the ticket.

### Response

See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-post-response

## `jira.add_comment`

This action creates a comment in the specified Jira ticket.

### Arguments

Example policy:

```rego
run[job] {
  job := {
    id: "your-action",
    uses: "jira.add_comment",
    args: {
      "account_id": "5f6c3a1c1a2b3c4d5e6f7a8b",
      "user": "mizutani@hey.com",
      "secret_token": input.env.JIRA_API_TOKEN,
      "base_url": "https://your-domain.atlassian.net",
      "issue_id": "SEC-1",
      "body": "This is a comment",
    },
  },
}
```

- `account_id` (string, required): Specifies the account ID of the Jira user.
- `user` (string, required): Specifies the Jira user name as email address.
- `secret_token` (string, required): Specifies the Jira API token.
- `base_url` (string, required): Specifies the base URL of the Jira instance. e.g. `https://your-domain.atlassian.net`
- `issue_id` (string, required): Specifies the issue ID of the Jira ticket where the comment will be added.
- `body` (string, required): Specifies the body of the comment.

### Response

See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-comments/#api-rest-api-3-issue-issueidorkey-comment-post-response

## `jira.add_attachment`

This action adds an attachment to the specified Jira ticket.

### Arguments

Example policy:

```rego
run[job] {
  job := {
    id: "your-action",
    uses: "jira.add_attachment",
    args: {
      "account_id": "5f6c3a1c1a2b3c4d5e6f7a8b",
      "user": "mizutani@hey.com",
      "secret_token": input.env.JIRA_API_TOKEN,
      "base_url": "https://your-domain.atlassian.net",
      "issue_id": "SEC-1",
      "file_name": "attachment.txt",
      "data": "This is an attachment",
    },
  },
}
```

- `account_id` (string, required): Specifies the account ID of the Jira user.
- `user` (string, required): Specifies the Jira user name as email address.
- `secret_token` (string, required): Specifies the Jira API token.
- `base_url` (string, required): Specifies the base URL of the Jira instance. e.g. `https://your-domain.atlassian.net`
- `issue_id` (string, required): Specifies the issue ID of the Jira ticket where the attachment will be added.
- `file_name` (string, required): Specifies the name of the attachment file.
- `data` (string, required): Specifies the data of the attachment.

### Response

See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-attachments/#api-rest-api-3-issue-issueidorkey-attachments-post-response


