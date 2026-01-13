package presentation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Tattsum/github-analytics/domain"
)

// OutputFormatter は出力フォーマットを担当するフォーマッターです.
type OutputFormatter struct {
	outputDir string
}

// NewOutputFormatter は新しいOutputFormatterを作成します.
func NewOutputFormatter(outputDir string) *OutputFormatter {
	return &OutputFormatter{
		outputDir: outputDir,
	}
}

// FormatAll は全ての出力形式を生成します.
func (f *OutputFormatter) FormatAll(stats *domain.UserStatistics) error {
	if err := os.MkdirAll(f.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// JSON形式で出力
	if err := f.OutputJSON(stats); err != nil {
		return fmt.Errorf("failed to output JSON: %w", err)
	}

	// CSV形式で出力
	if err := f.OutputCSV(stats); err != nil {
		return fmt.Errorf("failed to output CSV: %w", err)
	}

	// テキスト要約を出力
	if err := f.OutputTextSummary(stats); err != nil {
		return fmt.Errorf("failed to output text summary: %w", err)
	}

	// プレゼン用短文を出力
	if err := f.OutputPresentationSummary(stats); err != nil {
		return fmt.Errorf("failed to output presentation summary: %w", err)
	}

	return nil
}

// OutputJSON はJSON形式で出力します.
func (f *OutputFormatter) OutputJSON(stats *domain.UserStatistics) error {
	filename := fmt.Sprintf("%s/%s_statistics.json", f.outputDir, stats.User.Login)

	// JSON用の構造体を作成
	jsonData := struct {
		User                 string                 `json:"user"`
		TotalCommits         int                    `json:"total_commits"`
		TotalPRCreated       int                    `json:"total_pr_created"`
		TotalPRMerged        int                    `json:"total_pr_merged"`
		TotalIssues          int                    `json:"total_issues"`
		TotalReviews         int                    `json:"total_reviews"`
		TotalAdditions       int                    `json:"total_additions"`
		TotalDeletions       int                    `json:"total_deletions"`
		FirstActivityYear    int                    `json:"first_activity_year"`
		PeakActivityYear     int                    `json:"peak_activity_year"`
		PeakActivityCommits  int                    `json:"peak_activity_commits"`
		PRToReviewRatio      float64                `json:"pr_to_review_ratio"`
		YearlyStats          map[string]interface{} `json:"yearly_stats"`
		TopRepositories      []interface{}          `json:"top_repositories"`
		LongTermRepositories []interface{}          `json:"long_term_repositories"`
		RoleTransition       []interface{}          `json:"role_transition"`
	}{
		User:                 stats.User.Login,
		TotalCommits:         stats.TotalCommits,
		TotalPRCreated:       stats.TotalPRCreated,
		TotalPRMerged:        stats.TotalPRMerged,
		TotalIssues:          stats.TotalIssues,
		TotalReviews:         stats.TotalReviews,
		TotalAdditions:       stats.TotalAdditions,
		TotalDeletions:       stats.TotalDeletions,
		FirstActivityYear:    stats.FirstActivityYear,
		PeakActivityYear:     stats.PeakActivityYear,
		PeakActivityCommits:  stats.PeakActivityCommits,
		PRToReviewRatio:      stats.PRToReviewRatio,
		YearlyStats:          make(map[string]interface{}),
		TopRepositories:      make([]interface{}, 0),
		LongTermRepositories: make([]interface{}, 0),
		RoleTransition:       make([]interface{}, 0),
	}

	// 年別統計を変換
	for year, yearlyStat := range stats.YearlyStats {
		jsonData.YearlyStats[fmt.Sprintf("%d", year)] = map[string]interface{}{
			"year":         yearlyStat.Year,
			"commit_count": yearlyStat.CommitCount,
			"pr_created":   yearlyStat.PRCreated,
			"pr_merged":    yearlyStat.PRMerged,
			"issue_count":  yearlyStat.IssueCount,
			"review_count": yearlyStat.ReviewCount,
			"additions":    yearlyStat.TotalAdditions,
			"deletions":    yearlyStat.TotalDeletions,
		}
	}

	// TOP3リポジトリを変換
	for _, repo := range stats.TopRepositories {
		jsonData.TopRepositories = append(jsonData.TopRepositories, map[string]interface{}{
			"repository":     repo.Repository,
			"commit_count":   repo.CommitCount,
			"pr_count":       repo.PRCount,
			"issue_count":    repo.IssueCount,
			"review_count":   repo.ReviewCount,
			"additions":      repo.TotalAdditions,
			"deletions":      repo.TotalDeletions,
			"first_activity": repo.FirstActivity.Format(time.RFC3339),
			"last_activity":  repo.LastActivity.Format(time.RFC3339),
		})
	}

	// 長期間関与リポジトリを変換
	for _, repo := range stats.LongTermRepositories {
		jsonData.LongTermRepositories = append(jsonData.LongTermRepositories, map[string]interface{}{
			"repository":     repo.Repository,
			"commit_count":   repo.CommitCount,
			"first_activity": repo.FirstActivity.Format(time.RFC3339),
			"last_activity":  repo.LastActivity.Format(time.RFC3339),
			"duration_days":  int(repo.LastActivity.Sub(repo.FirstActivity).Hours() / 24),
		})
	}

	// ロール変化を変換
	for _, transition := range stats.RoleTransition {
		jsonData.RoleTransition = append(jsonData.RoleTransition, map[string]interface{}{
			"year":         transition.Year,
			"pr_created":   transition.PRCreated,
			"review_count": transition.ReviewCount,
			"ratio":        transition.Ratio,
			"description":  transition.Description,
		})
	}

	// JSONファイルに書き込み
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// OutputCSV はCSV形式で出力します.
func (f *OutputFormatter) OutputCSV(stats *domain.UserStatistics) error {
	filename := fmt.Sprintf("%s/%s_statistics.csv", f.outputDir, stats.User.Login)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 基本統計のヘッダー
	headers := []string{
		"Metric", "Value",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// 基本統計を書き込み
	basicStats := [][]string{
		{"Total Commits", fmt.Sprintf("%d", stats.TotalCommits)},
		{"Total PR Created", fmt.Sprintf("%d", stats.TotalPRCreated)},
		{"Total PR Merged", fmt.Sprintf("%d", stats.TotalPRMerged)},
		{"Total Issues", fmt.Sprintf("%d", stats.TotalIssues)},
		{"Total Reviews", fmt.Sprintf("%d", stats.TotalReviews)},
		{"Total Additions", fmt.Sprintf("%d", stats.TotalAdditions)},
		{"Total Deletions", fmt.Sprintf("%d", stats.TotalDeletions)},
		{"First Activity Year", fmt.Sprintf("%d", stats.FirstActivityYear)},
		{"Peak Activity Year", fmt.Sprintf("%d", stats.PeakActivityYear)},
		{"Peak Activity Commits", fmt.Sprintf("%d", stats.PeakActivityCommits)},
		{"PR to Review Ratio", fmt.Sprintf("%.2f", stats.PRToReviewRatio)},
	}

	for _, row := range basicStats {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	// 年別統計のセクション
	if err := writer.Write([]string{"", ""}); err != nil {
		return err
	}

	if err := writer.Write([]string{"Yearly Statistics", ""}); err != nil {
		return err
	}

	yearlyHeaders := []string{
		"Year", "Commits", "PR Created", "PR Merged", "Issues", "Reviews", "Additions", "Deletions",
	}
	if err := writer.Write(yearlyHeaders); err != nil {
		return err
	}

	// 年をソート
	years := make([]int, 0, len(stats.YearlyStats))
	for year := range stats.YearlyStats {
		years = append(years, year)
	}

	sort.Ints(years)

	for _, year := range years {
		yearlyStat := stats.YearlyStats[year]

		row := []string{
			fmt.Sprintf("%d", year),
			fmt.Sprintf("%d", yearlyStat.CommitCount),
			fmt.Sprintf("%d", yearlyStat.PRCreated),
			fmt.Sprintf("%d", yearlyStat.PRMerged),
			fmt.Sprintf("%d", yearlyStat.IssueCount),
			fmt.Sprintf("%d", yearlyStat.ReviewCount),
			fmt.Sprintf("%d", yearlyStat.TotalAdditions),
			fmt.Sprintf("%d", yearlyStat.TotalDeletions),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// OutputTextSummary はテキスト要約を出力します.
func (f *OutputFormatter) OutputTextSummary(stats *domain.UserStatistics) error {
	filename := fmt.Sprintf("%s/%s_summary.txt", f.outputDir, stats.User.Login)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== %s のGitHub活動統計 ===\n\n", stats.User.Login))
	sb.WriteString("この人を数字で表すと\n")

	years := time.Now().Year() - stats.FirstActivityYear + 1
	sb.WriteString(fmt.Sprintf("・%d年間で%d回のコミット\n", years, stats.TotalCommits))
	sb.WriteString(fmt.Sprintf("・%d件のPull Requestを作成し、%d件をマージ\n", stats.TotalPRCreated, stats.TotalPRMerged))
	sb.WriteString(fmt.Sprintf("・%d件のIssueを作成\n", stats.TotalIssues))
	sb.WriteString(fmt.Sprintf("・%d件のPRレビューを実施\n", stats.TotalReviews))
	sb.WriteString(fmt.Sprintf("・合計%d行の追加、%d行の削除\n\n", stats.TotalAdditions, stats.TotalDeletions))

	sb.WriteString("エンジニアとしての特徴\n")

	if stats.TotalReviews > stats.TotalPRCreated {
		sb.WriteString("・レビュー活動が活発で、チームのコード品質向上に大きく貢献\n")
	}

	if stats.PRToReviewRatio > 1.0 {
		sb.WriteString("・PR作成数よりもレビュー数が多く、メンター的な役割を果たしている\n")
	}

	if len(stats.LongTermRepositories) > 0 {
		sb.WriteString(fmt.Sprintf("・%d個のリポジトリに長期間（1年以上）関与し、継続的な貢献を実現\n", len(stats.LongTermRepositories)))
	}

	if stats.PeakActivityCommits > 0 {
		sb.WriteString(fmt.Sprintf("・%d年が最も活動的で、%d回のコミットを実施\n", stats.PeakActivityYear, stats.PeakActivityCommits))
	}

	sb.WriteString("\n役割の変化が読み取れるポイント\n")

	for _, transition := range stats.RoleTransition {
		if transition.PRCreated > 0 || transition.ReviewCount > 0 {
			sb.WriteString(fmt.Sprintf("・%d年: %s (PR作成: %d, レビュー: %d)\n",
				transition.Year, transition.Description, transition.PRCreated, transition.ReviewCount))
		}
	}

	sb.WriteString("\n最も貢献したリポジトリ TOP3\n")

	for i, repo := range stats.TopRepositories {
		sb.WriteString(fmt.Sprintf("%d. %s: %dコミット\n", i+1, repo.Repository, repo.CommitCount))
	}

	if len(stats.LongTermRepositories) > 0 {
		sb.WriteString("\n長期間関与しているリポジトリ\n")

		for _, repo := range stats.LongTermRepositories {
			duration := repo.LastActivity.Sub(repo.FirstActivity)
			sb.WriteString(fmt.Sprintf("・%s: %d日間 (初回: %s, 最終: %s)\n",
				repo.Repository,
				int(duration.Hours()/24),
				repo.FirstActivity.Format("2006-01-02"),
				repo.LastActivity.Format("2006-01-02")))
		}
	}

	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text summary: %w", err)
	}

	return nil
}

// OutputPresentationSummary はプレゼン用短文を出力します.
func (f *OutputFormatter) OutputPresentationSummary(stats *domain.UserStatistics) error {
	filename := fmt.Sprintf("%s/%s_presentation.txt", f.outputDir, stats.User.Login)

	var sb strings.Builder

	years := time.Now().Year() - stats.FirstActivityYear + 1

	sb.WriteString(fmt.Sprintf("=== %s の送別会用プレゼン素材 ===\n\n", stats.User.Login))
	sb.WriteString("【スライド1枚用の短文】\n\n")

	// 箇条書き3-4行を生成
	bullets := make([]string, 0)

	// 活動年数とコミット数
	bullets = append(bullets, fmt.Sprintf("・%d年間で%d回のコミットを実施", years, stats.TotalCommits))

	// PR関連
	if stats.TotalPRCreated > 0 {
		bullets = append(bullets, fmt.Sprintf("・%d件のPull Requestを作成し、%d件をマージ", stats.TotalPRCreated, stats.TotalPRMerged))
	}

	// レビュー関連
	if stats.TotalReviews > 0 {
		bullets = append(bullets, fmt.Sprintf("・%d件のPRレビューを実施し、チームの品質向上に貢献", stats.TotalReviews))
	}

	// 変更行数
	if stats.TotalAdditions > 0 || stats.TotalDeletions > 0 {
		bullets = append(bullets, fmt.Sprintf("・合計%d行の追加、%d行の削除でコードベースを進化", stats.TotalAdditions, stats.TotalDeletions))
	}

	// 長期間関与リポジトリ
	if len(stats.LongTermRepositories) > 0 {
		bullets = append(bullets, fmt.Sprintf("・%d個のリポジトリに長期間関与し、継続的な価値を創出", len(stats.LongTermRepositories)))
	}

	// 最大4行まで
	maxBullets := 4
	if len(bullets) > maxBullets {
		bullets = bullets[:maxBullets]
	}

	for _, bullet := range bullets {
		sb.WriteString(bullet + "\n")
	}

	sb.WriteString("\n【補足情報】\n")
	sb.WriteString(fmt.Sprintf("・最も活動的だった年: %d年（%dコミット）\n", stats.PeakActivityYear, stats.PeakActivityCommits))

	if len(stats.TopRepositories) > 0 {
		sb.WriteString(fmt.Sprintf("・最も貢献したリポジトリ: %s（%dコミット）\n", stats.TopRepositories[0].Repository, stats.TopRepositories[0].CommitCount))
	}

	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write presentation summary: %w", err)
	}

	return nil
}
