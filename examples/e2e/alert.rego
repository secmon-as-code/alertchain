package alert.my_alert

alert contains msg if {
    input.color == "blue"

    msg := {
        "title": "Test alert",
        "namespace": "test_alert",
    }
}
