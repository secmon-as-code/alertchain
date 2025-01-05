package action

# This is a sample policy for creating an issue on GitHub and adding a comment to it.

# base_github_args is a common argument set for GitHub actions.
# Please see https://github.com/secmon-lab/alertchain/tree/main/action/github for more details.
base_github_args := {
	"app_id": 111111,
	"install_id": 222222,
	"secret_private_key": input.env.GITHUB_APP_PRIVATE_KEY,
	"owner": "your-org",
	"repo": "your-repo",
}

# github_issue_id_key is a key to store the issue ID in the alert.
github_issue_id_key := "github_issue_id"

# has_github_issue checks if the alert has a GitHub issue ID as a persistent attribute
has_github_issue if {
	input.alert.attrs[_].key == github_issue_id_key
}

# If the alert does not have a GitHub issue ID, create a new issue on GitHub.
run contains {
	"id": "create_issue",
	"uses": "github.create_issue",
	"args": base_github_args,
	"commit": [{
		"key": github_issue_id_key,
		"persist": true,
		"path": "number",
	}],
} if {
	input.seq == 0
	not has_github_issue
}

# If GitHub issue ID exists, add a comment to the issue.
run contains {
	"id": "add_comment",
	"uses": "github.add_comment",
	"args": object.union(base_github_args, {
		"issue_number": input.alert.attrs[_].value,
		"body": "This is a comment",
	}),
} if {
	input.seq == 0
	has_github_issue
}
