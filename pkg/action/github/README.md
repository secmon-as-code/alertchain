# GitHub

Actions for github.com

## `github.create_issue`

This action creates an issue in the specified GitHub repository to serve as an alert handling ticket.

### Prerequisite

You need to create a GitHub App. You can find instructions on how to do so [here](https://docs.github.com/en/apps/creating-github-apps/creating-github-apps/creating-a-github-app).

The GitHub App requires `Read and Write` permissions for `Issues`, and you need to install it into the target repository.

### Arguments

Example policy:

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

- `app_id` (number, required): Specifies the ID of the GitHub App.
- `install_id` (number, required): Specifies the installation ID of the GitHub account where the action will be executed.
- `owner` (string, required): Specifies the owner name of the GitHub account where the action will be executed.
- `repo` (string, required): Specifies the repository name of the GitHub account where the action will be executed.
- `secret_private_key` (string, required): Specifies the private key of the GitHub App.
- `assignee` (string, optional): Specifies the GitHub user to be assigned to the issue.
- `labels` (array of strings, optional): Specifies the labels to be applied to the issue.

Note: If you wish to use `assignee` or `labels`, the GitHub App must also have `Read and Write` permissions for `Content`.