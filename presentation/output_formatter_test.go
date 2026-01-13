package presentation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Tattsum/github-analytics/domain"
)

func TestNewOutputFormatter(t *testing.T) {
	t.Parallel()

	formatter := NewOutputFormatter("test-output")
	assert.NotNil(t, formatter, "NewOutputFormatter() should not return nil")
	assert.Equal(t, "test-output", formatter.outputDir, "outputDir should match")
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
	require.NoError(t, err, "FormatAll() should not return error")

	// ファイルが作成されていることを確認
	expectedFiles := []string{
		"testuser_statistics.json",
		"testuser_statistics.csv",
		"testuser_summary.txt",
		"testuser_presentation.txt",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(tmpDir, filename)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "FormatAll() should create file: %s", filename)
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
	require.NoError(t, err, "OutputJSON() should not return error")

	filePath := filepath.Join(tmpDir, "testuser_statistics.json")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "OutputJSON() should create JSON file")
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
	require.NoError(t, err, "OutputCSV() should not return error")

	filePath := filepath.Join(tmpDir, "testuser_statistics.csv")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "OutputCSV() should create CSV file")
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
	require.NoError(t, err, "OutputTextSummary() should not return error")

	filePath := filepath.Join(tmpDir, "testuser_summary.txt")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "OutputTextSummary() should create text summary file")
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
	require.NoError(t, err, "OutputPresentationSummary() should not return error")

	filePath := filepath.Join(tmpDir, "testuser_presentation.txt")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "OutputPresentationSummary() should create presentation summary file")
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
	require.NoError(t, err, "FormatAll() should not return error")

	// ディレクトリが作成されていることを確認
	_, err = os.Stat(newDir)
	assert.NoError(t, err, "FormatAll() should create output directory")
}
