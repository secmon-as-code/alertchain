package alert.my_alert

alert[msg] {
    input.color == "blue"

    msg := {
        "title": "Test alert",
        "namespace": "test_alert",
    }
}
