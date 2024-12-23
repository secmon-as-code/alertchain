package action

count_attr := input.alert.attrs[x] if {
    input.alert.attrs[x].key == "count"
} else := {
    "id": null,
    "value": 0,
}

run contains res if {
    res := {
        "commit": [{
            "id": count_attr.id,
            "key": "count",
            "value": count_attr.value + 1,
            "persist": true,
        }],
    }
}

run contains job if {
    job := {
        "id": "test",
        "uses": "http.fetch",
        "args": {
            "method": "GET",
            "url": "http://localhost:9876/test",
        },
    }
}
