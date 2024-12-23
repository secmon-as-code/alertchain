# Instructions

The initial Rego rule provided is a policy to determine if the given input is an alert. The following JSON-formatted data represents false positives. Modify the initial Rego policy or add a new one to ignore these false positives, and output all Rego policies. If adding a new policy, ensure it aligns with the existing one. Please adhere to Rego syntax when writing.

# Rule

## Input

`input` (object): The input object. This is the received JSON data.

## Output

`alert` (set): The set containing the alert data. If the input is an alert (that should be triaged as security issue), the set should contain the alert data. Otherwise, it should be empty.

# Constraints

The new Rego policy file must include the content of all existing rules.
Integrate rules if possible.
The output should be in Rego code format only, not Markdown.
Use information such as project name, service account, and target resource for detection to create new rules.
Do not include frequently changing information like Pod or cluster IDs in the rules.
Use tab indentation for the rules instead of spaces.
Output only the Rego code. Do not include any other characters, spaces, or quotes like markdown in the output.
