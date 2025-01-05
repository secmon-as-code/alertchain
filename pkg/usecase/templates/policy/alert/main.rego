package alert.your_schema

alert contains {
	"title": input.name,
	"description": "Your description here",
	"source": "your_source",
	"namespace": input.key,
} if {
	input.severity == ["HIGH", "CRITICAL"][_]
}
