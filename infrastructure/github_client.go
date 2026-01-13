// Package infrastructure はインフラ層を提供します.
// このパッケージは外部API、データベース、ファイルシステムなどの実装を提供します.
package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

// GitHubClient はGitHub APIとの通信を担当するクライアントです.
type GitHubClient struct {
	client             *githubv4.Client
	limiter            *rate.Limiter
	mu                 sync.Mutex
	lastRateLimitReset time.Time
}

// RateLimitInfo はAPIレート制限情報を表します.
type RateLimitInfo struct {
	Remaining int
	ResetAt   time.Time
}

// NewGitHubClient は新しいGitHubClientを作成します
// token: GitHub Personal Access Token（環境変数 GITHUB_TOKEN から読み込む想定）
func NewGitHubClient(token string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// GraphQL APIのrate limitは5000リクエスト/時
	// 安全のため、4500リクエスト/時に制限
	const requestsPerHour = 4500

	limiter := rate.NewLimiter(rate.Every(time.Hour/requestsPerHour), 1)

	return &GitHubClient{
		client:  githubv4.NewClient(tc),
		limiter: limiter,
	}
}

// WaitForRateLimit はrate limitを考慮して待機します.
func (c *GitHubClient) WaitForRateLimit(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// レート制限のリセット時刻を確認
	if !c.lastRateLimitReset.IsZero() && time.Now().Before(c.lastRateLimitReset) {
		waitTime := time.Until(c.lastRateLimitReset)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(waitTime):
		}
	}

	if err := c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter wait failed: %w", err)
	}

	return nil
}

// Query はGraphQLクエリを実行します（rate limit対応）.
func (c *GitHubClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if err := c.WaitForRateLimit(ctx); err != nil {
		return fmt.Errorf("rate limit wait failed: %w", err)
	}

	err := c.client.Query(ctx, q, variables)
	if err != nil {
		return fmt.Errorf("graphql query failed: %w", err)
	}

	return nil
}

// GetRateLimitInfo は現在のrate limit情報を取得します.
func (c *GitHubClient) GetRateLimitInfo(ctx context.Context) (*RateLimitInfo, error) {
	var query struct {
		RateLimit struct {
			Remaining int
			ResetAt   time.Time
		}
	}

	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.lastRateLimitReset = query.RateLimit.ResetAt
	c.mu.Unlock()

	return &RateLimitInfo{
		Remaining: query.RateLimit.Remaining,
		ResetAt:   query.RateLimit.ResetAt,
	}, nil
}
