package action

run contains job if {
    input.seq == 0

    job := {
        "uses": "test.output_raw",
        "args": {
            "raw": input.alert.raw,
        }
    }
}
