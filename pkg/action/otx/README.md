# OTX

Actions for AlienVault OTX (Open Threat Exchange) API.

## `otx.indicator`

This action retrieves information about the specified indicator from the AlienVault OTX API.

### Prerequisite

You need an AlienVault OTX API key. You can sign up for a free account and get an API key [here](https://otx.alienvault.com/).

### Arguments

Example policy:

```rego
run[res] {
  res := {
    id: "your-action",
    uses: "otx.indicator",
    args: {
      "secret_api_key": input.env.OTX_API_KEY,
      "type": "domain",
      "indicator": "example.com",
      "section": "general",
    },
  },
}
```

- `secret_api_key` (string, required): Specifies the AlienVault OTX API key.
- `type` (string, required): Specifies the indicator type. Must be one of "ipv4", "ipv6", "domain", "hostname", "file", or "url".
- `indicator` (string, required): Specifies the indicator value to be queried.
- `section` (string, required): Specifies the section of information to be retrieved. Must be one of "general", "reputation", "geo", "malware", "url_list", "passive_dns", or "http_scans".

### Response

The response is a JSON object containing the requested information about the indicator. The structure of the JSON object depends on the specified section. For more information, see the [AlienVault OTX API documentation](https://otx.alienvault.com/assets/static/external_api.html).

**Example**

```json
{
  "whois": "http://whois.domaintools.com/192.0.2.1",
  "reputation": 0,
  "indicator": "192.0.2.1",
  "type": "IPv4",
  "type_title": "IPv4",
  "base_indicator": {
    "id": 99999999,
    "indicator": "192.0.2.1",
    "type": "IPv4",
    "title": "",
    "description": "",
    "content": "",
    "access_type": "public",
    "access_reason": ""
  },
  "pulse_info": {
    "count": 29,
    "pulses": [
      {
        "id": "xxxxxxxxxxxxxxx",
        "name": "IP Addresses Logged by the Rosethorn PotNet",
        "description": "Malicious activity detections from a small network of honeypots that spans multiple ISPs and geographic locations.\n\nBehavior is logged on ports 23, 80, 3306, and 5900.",
        "modified": "2023-05-13T00:18:43.306000",
        "created": "2023-03-27T16:17:36.094000",
        "tags": [],
        "references": [],
        "public": 1,
        "adversary": "",
        "targeted_countries": [
          "United States of America"
        ],
        "malware_families": [],
        "attack_ids": [],
        "industries": [],
        "TLP": "green",
        "cloned_from": null,
        "export_count": 46,
        "upvotes_count": 0,
        "downvotes_count": 0,
        "votes_count": 0,
        "locked": false,
        "pulse_source": "web",
        "validator_count": 0,
        "comment_count": 0,
        "follower_count": 0,
        "vote": 0,
        "author": {
          "username": "xxxxx",
          "id": "1111111",
          "avatar_url": "/otxapi/users/avatar_image/media/avatars/user_217809/resized/80/avatar_3b9c358f36.png",
          "is_subscribed": false,
          "is_following": false
        },
        "indicator_type_counts": {
          "IPv4": 31038
        },
        "indicator_count": 31038,
        "is_author": false,
        "is_subscribing": null,
        "subscriber_count": 82,
        "modified_text": "4 minutes ago ",
        "is_modified": true,
        "groups": [],
        "in_group": false,
        "threat_hunter_scannable": true,
        "threat_hunter_has_agents": 1,
        "related_indicator_type": "IPv4",
        "related_indicator_is_active": 1
      }
    ],
    "references": [],
    "related": {
      "alienvault": {
        "adversary": [],
        "malware_families": [],
        "industries": []
      },
      "other": {
        "adversary": [],
        "malware_families": [],
        "industries": [
          "Government",
          "Industrial",
          "Defense"
        ]
      }
    }
  },
  "false_positive": [],
  "validation": [],
  "asn": "AS29529 itecom bvba",
  "city_data": true,
  "city": null,
  "region": null,
  "continent_code": "EU",
  "country_code3": "BEL",
  "country_code2": "BE",
  "subdivision": null,
  "latitude": 50.8509,
  "postal_code": null,
  "longitude": 4.3447,
  "accuracy_radius": 100,
  "country_code": "BE",
  "country_name": "Belgium",
  "dma_code": 0,
  "charset": 0,
  "area_code": 0,
  "flag_url": "/assets/images/flags/be.png",
  "flag_title": "Belgium",
  "sections": [
    "general",
    "geo",
    "reputation",
    "url_list",
    "passive_dns",
    "malware",
    "nids_list",
    "http_scans"
  ]
}

```