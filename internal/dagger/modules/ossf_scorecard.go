package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// OSSFScorecardModule runs OSSF Scorecard for open source security scoring
type OSSFScorecardModule struct {
	client *dagger.Client
	name   string
}

const scorecardBinary = "/usr/local/bin/scorecard"

// NewOSSFScorecardModule creates a new OSSF Scorecard module
func NewOSSFScorecardModule(client *dagger.Client) *OSSFScorecardModule {
	return &OSSFScorecardModule{
		client: client,
		name:   "ossf-scorecard",
	}
}

// ScoreRepository scores a repository's security posture
func (m *OSSFScorecardModule) ScoreRepository(ctx context.Context, repoURL string, githubToken string) (string, error) {
	container := m.client.Container().
		From("gcr.io/openssf/scorecard:stable").
		WithEnvVariable("GITHUB_TOKEN", githubToken).
		WithExec([]string{
			scorecardBinary,
			"--repo", repoURL,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run scorecard: %w", err)
	}

	return output, nil
}

// ScoreWithChecks scores repository with specific checks
func (m *OSSFScorecardModule) ScoreWithChecks(ctx context.Context, repoURL string, checks []string, githubToken string) (string, error) {
	args := []string{
		scorecardBinary,
		"--repo", repoURL,
		"--format", "json",
	}

	for _, check := range checks {
		args = append(args, "--checks", check)
	}

	container := m.client.Container().
		From("gcr.io/openssf/scorecard:stable").
		WithEnvVariable("GITHUB_TOKEN", githubToken).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run scorecard with checks: %w", err)
	}

	return output, nil
}

// ListChecks lists available scorecard checks
func (m *OSSFScorecardModule) ListChecks(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("gcr.io/openssf/scorecard:stable").
		WithExec([]string{
			scorecardBinary,
			"--show-details",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list scorecard checks: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of OSSF Scorecard
func (m *OSSFScorecardModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("gcr.io/openssf/scorecard:stable").
		WithExec([]string{scorecardBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scorecard version: %w", err)
	}

	return output, nil
}
