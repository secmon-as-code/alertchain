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
}
```

`alertchain` imports all environment variables for external variable of jsonnet. They can be loaded by `std.extVar` function of jsonnet. Therefore, a configuration file does not need to have secret values, such as private key, API key, etc.

## policy

The policy object contains the settings for the alert policy.

### path

The path field specifies the path where the policy definition file is stored. In this example, the policy definition file is stored in `tmp/policy` directory.

### package

The package field defines the names of the packages used in the policy. In this example, three packages, alert, enrich, and action, are defined.
