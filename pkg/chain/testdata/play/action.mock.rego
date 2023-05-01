package action.my_action

action[res] {
    p := input.alert.params[_]
    p.key == "c"
    p.value < 2

    res := {
        "id": "my_action",
        "params": {
            {
                "key": "c",
                "value": p.value + 1,
            },
        }
    }
}
