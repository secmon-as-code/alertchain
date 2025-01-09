package github

import (
	"context"
	_ "embed"

	"net/http"

	"log/slog"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gots/ptr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

func CreateComment(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	// Required arguments
	appID, ok := args["app_id"].(float64)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "app_id is required")
	}

	installID, ok := args["install_id"].(float64)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "install_id is required")
	}

	privateKey, ok := args["secret_private_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "private_key is required")
	} else if !isRSAPrivateKey(privateKey) {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "private_key must be RSA private key")
	}

	owner, ok := args["owner"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "owner is required")
	}

	repo, ok := args["repo"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "repo is required")
	}

	issue_number, ok := args["issue_number"].(float64)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "issue_number is required")
	}

	body, ok := args["body"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "body is required")
	}

	if ctxutil.IsDryRun(ctx) {
		return nil, nil
	}

	req := &github.IssueComment{
		Body: &body,
	}

	rt := http.DefaultTransport

	itr, err := ghinstallation.New(rt, int64(appID), int64(installID), []byte(privateKey))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create GitHub App installation transport")
	}

	client := github.NewClient(&http.Client{Transport: itr})

	comment, resp, err := client.Issues.CreateComment(ctx, owner, repo, int(issue_number), req)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create GitHub comment")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, goerr.New("Failed to create GitHub comment (unexpected status code)", goerr.V("status", resp.StatusCode))
	}

	ctxutil.Logger(ctx).Debug("Created GitHub comment",
		slog.Any("comment_id", ptr.From(comment.ID)),
	)

	return comment, nil
}
