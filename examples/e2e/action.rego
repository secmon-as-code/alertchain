package action

count_attr := input.alert.attrs[x] {
    input.alert.attrs[x].key == "count"
} else := {
    "id": null,
    "value": 0,
}

init[res] {
    input.seq == 0

    res := {
        "abort": count_attr.value < 2,
        "attrs": [{
            "id": count_attr.id,
            "key": "count",
            "value": count_attr.value + 1,
            "global": true,
        }]
    }
    print(res)
}

run[job] {
    job := {
        "id": "testing",
        "uses": "http.fetch",
        "args": {
            "method": "GET",
            "url": "http://localhost:9876/test",
        },
    }
}
