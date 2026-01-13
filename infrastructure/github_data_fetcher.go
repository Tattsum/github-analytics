package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/shurcooL/githubv4"
)

// GitHubDataFetcher はGitHub APIから各種データを取得するフェッチャーです.
type GitHubDataFetcher struct {
	repo *GitHubRepository
}

// NewGitHubDataFetcher は新しいGitHubDataFetcherを作成します.
func NewGitHubDataFetcher(repo *GitHubRepository) *GitHubDataFetcher {
	return &GitHubDataFetcher{
		repo: repo,
	}
}

// UserActivityData はユーザーの全活動データを表します.
type UserActivityData struct {
	User    *domain.User
	Commits []*domain.Activity
	PRs     []*domain.Activity
	Issues  []*domain.Activity
	Reviews []*domain.Activity
}

// FetchAllUserActivity はユーザーの全活動データを取得します.
func (f *GitHubDataFetcher) FetchAllUserActivity(ctx context.Context, username string, includePrivate bool) (*UserActivityData, error) {
	user, err := f.repo.FetchUserInfo(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// 並列で各種データを取得
	type result struct {
		commits []*domain.Activity
		prs     []*domain.Activity
		issues  []*domain.Activity
		reviews []*domain.Activity
		err     error
	}

	resultChan := make(chan result, 1)

	go func() {
		commits, err1 := f.FetchCommits(ctx, username, includePrivate)
		prs, err2 := f.FetchPullRequests(ctx, username, includePrivate)
		issues, err3 := f.FetchIssues(ctx, username, includePrivate)
		reviews, err4 := f.FetchReviews(ctx, username, includePrivate)

		err := err1
		if err == nil {
			err = err2
		}

		if err == nil {
			err = err3
		}

		if err == nil {
			err = err4
		}

		resultChan <- result{
			commits: commits,
			prs:     prs,
			issues:  issues,
			reviews: reviews,
			err:     err,
		}
	}()

	r := <-resultChan
	if r.err != nil {
		return nil, r.err
	}

	return &UserActivityData{
		User:    user,
		Commits: r.commits,
		PRs:     r.prs,
		Issues:  r.issues,
		Reviews: r.reviews,
	}, nil
}



// findRepositoryInQuery はクエリ結果から指定されたリポジトリを探します.
func (f *GitHubDataFetcher) findRepositoryInQuery(
	repoContribs []struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt githubv4.DateTime
				Commit     struct {
					OID string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
	repoName string,
) (int, bool) {
	for i, repoContrib := range repoContribs {
		if repoContrib.Repository.NameWithOwner == repoName {
			return i, true
		}
	}

	return -1, false
}

// processCommitContributionsFromRepo はリポジトリのContributionsを処理します.
func (f *GitHubDataFetcher) processCommitContributionsFromRepo(
	repoContrib struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt githubv4.DateTime
				Commit     struct {
					OID string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
) []*domain.Activity {
	activities := make([]*domain.Activity, 0)

	for _, contrib := range repoContrib.Contributions.Nodes {
		activity := domain.NewActivity(
			domain.ActivityTypeCommit,
			repoContrib.Repository.NameWithOwner,
			contrib.OccurredAt.Time,
			0, // Additions: GraphQL APIの制限により取得不可
			0, // Deletions: GraphQL APIの制限により取得不可
		)
		activities = append(activities, activity)
	}

	return activities
}

// fetchRepositoryCommitContributionsPage は特定リポジトリのコミット貢献の1ページを取得します.
func (f *GitHubDataFetcher) fetchRepositoryCommitContributionsPage(
	ctx context.Context,
	username string,
	repoName string,
	from, to githubv4.DateTime,
	after *githubv4.String,
) ([]*domain.Activity, *githubv4.String, bool, error) {
	var query struct {
		User struct {
			ContributionsCollection struct {
				CommitContributionsByRepository []struct {
					Repository struct {
						NameWithOwner string
					}
					Contributions struct {
						TotalCount int
						Nodes      []struct {
							OccurredAt githubv4.DateTime
							Commit     struct {
								OID string
							}
						}
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
					} `graphql:"contributions(first: $first, after: $after)"`
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	first := 100
	variables := map[string]interface{}{
		"login": githubv4.String(username),
		"from":  from,
		"to":    to,
		"first": githubv4.Int(first),
		"after": after,
	}

	if err := f.repo.client.Query(ctx, &query, variables); err != nil {
		return nil, nil, false, fmt.Errorf("failed to fetch commit contributions for repository %s: %w", repoName, err)
	}

	idx, found := f.findRepositoryInQuery(query.User.ContributionsCollection.CommitContributionsByRepository, repoName)
	if !found {
		return nil, nil, false, nil
	}

	repoContrib := query.User.ContributionsCollection.CommitContributionsByRepository[idx]
	activities := f.processCommitContributionsFromRepo(repoContrib)

	var nextAfter *githubv4.String

	if repoContrib.Contributions.PageInfo.HasNextPage {
		cursor := githubv4.String(repoContrib.Contributions.PageInfo.EndCursor)
		nextAfter = &cursor
	}

	return activities, nextAfter, repoContrib.Contributions.PageInfo.HasNextPage, nil
}

// fetchRepositoryCommitContributionsPaginated は特定リポジトリのコミット貢献をページネーションで取得します.
// 最初のページは既に取得済みのため、2ページ目以降を取得します.
func (f *GitHubDataFetcher) fetchRepositoryCommitContributionsPaginated(
	ctx context.Context,
	username string,
	repoName string,
	from, to githubv4.DateTime,
	startAfter *githubv4.String,
) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0)
	after := startAfter

	for {
		pageActivities, nextAfter, hasNext, err := f.fetchRepositoryCommitContributionsPage(ctx, username, repoName, from, to, after)
		if err != nil {
			return nil, err
		}

		if pageActivities == nil {
			break
		}

		activities = append(activities, pageActivities...)

		if !hasNext {
			break
		}

		after = nextAfter
	}

	return activities, nil
}

// FetchCommits はコミットを取得します
// 注意: GitHub GraphQL APIのContributionsCollectionでは、コミットのAdditions/Deletionsを
// 直接取得できないため、コミット数と日時のみを取得します。
// 変更行数の詳細が必要な場合は、各リポジトリのコミット履歴を個別に取得する必要があります。
// ページネーション: 各リポジトリのContributionsをページネーションで取得します。
func (f *GitHubDataFetcher) FetchCommits(ctx context.Context, username string, _ bool) ([]*domain.Activity, error) {
	var query struct {
		User struct {
			ContributionsCollection struct {
				CommitContributionsByRepository []struct {
					Repository struct {
						NameWithOwner string
					}
					Contributions struct {
						TotalCount int
						Nodes      []struct {
							OccurredAt githubv4.DateTime
							Commit     struct {
								OID string
							}
						}
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
					} `graphql:"contributions(first: $first, after: $after)"`
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	from := githubv4.DateTime{Time: time.Now().AddDate(-10, 0, 0)}
	to := githubv4.DateTime{Time: time.Now()}
	first := 100
	after := (*githubv4.String)(nil)

	// 最初のクエリでリポジトリのリストを取得
	variables := map[string]interface{}{
		"login": githubv4.String(username),
		"from":  from,
		"to":    to,
		"first": githubv4.Int(first),
		"after": after,
	}

	if err := f.repo.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}

	return f.processCommitsWithPagination(ctx, username, query.User.ContributionsCollection.CommitContributionsByRepository, from, to)
}

// processCommitsWithPagination はコミット貢献をページネーションで処理します.
func (f *GitHubDataFetcher) processCommitsWithPagination(
	ctx context.Context,
	username string,
	repoContribs []struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt githubv4.DateTime
				Commit     struct {
					OID string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
	from, to githubv4.DateTime,
) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0)

	// 各リポジトリのContributionsをページネーションで取得
	for _, repoContrib := range repoContribs {
		repoName := repoContrib.Repository.NameWithOwner

		// 最初のページのContributionsを処理
		for _, contrib := range repoContrib.Contributions.Nodes {
			activity := domain.NewActivity(
				domain.ActivityTypeCommit,
				repoContrib.Repository.NameWithOwner,
				contrib.OccurredAt.Time,
				0, // Additions: GraphQL APIの制限により取得不可
				0, // Deletions: GraphQL APIの制限により取得不可
			)
			activities = append(activities, activity)
		}

		// このリポジトリに次のページがある場合、ページネーションで取得
		if repoContrib.Contributions.PageInfo.HasNextPage {
			startAfter := githubv4.String(repoContrib.Contributions.PageInfo.EndCursor)

			repoActivities, err := f.fetchRepositoryCommitContributionsPaginated(ctx, username, repoName, from, to, &startAfter)
			if err != nil {
				// エラーが発生しても他のリポジトリの処理は続行
				continue
			}

			activities = append(activities, repoActivities...)
		}
	}

	return activities, nil
}

// FetchPullRequests はPull Requestを取得します.
func (f *GitHubDataFetcher) FetchPullRequests(ctx context.Context, username string, _ bool) ([]*domain.Activity, error) {
	var query struct {
		User struct {
			PullRequests struct {
				TotalCount int
				Nodes      []struct {
					Title      string
					CreatedAt  githubv4.DateTime
					MergedAt   *githubv4.DateTime
					Repository struct {
						NameWithOwner string
					}
					Additions int
					Deletions int
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"pullRequests(first: $first, after: $after, states: [OPEN, CLOSED, MERGED])"`
		} `graphql:"user(login: $login)"`
	}

	activities := make([]*domain.Activity, 0)
	first := 100
	after := (*githubv4.String)(nil)

	for {
		variables := map[string]interface{}{
			"login": githubv4.String(username),
			"first": githubv4.Int(first),
			"after": after,
		}

		if err := f.repo.client.Query(ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		for _, pr := range query.User.PullRequests.Nodes {
			activity := domain.NewActivity(
				domain.ActivityTypePR,
				pr.Repository.NameWithOwner,
				pr.CreatedAt.Time,
				pr.Additions,
				pr.Deletions,
			)
			activity.IsMerged = pr.MergedAt != nil
			activities = append(activities, activity)
		}

		if !query.User.PullRequests.PageInfo.HasNextPage {
			break
		}

		cursor := githubv4.String(query.User.PullRequests.PageInfo.EndCursor)
		after = &cursor
	}

	return activities, nil
}

// FetchIssues はIssueを取得します.
func (f *GitHubDataFetcher) FetchIssues(ctx context.Context, username string, _ bool) ([]*domain.Activity, error) {
	var query struct {
		User struct {
			Issues struct {
				TotalCount int
				Nodes      []struct {
					Title      string
					CreatedAt  githubv4.DateTime
					Repository struct {
						NameWithOwner string
					}
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"issues(first: $first, after: $after, states: [OPEN, CLOSED])"`
		} `graphql:"user(login: $login)"`
	}

	activities := make([]*domain.Activity, 0)
	first := 100
	after := (*githubv4.String)(nil)

	for {
		variables := map[string]interface{}{
			"login": githubv4.String(username),
			"first": githubv4.Int(first),
			"after": after,
		}

		if err := f.repo.client.Query(ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch issues: %w", err)
		}

		for _, issue := range query.User.Issues.Nodes {
			activity := domain.NewActivity(
				domain.ActivityTypeIssue,
				issue.Repository.NameWithOwner,
				issue.CreatedAt.Time,
				0,
				0,
			)
			activities = append(activities, activity)
		}

		if !query.User.Issues.PageInfo.HasNextPage {
			break
		}

		cursor := githubv4.String(query.User.Issues.PageInfo.EndCursor)
		after = &cursor
	}

	return activities, nil
}

// findReviewRepositoryInQuery はクエリ結果から指定されたリポジトリを探します.
func (f *GitHubDataFetcher) findReviewRepositoryInQuery(
	repoContribs []struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt        githubv4.DateTime
				PullRequestReview struct {
					State string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
	repoName string,
) (int, bool) {
	for i, repoContrib := range repoContribs {
		if repoContrib.Repository.NameWithOwner == repoName {
			return i, true
		}
	}

	return -1, false
}

// processReviewContributionsFromRepo はリポジトリのレビューContributionsを処理します.
func (f *GitHubDataFetcher) processReviewContributionsFromRepo(
	repoContrib struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt        githubv4.DateTime
				PullRequestReview struct {
					State string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
) []*domain.Activity {
	activities := make([]*domain.Activity, 0)

	for _, contrib := range repoContrib.Contributions.Nodes {
		activity := domain.NewActivity(
			domain.ActivityTypeReview,
			repoContrib.Repository.NameWithOwner,
			contrib.OccurredAt.Time,
			0,
			0,
		)
		activity.IsReview = true
		activities = append(activities, activity)
	}

	return activities
}

// fetchRepositoryReviewContributionsPage は特定リポジトリのレビュー貢献の1ページを取得します.
func (f *GitHubDataFetcher) fetchRepositoryReviewContributionsPage(
	ctx context.Context,
	username string,
	repoName string,
	from, to githubv4.DateTime,
	after *githubv4.String,
) ([]*domain.Activity, *githubv4.String, bool, error) {
	var query struct {
		User struct {
			ContributionsCollection struct {
				PullRequestReviewContributionsByRepository []struct {
					Repository struct {
						NameWithOwner string
					}
					Contributions struct {
						TotalCount int
						Nodes      []struct {
							OccurredAt        githubv4.DateTime
							PullRequestReview struct {
								State string
							}
						}
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
					} `graphql:"contributions(first: $first, after: $after)"`
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	first := 100
	variables := map[string]interface{}{
		"login": githubv4.String(username),
		"from":  from,
		"to":    to,
		"first": githubv4.Int(first),
		"after": after,
	}

	if err := f.repo.client.Query(ctx, &query, variables); err != nil {
		return nil, nil, false, fmt.Errorf("failed to fetch review contributions for repository %s: %w", repoName, err)
	}

	idx, found := f.findReviewRepositoryInQuery(query.User.ContributionsCollection.PullRequestReviewContributionsByRepository, repoName)
	if !found {
		return nil, nil, false, nil
	}

	repoContrib := query.User.ContributionsCollection.PullRequestReviewContributionsByRepository[idx]
	activities := f.processReviewContributionsFromRepo(repoContrib)

	var nextAfter *githubv4.String

	if repoContrib.Contributions.PageInfo.HasNextPage {
		cursor := githubv4.String(repoContrib.Contributions.PageInfo.EndCursor)
		nextAfter = &cursor
	}

	return activities, nextAfter, repoContrib.Contributions.PageInfo.HasNextPage, nil
}

// fetchRepositoryReviewContributionsPaginated は特定リポジトリのレビュー貢献をページネーションで取得します.
// 最初のページは既に取得済みのため、2ページ目以降を取得します.
func (f *GitHubDataFetcher) fetchRepositoryReviewContributionsPaginated(
	ctx context.Context,
	username string,
	repoName string,
	from, to githubv4.DateTime,
	startAfter *githubv4.String,
) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0)
	after := startAfter

	for {
		pageActivities, nextAfter, hasNext, err := f.fetchRepositoryReviewContributionsPage(ctx, username, repoName, from, to, after)
		if err != nil {
			return nil, err
		}

		if pageActivities == nil {
			break
		}

		activities = append(activities, pageActivities...)

		if !hasNext {
			break
		}

		after = nextAfter
	}

	return activities, nil
}

// FetchReviews はPRレビューを取得します.
// ページネーション: 各リポジトリのContributionsをページネーションで取得します。
func (f *GitHubDataFetcher) FetchReviews(ctx context.Context, username string, _ bool) ([]*domain.Activity, error) {
	var query struct {
		User struct {
			ContributionsCollection struct {
				PullRequestReviewContributionsByRepository []struct {
					Repository struct {
						NameWithOwner string
					}
					Contributions struct {
						TotalCount int
						Nodes      []struct {
							OccurredAt        githubv4.DateTime
							PullRequestReview struct {
								State string
							}
						}
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
					} `graphql:"contributions(first: $first, after: $after)"`
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	from := githubv4.DateTime{Time: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)}
	to := githubv4.DateTime{Time: time.Now()}
	first := 100
	after := (*githubv4.String)(nil)

	// 最初のクエリでリポジトリのリストを取得
	variables := map[string]interface{}{
		"login": githubv4.String(username),
		"from":  from,
		"to":    to,
		"first": githubv4.Int(first),
		"after": after,
	}

	if err := f.repo.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	return f.processReviewsWithPagination(ctx, username, query.User.ContributionsCollection.PullRequestReviewContributionsByRepository, from, to)
}

// processReviewsWithPagination はレビュー貢献をページネーションで処理します.
func (f *GitHubDataFetcher) processReviewsWithPagination(
	ctx context.Context,
	username string,
	repoContribs []struct {
		Repository struct {
			NameWithOwner string
		}
		Contributions struct {
			TotalCount int
			Nodes      []struct {
				OccurredAt        githubv4.DateTime
				PullRequestReview struct {
					State string
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"contributions(first: $first, after: $after)"`
	},
	from, to githubv4.DateTime,
) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0)

	// 各リポジトリのContributionsをページネーションで取得
	for _, repoContrib := range repoContribs {
		repoName := repoContrib.Repository.NameWithOwner

		// 最初のページのContributionsを処理
		for _, contrib := range repoContrib.Contributions.Nodes {
			activity := domain.NewActivity(
				domain.ActivityTypeReview,
				repoContrib.Repository.NameWithOwner,
				contrib.OccurredAt.Time,
				0,
				0,
			)
			activity.IsReview = true
			activities = append(activities, activity)
		}

		// このリポジトリに次のページがある場合、ページネーションで取得
		if repoContrib.Contributions.PageInfo.HasNextPage {
			startAfter := githubv4.String(repoContrib.Contributions.PageInfo.EndCursor)

			repoActivities, err := f.fetchRepositoryReviewContributionsPaginated(ctx, username, repoName, from, to, &startAfter)
			if err != nil {
				// エラーが発生しても他のリポジトリの処理は続行
				continue
			}

			activities = append(activities, repoActivities...)
		}
	}

	return activities, nil
}
