# github-issuer

GitHub issuer will create an issue as alert handling ticket.

## Prerequisite

You need to create your GitHub App. An instruction is [here](https://docs.github.com/en/apps/creating-github-apps/creating-github-apps/creating-a-github-app).

The GitHub App requires `Read and Write` permission for `Issues` and you need to install into target repository.

## Config

Example
```jsonnet
    {
      id: 'create-github-issue',
      use: 'github-issuer',
      config: {
        app_id: 134650,
        install_id: 19102538,
        private_key: std.extVar('GITHUB_PRIVATE_KEY'),
        owner: 'm-mizutani',
        repo: 'security-alert',
      },
    },
```

### `app_id` (required)

The app_id field specifies the ID of the GitHub App.

### `install_id` (required)

The install_id field specifies the installation ID of the GitHub account to run the action.

### `private_key` (required)

The private_key field specifies the private key of the GitHub App. In this example, it is expected that the private key will be loaded from an external variable named GITHUB_PRIVATE_KEY.

### `owner` (required)

The owner field specifies the owner name of the GitHub account to run the action.

### `repo` (required)

The repo field specifies the repository name of the GitHub account to run the action.
