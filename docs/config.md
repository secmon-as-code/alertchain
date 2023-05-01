# Configuration

The document describes a specification of configuration. A configuration file is written as [jsonnet](https://jsonnet.org/) format. An example of configuration is following.

```jsonnet
{
  policy: {
    path: 'tmp/policy',
    package: {
      alert: 'alert',
      action: 'action',
    },
  },
  actions: [
    {
      id: 'create_github_issue',
      uses: 'github-issuer',
      config: {
        app_id: 134650,
        install_id: 19102538,
        private_key: std.extVar('GITHUB_PRIVATE_KEY'),
        owner: 'm-mizutani',
        repo: 'alert-test',
      },
    },
  ],
}
```

`alertchain` imports all environment variables for external variable of jsonnet. They can be loaded by `std.extVar` function of jsonnet. Therefore, a configuration file does not need to have secret values, such as private key, API key, etc.

## policy

The policy object contains the settings for the alert policy.

### path

The path field specifies the path where the policy definition file is stored. In this example, the policy definition file is stored in `tmp/policy` directory.

### package

The package field defines the names of the packages used in the policy. In this example, three packages, alert, enrich, and action, are defined.

## actions

`actions` define a list of actions to be executed when an alert is triggered.

- `id`: The field specifies a unique identifier for the action. This field must be unique to the other actions and allows the user to specify any string.
- `uses`: The use field specifies the tool or service to be used for the action. In this example, `github-issuer` is used.
- `config`: The config field defines the settings for the action. This field is key-value map and different action by action.

Documents of action for more detail are in [action](./action/README.md).
