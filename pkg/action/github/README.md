# GitHub

Actions for github.com

## `github.create_issue`

The action will create an issue as alert handling ticket.

### Prerequisite

You need to create your GitHub App. An instruction is [here](https://docs.github.com/en/apps/creating-github-apps/creating-github-apps/creating-a-github-app).

The GitHub App requires `Read and Write` permission for `Issues` and you need to install into target repository.

### Arguments

Example policy

```rego
run[res] {
  res := {
    id: "your-action",
    uses: "github.create_issue",
    config: {
      app_id: 134650,
      install_id: 19102538,
      owner: "m-mizutani",
      repo: "security-alert",
      secret_private_key: input.env.GITHUB_PRIVATE_KEY,
    },
  },
}
```

- `app_id` (number, required): The app_id field specifies the ID of the GitHub App.
- `install_id` (number, required): The install_id field specifies the installation ID of the GitHub account to run the action.
- `owner` (string, required): The owner field specifies the owner name of the GitHub account to run the action.
- `repo` (string, required): The repo field specifies the repository name of the GitHub account to run the action.
- `secret_private_key` (string, required): The private_key field specifies the private key of the GitHub App.
- `assignee` (string, optional): Specify the GitHub user to be assigned.
- `labels` (array of strings, optional): Specify the labels.

Note: If you wish to use `assignee` or `labels`, the GitHub App must also have `Read and Write` permissions for `Content`.