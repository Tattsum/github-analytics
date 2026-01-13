// Package main はGitHub Analyticsアプリケーションのエントリーポイントです.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
	"github.com/Tattsum/github-analytics/presentation"
)

// 実行手順・前提条件:
// 1. GitHub Personal Access Tokenを環境変数 GITHUB_TOKEN に設定してください
//    - トークンの作成: https://github.com/settings/tokens
//    - 必要なスコープ: public_repo, repo (privateリポジトリも対象にする場合)
//    - 組織のメンバーを取得する場合: read:org スコープが必要
// 2. このスクリプトを実行すると、output/ ディレクトリに結果が出力されます
// 3. API制限を考慮して実装されているため、大量のデータがある場合は時間がかかる場合があります
// 4. privateリポジトリも対象にする場合は、トークンに適切な権限が必要です

// showHelp はヘルプを表示します.
func showHelp() {
	flag.Usage()
	fmt.Println("\n使用例:")
	fmt.Println("  # 特定のユーザーを分析")
	fmt.Println("  ./github-analytics -users user1,user2")
	fmt.Println("  # 組織のメンバーを分析")
	fmt.Println("  ./github-analytics -org myorg")
	fmt.Println("  # privateリポジトリも含める")
	fmt.Println("  ./github-analytics -users user1 -private")
	os.Exit(0)
}

// getUsers はユーザーリストを取得します.
func getUsers(orgName, usersStr, token *string) []string {
	var users []string

	switch {
	case *orgName != "":
		var err error

		users, err = fetchOrganizationMembers(*token, *orgName)
		if err != nil {
			log.Fatalf("Failed to fetch organization members: %v", err)
		}

		if len(users) == 0 {
			log.Fatalf("No members found in organization: %s", *orgName)
		}

		fmt.Printf("Found %d members in organization: %s\n", len(users), *orgName)
	case *usersStr != "":
		users = strings.Split(*usersStr, ",")
		for i := range users {
			users[i] = strings.TrimSpace(users[i])
		}
	default:
		log.Fatal("Either -users or -org flag must be specified. Use -help for usage.")
	}

	if len(users) == 0 {
		log.Fatal("No users specified for analysis.")
	}

	return users
}

// processUser はユーザーの統計を処理します.
func processUser(
	ctx context.Context,
	user string,
	includePrivate bool,
	fetcher *infrastructure.GitHubDataFetcher,
	statsService *application.StatisticsService,
) (*domain.UserStatistics, error) {
	fmt.Printf("Processing user: %s\n", user)

	data, err := fetcher.FetchAllUserActivity(ctx, user, includePrivate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user activity: %w", err)
	}

	stats, err := statsService.CalculateStatistics(data)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate statistics: %w", err)
	}

	return stats, nil
}

// collectResults は結果を収集して出力します.
func collectResults(
	results chan userResult,
	users []string,
	formatter *presentation.OutputFormatter,
) map[string]interface{} {
	allStats := make(map[string]interface{})

	for i := 0; i < len(users); i++ {
		result := <-results
		if result.err != nil {
			log.Printf("Error processing user %s: %v", result.username, result.err)
			continue
		}

		stats, ok := result.stats.(*domain.UserStatistics)
		if !ok {
			log.Printf("Error: invalid stats type for user %s", result.username)
			continue
		}

		allStats[result.username] = stats

		if err := formatter.FormatAll(stats); err != nil {
			log.Printf("Error formatting output for user %s: %v", result.username, err)
			continue
		}

		fmt.Printf("Completed processing user: %s\n", result.username)
	}

	return allStats
}

type userResult struct {
	username string
	stats    interface{}
	err      error
}

func main() {
	var (
		usersStr       = flag.String("users", "", "分析対象のGitHubユーザー名（カンマ区切り、例: user1,user2）")
		orgName        = flag.String("org", "", "分析対象のGitHub組織名（指定した場合、組織のメンバーを分析）")
		outputDir      = flag.String("output", "output", "出力ディレクトリ")
		includePrivate = flag.Bool("private", false, "privateリポジトリも対象にする")
		help           = flag.Bool("help", false, "ヘルプを表示")
	)

	flag.Parse()

	if *help {
		showHelp()
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set. Please set your GitHub Personal Access Token.")
	}

	users := getUsers(orgName, usersStr, &token)

	setupAndProcessUsers(users, *outputDir, *includePrivate, token)

	fmt.Println("\n=== 処理完了 ===")
	fmt.Printf("結果は %s/ ディレクトリに出力されました。\n", *outputDir)
}

// setupAndProcessUsers はユーザー処理のセットアップと実行を行います.
func setupAndProcessUsers(users []string, outputDir string, includePrivate bool, token string) {
	const (
		dirPerm        = 0750
		timeoutMinutes = 30
	)

	if err := os.MkdirAll(outputDir, dirPerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes*time.Minute)
	defer cancel()

	client := infrastructure.NewGitHubClient(token)
	repo := infrastructure.NewGitHubRepository(client)
	fetcher := infrastructure.NewGitHubDataFetcher(repo)
	statsService := application.NewStatisticsService()

	results := make(chan userResult, len(users))

	for _, username := range users {
		go func(user string) {
			stats, err := processUser(ctx, user, includePrivate, fetcher, statsService)
			if err != nil {
				results <- userResult{username: user, err: err}
				return
			}

			results <- userResult{username: user, stats: stats}
		}(username)
	}

	formatter := presentation.NewOutputFormatter(outputDir)
	allStats := collectResults(results, users, formatter)

	if err := generateCombinedReport(outputDir, allStats); err != nil {
		log.Printf("Error generating combined report: %v", err)
	}
}

// fetchOrganizationMembers は組織のメンバー一覧を取得します.
func fetchOrganizationMembers(token, orgName string) ([]string, error) {
	const orgFetchTimeoutMinutes = 5

	ctx, cancel := context.WithTimeout(context.Background(), orgFetchTimeoutMinutes*time.Minute)
	defer cancel()

	client := infrastructure.NewGitHubClient(token)

	var query struct {
		Organization struct {
			MembersWithRole struct {
				Nodes []struct {
					Login string
				}
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"membersWithRole(first: $first, after: $after)"`
		} `graphql:"organization(login: $login)"`
	}

	members := make([]string, 0)
	first := 100
	after := (*githubv4.String)(nil)

	for {
		variables := map[string]interface{}{
			"login": githubv4.String(orgName),
			"first": githubv4.Int(first),
			"after": after,
		}

		if err := client.Query(ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch organization members: %w", err)
		}

		for _, node := range query.Organization.MembersWithRole.Nodes {
			members = append(members, node.Login)
		}

		if !query.Organization.MembersWithRole.PageInfo.HasNextPage {
			break
		}

		cursor := githubv4.String(query.Organization.MembersWithRole.PageInfo.EndCursor)
		after = &cursor
	}

	return members, nil
}

// generateCombinedReport は複数ユーザーの統合レポートを生成します.
func generateCombinedReport(_ string, _ map[string]interface{}) error {
	// 簡易的な統合レポート
	// より詳細な実装が必要な場合は拡張可能
	return nil
}
