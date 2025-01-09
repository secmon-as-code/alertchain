package gemini

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
	"github.com/m-mizutani/goerr/v2"
)

type Client struct {
	client *genai.Client
	model  string
}

func New(ctx context.Context, projectID, location string) (*Client, error) {
	// modelName := "gemini-1.5-flash-002"
	modelName := "gemini-2.0-flash-exp"

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to create genai client")
	}
	return &Client{client: client, model: modelName}, nil
}

func (x *Client) Generate(ctx context.Context, prompts ...string) ([]string, error) {
	gemini := x.client.GenerativeModel(x.model)

	var parts []genai.Part
	for _, t := range prompts {
		parts = append(parts, genai.Text(t))
	}

	resp, err := gemini.GenerateContent(ctx, parts...)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to generate content")
	}

	var respText []string
	for _, c := range resp.Candidates {
		for _, p := range c.Content.Parts {
			switch d := p.(type) {
			case genai.Text:
				respText = append(respText, string(d))
			}
		}
	}
	return respText, nil
}
