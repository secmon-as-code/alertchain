# Authorization

## Introduction

AlertChain receives alert data from various sources. For example, you can receive alerts from AWS GuardDuty, Google Cloud Security Command Center, or your own SIEM. AlertChain can also receive alerts from multiple sources at the same time.

AlertChain receives alerts via HTTP API. Then, AlertChain requires to expose HTTP port to Internet for public SaaS. Therefore, AlertChain needs to authenticate and authorization mechanism for the sender of the alert data to prevent unauthorized data input.

The authentication and authorization mechanism of AlertChain is based on [Open Policy Agent](https://www.openpolicyagent.org/). Open Policy Agent is a policy engine that can be used to implement fine-grained access control. AlertChain uses Open Policy Agent to implement authentication and authorization for alert data.

## Authorization Policy

The authorization policy is written in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). Here is HTTP authorization policy example:

```rego
package authz.http

default deny = false

deny {
    not net.cidr_contains("10.0.0.0/8", input.remote)
}
```

This policy allows access from the `10.0.0.0/8`` network.

## Policy Specification

### Input

The input of the authorization policy is as follows:

- `input.remote` (string): IP address of the sender of the alert data
- `input.method` (string): HTTP method of the request
- `input.path` (string): HTTP path of the request
- `input.query` (map of string array): HTTP query of the request
- `input.header` (map of string array): HTTP headers of the request

### Output

The output of the authorization policy is as follows:

- `deny` (boolean): Deny access if `true` is returned. `false` and undefined are treated as allow.

When `deny` is `true`, HTTP response is as follows:

- Status code: 403
- Message: `Access denied`

## Examples

### Validate Google Cloud Service

```rego
package authz.http

default deny = true

deny := false { allow }

jwks_request(url) := http.send({
    "url": url,
    "method": "GET",
    "force_cache": true,
    "force_cache_duration_seconds": 3600 # Cache response for an hour
}).raw_body

allow {
    startswith(input.path, "/alert/")

    ahthHdr := input.header["Authorization"]
    count(ahthHdr) == 1
    authHdrValues := split(ahthHdr[0], " ")
    count(authHdrValues) == 2
    lower(authHdrValues[0]) == "bearer"
    token := authHdrValues[1]

    jwks := jwks_request("https://www.googleapis.com/oauth2/v3/certs")

    io.jwt.verify_rs256(token, jwks)
    claims := io.jwt.decode(token)
    claims[1]["email"] == "xxxxxx-compute@developer.gserviceaccount.com"
}
```
