package action

run[res] {
    res := {
        "id": "run_mock",
        "uses": "mock",
    }
}

exit[res] {
    input.proc.id == "run_mock"
    input.alert.params[k2].name == "k2"
    res := {
        "params": {
            {
                "id": input.alert.params[k2].id,
                "name": "k2",
                "value": "v2a",
            },
            {
                "name": "k3",
                "value": "v3",
            },
        },
    }
}

run[res] {
    input.called[_].id == "run_mock"
    res := {
        "uses": "mock.after",
    }
}
