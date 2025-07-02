package hasura

import (
	"context"
	
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	gqlClient *graphql.Client
}

func NewClient(hasuraURL, adminSecret string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: adminSecret},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	
	return &Client{
		gqlClient: graphql.NewClient(hasuraURL, httpClient),
	}
}

// Add your custom query/mutation methods here
func (c *Client) ExecuteQuery(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	return c.gqlClient.Query(ctx, query, variables)
}

func (c *Client) ExecuteMutation(ctx context.Context, mutation interface{}, variables map[string]interface{}) error {
	return c.gqlClient.Mutate(ctx, mutation, variables)
}