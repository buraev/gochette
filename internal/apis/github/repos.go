package github

import (
	"context"
	"errors"
	"fmt"
	"lightweight-cache-proxy-service/internal/apis"
	"lightweight-cache-proxy-service/pgc/lcp"
	"net"
	"time"

	"github.com/buraev/barelog"
	"github.com/shurcooL/githubv4"
)

type pinnedItemsQuery struct {
	Viewer struct {
		PinnedItems struct {
			Nodes []struct {
				Repository struct {
					Name  githubv4.String
					Owner struct {
						Login githubv4.String
					}
					PrimaryLanguage *struct { // Сделано указателем, т.к. может быть nil
						Name  githubv4.String
						Color githubv4.String
					}
					Description githubv4.String
					UpdatedAt   githubv4.DateTime
					IsPrivate   githubv4.Boolean
					ID          githubv4.ID
					URL         githubv4.URI
				} `graphql:"... on Repository"`
			}
		} `graphql:"pinnedItems(first: 6, types: REPOSITORY)"`
	}
}

func fetchPinnedRepos(client *githubv4.Client) ([]lcp.GitHubRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query pinnedItemsQuery
	err := client.Query(ctx, &query, nil)

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		barelog.Warn(cacheInstance.LogPrefix(), "connection timed out while fetching pinned repos")
		return nil, apis.WarningError
	}
	if err != nil {
		barelog.Warn(cacheInstance.LogPrefix(), "error querying GitHub GraphQL API: %v", err)
		return nil, fmt.Errorf("querying github's graphql API failed: %w", err)
	}

	nodes := query.Viewer.PinnedItems.Nodes
	repositories := make([]lcp.GitHubRepository, 0, len(nodes))

	for _, node := range nodes {
		repo := node.Repository

		langName := ""
		langColor := ""
		if repo.PrimaryLanguage != nil {
			langName = string(repo.PrimaryLanguage.Name)
			langColor = string(repo.PrimaryLanguage.Color)
		}

		repositories = append(repositories, lcp.GitHubRepository{
			Name:          string(repo.Name),
			Owner:         string(repo.Owner.Login),
			Language:      langName,
			LanguageColor: langColor,
			Description:   string(repo.Description),
			UpdatedAt:     repo.UpdatedAt.Time,
			ID:            fmt.Sprint(repo.ID),
			URL:           repo.URL.String(),
		})
	}

	return repositories, nil
}
