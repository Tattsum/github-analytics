package presentation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tattsum/github-analytics/domain"
)

func TestNewOutputFormatter(t *testing.T) {
	t.Parallel()

	formatter := NewOutputFormatter("test-output")
	if formatter == nil {
		t.Error("NewOutputFormatter() should not return nil")
		return
	}

	if formatter.outputDir != "test-output" {
		t.Errorf("NewOutputFormatter().outputDir = %v, want test-output", formatter.outputDir)
	}
}

func TestOutputFormatter_FormatAll(t *testing.T) {
	t.Parallel()

	// 一時ディレクトリを作成
	tmpDir := t.TempDir()
	formatter := NewOutputFormatter(tmpDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)
	stats.TotalCommits = 10
	stats.TotalPRCreated = 5
	stats.TotalPRMerged = 3
	stats.TotalIssues = 2
	stats.TotalReviews = 4

	err := formatter.FormatAll(stats)
	if err != nil {
		t.Fatalf("FormatAll() error = %v", err)
	}

	// ファイルが作成されていることを確認
	expectedFiles := []string{
		"testuser_statistics.json",
		"testuser_statistics.csv",
		"testuser_summary.txt",
		"testuser_presentation.txt",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("FormatAll() did not create file: %s", filename)
		}
	}
}

func TestOutputFormatter_OutputJSON(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	formatter := NewOutputFormatter(tmpDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)
	stats.TotalCommits = 10
	stats.FirstActivityYear = 2020

	err := formatter.OutputJSON(stats)
	if err != nil {
		t.Fatalf("OutputJSON() error = %v", err)
	}

	filePath := filepath.Join(tmpDir, "testuser_statistics.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("OutputJSON() did not create JSON file")
	}
}

func TestOutputFormatter_OutputCSV(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	formatter := NewOutputFormatter(tmpDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)
	stats.TotalCommits = 10
	stats.YearlyStats[2020] = domain.NewYearlyStatistics(2020)
	stats.YearlyStats[2020].CommitCount = 10

	err := formatter.OutputCSV(stats)
	if err != nil {
		t.Fatalf("OutputCSV() error = %v", err)
	}

	filePath := filepath.Join(tmpDir, "testuser_statistics.csv")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("OutputCSV() did not create CSV file")
	}
}

func TestOutputFormatter_OutputTextSummary(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	formatter := NewOutputFormatter(tmpDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)
	stats.TotalCommits = 10
	stats.TotalPRCreated = 5
	stats.TotalReviews = 4
	stats.FirstActivityYear = 2020
	stats.PeakActivityYear = 2021
	stats.PeakActivityCommits = 10

	err := formatter.OutputTextSummary(stats)
	if err != nil {
		t.Fatalf("OutputTextSummary() error = %v", err)
	}

	filePath := filepath.Join(tmpDir, "testuser_summary.txt")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("OutputTextSummary() did not create text summary file")
	}
}

func TestOutputFormatter_OutputPresentationSummary(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	formatter := NewOutputFormatter(tmpDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)
	stats.TotalCommits = 10
	stats.TotalPRCreated = 5
	stats.TotalReviews = 4
	stats.FirstActivityYear = 2020

	err := formatter.OutputPresentationSummary(stats)
	if err != nil {
		t.Fatalf("OutputPresentationSummary() error = %v", err)
	}

	filePath := filepath.Join(tmpDir, "testuser_presentation.txt")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("OutputPresentationSummary() did not create presentation summary file")
	}
}

func TestOutputFormatter_FormatAll_CreatesDirectory(t *testing.T) {
	t.Parallel()

	// 存在しないディレクトリを指定
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new-dir")
	formatter := NewOutputFormatter(newDir)

	user := domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	stats := domain.NewUserStatistics(user)

	err := formatter.FormatAll(stats)
	if err != nil {
		t.Fatalf("FormatAll() error = %v", err)
	}

	// ディレクトリが作成されていることを確認
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("FormatAll() did not create output directory")
	}
}
