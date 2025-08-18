package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// HistoryScrubModule runs Git history cleaning tools
type HistoryScrubModule struct {
	client *dagger.Client
	name   string
}

const (
	gitFilterRepoBinary = "/usr/local/bin/git-filter-repo"
	historyScrubGitleaksBinary = "/usr/local/bin/gitleaks"
)

// NewHistoryScrubModule creates a new Git history scrub module
func NewHistoryScrubModule(client *dagger.Client) *HistoryScrubModule {
	return &HistoryScrubModule{
		client: client,
		name:   "history-scrub",
	}
}

// RemoveSecretsWithBFG removes secrets using BFG Repo-Cleaner
func (m *HistoryScrubModule) RemoveSecretsWithBFG(ctx context.Context, repoPath string, secretsFile string, dryRun bool) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git", "openjdk11", "bash"}).
		WithExec([]string{"wget", "-O", "/usr/local/bin/bfg.jar", "https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar"}).
		WithDirectory("/repo", m.client.Host().Directory(repoPath)).
		WithWorkdir("/repo")

	if secretsFile != "" {
		container = container.WithFile("/secrets.txt", m.client.Host().File(secretsFile))
	}

	var cmd []string
	if dryRun {
		// BFG doesn't have direct dry-run, so we'll analyze instead
		cmd = []string{"sh", "-c", "git log --oneline --all | wc -l && echo 'Would run: java -jar /usr/local/bin/bfg.jar --replace-text /secrets.txt .'"}
	} else {
		cmd = []string{"sh", "-c", "java -jar /usr/local/bin/bfg.jar --replace-text /secrets.txt . && git reflog expire --expire=now --all && git gc --prune=now --aggressive"}
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run BFG repo cleaner: %w", err)
	}

	return output, nil
}

// RemoveSecretsWithGitFilter removes secrets using git-filter-repo
func (m *HistoryScrubModule) RemoveSecretsWithGitFilter(ctx context.Context, repoPath string, patternsFile string, dryRun bool) (string, error) {
	container := m.client.Container().
		From("python:3.11-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"pip", "install", "git-filter-repo"}).
		WithDirectory("/repo", m.client.Host().Directory(repoPath)).
		WithWorkdir("/repo")

	if patternsFile != "" {
		container = container.WithFile("/patterns.txt", m.client.Host().File(patternsFile))
	}

	cmd := []string{gitFilterRepoBinary, "--replace-text", "/patterns.txt"}
	if dryRun {
		cmd = append(cmd, "--dry-run")
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run git-filter-repo: %w", err)
	}

	return output, nil
}

// CreateBareClone creates a bare clone for safe history rewriting
func (m *HistoryScrubModule) CreateBareClone(ctx context.Context, sourceRepo string, clonePath string) (string, error) {
	container := m.client.Container().
		From("alpine/git:latest").
		WithExec([]string{"git", "clone", "--bare", sourceRepo, "/clone"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create bare clone: %w", err)
	}

	return output, nil
}

// VerifyHistoryClean verifies secrets have been removed from history
func (m *HistoryScrubModule) VerifyHistoryClean(ctx context.Context, repoPath string, scanTool string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithDirectory("/repo", m.client.Host().Directory(repoPath)).
		WithWorkdir("/repo")

	var cmd []string
	switch scanTool {
	case "gitleaks":
		container = container.WithExec([]string{"wget", "-O", "/usr/local/bin/gitleaks", "https://github.com/gitleaks/gitleaks/releases/latest/download/gitleaks_linux_x64"}).
			WithExec([]string{"chmod", "+x", "/usr/local/bin/gitleaks"})
		cmd = []string{historyScrubGitleaksBinary, "detect", "--source", ".", "--format", "json"}
	case "git-secrets":
		container = container.WithExec([]string{"apk", "add", "--no-cache", "bash"}).
			WithExec([]string{"git", "clone", "https://github.com/awslabs/git-secrets.git", "/git-secrets"}).
			WithExec([]string{"make", "-C", "/git-secrets", "install"})
		cmd = []string{"git", "secrets", "--scan-history"}
	default:
		cmd = []string{"git", "log", "--all", "--grep=secret", "--grep=password", "--grep=key", "--oneline"}
	}

	container = container.WithExec(cmd, dagger.ContainerWithExecOpts{
		Expect: "ANY", // May return non-zero if secrets found
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify history clean: %w", err)
	}

	return output, nil
}

// AnalyzeRepoSize analyzes repository size before and after cleaning
func (m *HistoryScrubModule) AnalyzeRepoSize(ctx context.Context, repoPath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithDirectory("/repo", m.client.Host().Directory(repoPath)).
		WithWorkdir("/repo").
		WithExec([]string{
			"sh", "-c",
			`
				echo "Repository Size Analysis:"
				echo "========================"
				echo "Total size: $(du -sh . | cut -f1)"
				echo "Git objects: $(git count-objects -v)"
				echo "Commits: $(git log --oneline --all | wc -l)"
				echo "Branches: $(git branch -a | wc -l)"
			`,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze repo size: %w", err)
	}

	return output, nil
}
