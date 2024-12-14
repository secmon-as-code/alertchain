# Instructions

The initial JSON data provided contains information about false positive alerts. Based on the code given thereafter, generate a new Rego policy file to ignore these alerts.

# Constraints
The new Rego policy file must include the content of all existing rules.
Integrate rules if possible.
The output should be in Rego code format only, not Markdown.
Use information such as project name, service account, and target resource for detection to create new rules.
Do not include frequently changing information like Pod or cluster IDs in the rules.
