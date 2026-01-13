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

func main() {
	// コマンドライン引数の定義
	var (
		usersStr       = flag.String("users", "", "分析対象のGitHubユーザー名（カンマ区切り、例: user1,user2）")
		orgName        = flag.String("org", "", "分析対象のGitHub組織名（指定した場合、組織のメンバーを分析）")
		outputDir      = flag.String("output", "output", "出力ディレクトリ")
		includePrivate = flag.Bool("private", false, "privateリポジトリも対象にする")
		help           = flag.Bool("help", false, "ヘルプを表示")
	)

	flag.Parse()

	// ヘルプ表示
	if *help {
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

	// 環境変数からトークンを取得
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set. Please set your GitHub Personal Access Token.")
	}

	// ユーザーリストを取得
	var users []string

	if *orgName != "" {
		// 組織のメンバーを取得
		var err error

		users, err = fetchOrganizationMembers(token, *orgName)
		if err != nil {
			log.Fatalf("Failed to fetch organization members: %v", err)
		}

		if len(users) == 0 {
			log.Fatalf("No members found in organization: %s", *orgName)
		}

		fmt.Printf("Found %d members in organization: %s\n", len(users), *orgName)
	} else if *usersStr != "" {
		// コマンドライン引数からユーザー名を取得
		users = strings.Split(*usersStr, ",")
		for i := range users {
			users[i] = strings.TrimSpace(users[i])
		}
	} else {
		log.Fatal("Either -users or -org flag must be specified. Use -help for usage.")
	}

	if len(users) == 0 {
		log.Fatal("No users specified for analysis.")
	}

	// 出力ディレクトリを作成
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// コンテキストを作成（タイムアウト: 30分）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// GitHubクライアントを作成
	client := infrastructure.NewGitHubClient(token)
	repo := infrastructure.NewGitHubRepository(client)
	fetcher := infrastructure.NewGitHubDataFetcher(repo)
	statsService := application.NewStatisticsService()

	// 各ユーザーの統計を並列で取得
	type userResult struct {
		username string
		stats    interface{}
		err      error
	}

	results := make(chan userResult, len(users))

	for _, username := range users {
		go func(user string) {
			fmt.Printf("Processing user: %s\n", user)

			// データ取得
			data, err := fetcher.FetchAllUserActivity(ctx, user, *includePrivate)
			if err != nil {
				results <- userResult{username: user, err: err}
				return
			}

			// 統計計算
			stats, err := statsService.CalculateStatistics(data)
			if err != nil {
				results <- userResult{username: user, err: err}
				return
			}

			results <- userResult{username: user, stats: stats}
		}(username)
	}

	// 結果を収集して出力
	formatter := presentation.NewOutputFormatter(*outputDir)
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

		// 出力生成
		if err := formatter.FormatAll(stats); err != nil {
			log.Printf("Error formatting output for user %s: %v", result.username, err)
			continue
		}

		fmt.Printf("Completed processing user: %s\n", result.username)
	}

	// 統合レポートを生成（オプション）
	if err := generateCombinedReport(*outputDir, allStats); err != nil {
		log.Printf("Error generating combined report: %v", err)
	}

	fmt.Println("\n=== 処理完了 ===")
	fmt.Printf("結果は %s/ ディレクトリに出力されました。\n", *outputDir)
}

// fetchOrganizationMembers は組織のメンバー一覧を取得します.
func fetchOrganizationMembers(token, orgName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
func generateCombinedReport(outputDir string, allStats map[string]interface{}) error {
	// 簡易的な統合レポート
	// より詳細な実装が必要な場合は拡張可能
	return nil
}
