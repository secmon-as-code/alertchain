package alert

alert[msg] {
    print(input)
    msg := {
        "title": input.message.finding.category,
        "source": "scc",
    }
}
