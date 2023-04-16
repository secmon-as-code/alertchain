package action.mock

action[res] {
    p := input.alert.params[_]
    p.key == "c"
    p.value < 10

    res := {
        "id": "mock",
        "params": {
            {
                "key": "c",
                "value": p.value + 1,
            },
        }
    }
}
