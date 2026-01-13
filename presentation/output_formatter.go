// Package presentation はプレゼンテーション層を提供します.
// このパッケージは出力フォーマット、UI、CLIなどの実装を提供します.
package presentation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

const (
	dirPerm     = 0750
	filePerm    = 0600
	hoursPerDay = 24
)

// FormatAll は全ての出力形式を生成します.
func (f *OutputFormatter) FormatAll(stats *domain.UserStatistics) error {
	if err := os.MkdirAll(f.outputDir, dirPerm); err != nil {
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

// buildJSONData はJSON用のデータ構造を構築します.
func (f *OutputFormatter) buildJSONData(stats *domain.UserStatistics) map[string]interface{} {
	jsonData := map[string]interface{}{
		"user":                   stats.User.Login,
		"total_commits":          stats.TotalCommits,
		"total_pr_created":       stats.TotalPRCreated,
		"total_pr_merged":        stats.TotalPRMerged,
		"total_issues":           stats.TotalIssues,
		"total_reviews":          stats.TotalReviews,
		"total_additions":        stats.TotalAdditions,
		"total_deletions":        stats.TotalDeletions,
		"first_activity_year":    stats.FirstActivityYear,
		"peak_activity_year":     stats.PeakActivityYear,
		"peak_activity_commits":  stats.PeakActivityCommits,
		"pr_to_review_ratio":     stats.PRToReviewRatio,
		"yearly_stats":           f.buildYearlyStatsJSON(stats),
		"top_repositories":       f.buildTopRepositoriesJSON(stats),
		"long_term_repositories": f.buildLongTermRepositoriesJSON(stats),
		"role_transition":        f.buildRoleTransitionJSON(stats),
	}

	return jsonData
}

// buildYearlyStatsJSON は年別統計のJSONデータを構築します.
func (f *OutputFormatter) buildYearlyStatsJSON(stats *domain.UserStatistics) map[string]interface{} {
	yearlyStats := make(map[string]interface{})
	for year, yearlyStat := range stats.YearlyStats {
		yearlyStats[fmt.Sprintf("%d", year)] = map[string]interface{}{
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

	return yearlyStats
}

// buildTopRepositoriesJSON はTOP3リポジトリのJSONデータを構築します.
func (f *OutputFormatter) buildTopRepositoriesJSON(stats *domain.UserStatistics) []interface{} {
	repos := make([]interface{}, 0, len(stats.TopRepositories))
	for _, repo := range stats.TopRepositories {
		repos = append(repos, map[string]interface{}{
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

	return repos
}

// buildLongTermRepositoriesJSON は長期間関与リポジトリのJSONデータを構築します.
func (f *OutputFormatter) buildLongTermRepositoriesJSON(stats *domain.UserStatistics) []interface{} {
	repos := make([]interface{}, 0, len(stats.LongTermRepositories))
	for _, repo := range stats.LongTermRepositories {
		repos = append(repos, map[string]interface{}{
			"repository":     repo.Repository,
			"commit_count":   repo.CommitCount,
			"first_activity": repo.FirstActivity.Format(time.RFC3339),
			"last_activity":  repo.LastActivity.Format(time.RFC3339),
			"duration_days":  int(repo.LastActivity.Sub(repo.FirstActivity).Hours() / hoursPerDay),
		})
	}

	return repos
}

// buildRoleTransitionJSON はロール変化のJSONデータを構築します.
func (f *OutputFormatter) buildRoleTransitionJSON(stats *domain.UserStatistics) []interface{} {
	transitions := make([]interface{}, 0, len(stats.RoleTransition))
	for _, transition := range stats.RoleTransition {
		transitions = append(transitions, map[string]interface{}{
			"year":         transition.Year,
			"pr_created":   transition.PRCreated,
			"review_count": transition.ReviewCount,
			"ratio":        transition.Ratio,
			"description":  transition.Description,
		})
	}

	return transitions
}

// OutputJSON はJSON形式で出力します.
func (f *OutputFormatter) OutputJSON(stats *domain.UserStatistics) error {
	filename := filepath.Join(f.outputDir, fmt.Sprintf("%s_statistics.json", stats.User.Login))

	jsonData := f.buildJSONData(stats)

	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, filePerm); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// writeCSVBasicStats は基本統計をCSVに書き込みます.
func (f *OutputFormatter) writeCSVBasicStats(writer *csv.Writer, stats *domain.UserStatistics) error {
	headers := []string{"Metric", "Value"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

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

	return nil
}

// writeCSVYearlyStats は年別統計をCSVに書き込みます.
func (f *OutputFormatter) writeCSVYearlyStats(writer *csv.Writer, stats *domain.UserStatistics) error {
	if err := writer.Write([]string{"", ""}); err != nil {
		return fmt.Errorf("failed to write CSV empty row: %w", err)
	}

	if err := writer.Write([]string{"Yearly Statistics", ""}); err != nil {
		return fmt.Errorf("failed to write CSV section header: %w", err)
	}

	yearlyHeaders := []string{
		"Year", "Commits", "PR Created", "PR Merged", "Issues", "Reviews", "Additions", "Deletions",
	}
	if err := writer.Write(yearlyHeaders); err != nil {
		return fmt.Errorf("failed to write CSV yearly headers: %w", err)
	}

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
			return fmt.Errorf("failed to write CSV yearly row: %w", err)
		}
	}

	return nil
}

// OutputCSV はCSV形式で出力します.
func (f *OutputFormatter) OutputCSV(stats *domain.UserStatistics) error {
	filename := filepath.Join(f.outputDir, fmt.Sprintf("%s_statistics.csv", stats.User.Login))

	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close CSV file: %w", closeErr)
			}
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := f.writeCSVBasicStats(writer, stats); err != nil {
		return err
	}

	if err := f.writeCSVYearlyStats(writer, stats); err != nil {
		return err
	}

	return nil
}

// writeTextSummaryNumbers は数字で表すセクションを書き込みます.
func (f *OutputFormatter) writeTextSummaryNumbers(sb *strings.Builder, stats *domain.UserStatistics) {
	years := time.Now().Year() - stats.FirstActivityYear + 1
	fmt.Fprintf(sb, "・%d年間で%d回のコミット\n", years, stats.TotalCommits)
	fmt.Fprintf(sb, "・%d件のPull Requestを作成し、%d件をマージ\n", stats.TotalPRCreated, stats.TotalPRMerged)
	fmt.Fprintf(sb, "・%d件のIssueを作成\n", stats.TotalIssues)
	fmt.Fprintf(sb, "・%d件のPRレビューを実施\n", stats.TotalReviews)
	fmt.Fprintf(sb, "・合計%d行の追加、%d行の削除\n\n", stats.TotalAdditions, stats.TotalDeletions)
}

// writeTextSummaryCharacteristics はエンジニアとしての特徴を書き込みます.
func (f *OutputFormatter) writeTextSummaryCharacteristics(sb *strings.Builder, stats *domain.UserStatistics) {
	if stats.TotalReviews > stats.TotalPRCreated {
		sb.WriteString("・レビュー活動が活発で、チームのコード品質向上に大きく貢献\n")
	}

	if stats.PRToReviewRatio > 1.0 {
		sb.WriteString("・PR作成数よりもレビュー数が多く、メンター的な役割を果たしている\n")
	}

	if len(stats.LongTermRepositories) > 0 {
		fmt.Fprintf(sb, "・%d個のリポジトリに長期間（1年以上）関与し、継続的な貢献を実現\n", len(stats.LongTermRepositories))
	}

	if stats.PeakActivityCommits > 0 {
		fmt.Fprintf(sb, "・%d年が最も活動的で、%d回のコミットを実施\n", stats.PeakActivityYear, stats.PeakActivityCommits)
	}
}

// writeTextSummaryRoleTransition は役割の変化を書き込みます.
func (f *OutputFormatter) writeTextSummaryRoleTransition(sb *strings.Builder, stats *domain.UserStatistics) {
	for _, transition := range stats.RoleTransition {
		if transition.PRCreated > 0 || transition.ReviewCount > 0 {
			fmt.Fprintf(sb, "・%d年: %s (PR作成: %d, レビュー: %d)\n",
				transition.Year, transition.Description, transition.PRCreated, transition.ReviewCount)
		}
	}
}

// writeTextSummaryRepositories はリポジトリ情報を書き込みます.
func (f *OutputFormatter) writeTextSummaryRepositories(sb *strings.Builder, stats *domain.UserStatistics) {
	for i, repo := range stats.TopRepositories {
		fmt.Fprintf(sb, "%d. %s: %dコミット\n", i+1, repo.Repository, repo.CommitCount)
	}

	if len(stats.LongTermRepositories) > 0 {
		sb.WriteString("\n長期間関与しているリポジトリ\n")

		for _, repo := range stats.LongTermRepositories {
			duration := repo.LastActivity.Sub(repo.FirstActivity)
			fmt.Fprintf(sb, "・%s: %d日間 (初回: %s, 最終: %s)\n",
				repo.Repository,
				int(duration.Hours()/hoursPerDay),
				repo.FirstActivity.Format("2006-01-02"),
				repo.LastActivity.Format("2006-01-02"))
		}
	}
}

// OutputTextSummary はテキスト要約を出力します.
func (f *OutputFormatter) OutputTextSummary(stats *domain.UserStatistics) error {
	filename := filepath.Join(f.outputDir, fmt.Sprintf("%s_summary.txt", stats.User.Login))

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== %s のGitHub活動統計 ===\n\n", stats.User.Login))
	sb.WriteString("この人を数字で表すと\n")
	f.writeTextSummaryNumbers(&sb, stats)

	sb.WriteString("エンジニアとしての特徴\n")
	f.writeTextSummaryCharacteristics(&sb, stats)

	sb.WriteString("\n役割の変化が読み取れるポイント\n")
	f.writeTextSummaryRoleTransition(&sb, stats)

	sb.WriteString("\n最も貢献したリポジトリ TOP3\n")
	f.writeTextSummaryRepositories(&sb, stats)

	if err := os.WriteFile(filepath.Clean(filename), []byte(sb.String()), filePerm); err != nil {
		return fmt.Errorf("failed to write text summary: %w", err)
	}

	return nil
}

// OutputPresentationSummary はプレゼン用短文を出力します.
func (f *OutputFormatter) OutputPresentationSummary(stats *domain.UserStatistics) error {
	filename := filepath.Join(f.outputDir, fmt.Sprintf("%s_presentation.txt", stats.User.Login))

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

	if err := os.WriteFile(filepath.Clean(filename), []byte(sb.String()), filePerm); err != nil {
		return fmt.Errorf("failed to write presentation summary: %w", err)
	}

	return nil
}
