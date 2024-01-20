package action

run[job] {
    input.seq == 0

    job := {
        "uses": "test.output_raw",
        "args": {
            "raw": input.alert.raw,
        }
    }
}
