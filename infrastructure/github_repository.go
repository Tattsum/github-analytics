package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/shurcooL/githubv4"
)

// GitHubRepository はGitHub APIからデータを取得するリポジトリです.
type GitHubRepository struct {
	client *GitHubClient
}

// NewGitHubRepository は新しいGitHubRepositoryを作成します.
func NewGitHubRepository(client *GitHubClient) *GitHubRepository {
	return &GitHubRepository{
		client: client,
	}
}

// UserInfo はGraphQLから取得するユーザー情報の構造体です.
type UserInfo struct {
	Login     string
	Name      string
	CreatedAt githubv4.DateTime
}

// FetchUserInfo はユーザー情報を取得します.
func (r *GitHubRepository) FetchUserInfo(ctx context.Context, username string) (*domain.User, error) {
	var query struct {
		User struct {
			Login     string
			Name      string
			CreatedAt githubv4.DateTime
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(username),
	}

	if err := r.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	user := domain.NewUser(
		query.User.Login,
		query.User.Name,
		query.User.CreatedAt.Format(time.RFC3339),
	)

	return user, nil
}

// CommitNode はコミット情報を表すGraphQLノードです.
type CommitNode struct {
	OID     string
	Message string
	Author  struct {
		Date githubv4.DateTime
		User struct {
			Login string
		}
	}
	Additions int
	Deletions int
}

// RepositoryNode はリポジトリ情報を表すGraphQLノードです.
type RepositoryNode struct {
	Name          string
	IsPrivate     bool
	DefaultBranch struct {
		Target struct {
			Commit struct {
				History struct {
					TotalCount int
					Nodes      []CommitNode
					PageInfo   struct {
						HasNextPage bool
						EndCursor   string
					}
				} `graphql:"history(first: $first, after: $after, author: $author)"`
			} `graphql:"... on Commit"`
		}
	} `graphql:"defaultBranchRef @include(if: $includeBranch)"`
}

// FetchCommits は指定ユーザーのコミットを取得します.
func (r *GitHubRepository) FetchCommits(ctx context.Context, username string, includePrivate bool) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0)

	// まず、ユーザーがアクセス可能なリポジトリ一覧を取得
	repos, err := r.fetchUserRepositories(ctx, username, includePrivate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	// 各リポジトリのコミットを並列で取得
	type repoCommitResult struct {
		activities []*domain.Activity
		err        error
	}

	results := make(chan repoCommitResult, len(repos))

	const maxConcurrentRepos = 5

	semaphore := make(chan struct{}, maxConcurrentRepos) // 最大5並列

	for _, repo := range repos {
		go func(repoName string) {
			semaphore <- struct{}{}

			defer func() { <-semaphore }()

			commits, err := r.fetchRepositoryCommits(ctx, username, repoName)
			results <- repoCommitResult{activities: commits, err: err}
		}(repo)
	}

	for i := 0; i < len(repos); i++ {
		result := <-results
		if result.err != nil {
			// エラーはログに記録するが、他のリポジトリの処理は続行
			fmt.Printf("Warning: failed to fetch commits for repository: %v\n", result.err)
			continue
		}

		activities = append(activities, result.activities...)
	}

	return activities, nil
}

// fetchUserRepositories はユーザーがアクセス可能なリポジトリ一覧を取得します.
func (r *GitHubRepository) fetchUserRepositories(ctx context.Context, username string, includePrivate bool) ([]string, error) {
	var query struct {
		User struct {
			Repositories struct {
				Nodes []struct {
					Name      string
					IsPrivate bool
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"repositories(first: $first, after: $after, ownerAffiliations: OWNER, affiliations: OWNER_COLLABORATOR_ORGANIZATION_MEMBER)"`
		} `graphql:"user(login: $login)"`
	}

	repos := make([]string, 0)
	first := 100
	after := (*githubv4.String)(nil)

	for {
		variables := map[string]interface{}{
			"login": githubv4.String(username),
			"first": githubv4.Int(first),
			"after": after,
		}

		if err := r.client.Query(ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch repositories: %w", err)
		}

		for _, node := range query.User.Repositories.Nodes {
			if !includePrivate && node.IsPrivate {
				continue
			}

			repos = append(repos, node.Name)
		}

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}

		cursor := githubv4.String(query.User.Repositories.PageInfo.EndCursor)
		after = &cursor
	}

	return repos, nil
}

// fetchRepositoryCommits は特定リポジトリのコミットを取得します.
func (r *GitHubRepository) fetchRepositoryCommits(_ context.Context, _, _ string) ([]*domain.Activity, error) {
	// この実装は簡略化されています
	// 実際には、リポジトリのowner情報も必要です
	// より詳細な実装が必要な場合は、リポジトリの完全な情報（owner/repo）を取得する必要があります
	return nil, domain.ErrNotImplemented
}
