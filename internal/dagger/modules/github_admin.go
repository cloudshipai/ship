package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GitHubAdminModule provides GitHub administration tools
type GitHubAdminModule struct {
	client *dagger.Client
	name   string
}

// NewGitHubAdminModule creates a new GitHub admin module
func NewGitHubAdminModule(client *dagger.Client) *GitHubAdminModule {
	return &GitHubAdminModule{
		client: client,
		name:   "github-admin",
	}
}

// GetOrgMembers gets organization members
func (m *GitHubAdminModule) GetOrgMembers(ctx context.Context, org string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/orgs/%s/members | jq .`, org),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get org members: %w", err)
	}

	return output, nil
}

// GetRepoPermissions gets repository permissions
func (m *GitHubAdminModule) GetRepoPermissions(ctx context.Context, owner string, repo string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/%s/%s/collaborators | jq .`, owner, repo),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get repo permissions: %w", err)
	}

	return output, nil
}

// AuditOrgSecurity audits organization security settings
func (m *GitHubAdminModule) AuditOrgSecurity(ctx context.Context, org string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`
				echo "Organization Security Audit for %s"
				echo "================================"
				echo "Organization settings:"
				curl -s -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/orgs/%s | jq '{two_factor_requirement_enabled, members_can_create_repositories, default_repository_permission}'
				echo "Security advisories:"
				curl -s -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/orgs/%s/security-advisories | jq length
			`, org, org, org),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to audit org security: %w", err)
	}

	return output, nil
}

// ListOrgRepos lists repositories in an organization
func (m *GitHubAdminModule) ListOrgRepos(ctx context.Context, org string, visibility string, token string) (string, error) {
	args := []string{"gh", "repo", "list", org}
	if visibility != "" {
		args = append(args, "--visibility", visibility)
	}
	args = append(args, "--json", "name,visibility,isPrivate,createdAt")

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org repos: %w", err)
	}

	return output, nil
}

// CreateOrgRepo creates a repository in an organization
func (m *GitHubAdminModule) CreateOrgRepo(ctx context.Context, org string, repoName string, isPrivate bool, token string) (string, error) {
	args := []string{"gh", "repo", "create", org + "/" + repoName}
	if isPrivate {
		args = append(args, "--private")
	} else {
		args = append(args, "--public")
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create org repo: %w", err)
	}

	return output, nil
}

// GetRepoInfoDetailed gets detailed repository information
func (m *GitHubAdminModule) GetRepoInfoDetailed(ctx context.Context, owner string, repo string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{"gh", "repo", "view", owner + "/" + repo, "--json", "name,owner,visibility,createdAt,updatedAt,stargazerCount,forkCount"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get repo info: %w", err)
	}

	return output, nil
}

// ListOrgIssues lists issues across organization repositories
func (m *GitHubAdminModule) ListOrgIssues(ctx context.Context, org string, state string, token string) (string, error) {
	args := []string{"gh", "issue", "list", "--search", "org:" + org}
	if state != "" {
		args = append(args, "--state", state)
	}
	args = append(args, "--json", "number,title,state,createdAt,repository")

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org issues: %w", err)
	}

	return output, nil
}

// ListOrgPRs lists pull requests across organization repositories
func (m *GitHubAdminModule) ListOrgPRs(ctx context.Context, org string, state string, token string) (string, error) {
	args := []string{"gh", "pr", "list", "--search", "org:" + org}
	if state != "" {
		args = append(args, "--state", state)
	}
	args = append(args, "--json", "number,title,state,createdAt,repository")

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org PRs: %w", err)
	}

	return output, nil
}

// GetVersion returns GitHub CLI version
func (m *GitHubAdminModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec([]string{"gh", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub CLI version: %w", err)
	}

	return output, nil
}

// ListOrgReposSimple lists repositories in an organization (MCP compatible)
func (m *GitHubAdminModule) ListOrgReposSimple(ctx context.Context, organization string, visibility string) (string, error) {
	args := []string{"gh", "repo", "list", organization}
	if visibility != "" {
		args = append(args, "--visibility", visibility)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org repos simple: %w", err)
	}

	return output, nil
}

// CreateOrgRepoSimple creates a repository in an organization (MCP compatible)
func (m *GitHubAdminModule) CreateOrgRepoSimple(ctx context.Context, organization string, repoName string, visibility string, description string) (string, error) {
	args := []string{"gh", "repo", "create", organization + "/" + repoName}
	
	if visibility != "" {
		args = append(args, "--"+visibility)
	}
	if description != "" {
		args = append(args, "--description", description)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create org repo simple: %w", err)
	}

	return output, nil
}

// GetRepoInfoSimple gets repository information (MCP compatible)
func (m *GitHubAdminModule) GetRepoInfoSimple(ctx context.Context, repository string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec([]string{"gh", "repo", "view", repository})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get repo info simple: %w", err)
	}

	return output, nil
}

// ListOrgIssuesSimple lists issues across organization repositories (MCP compatible)
func (m *GitHubAdminModule) ListOrgIssuesSimple(ctx context.Context, organization string, state string) (string, error) {
	args := []string{"gh", "issue", "list", "--search", "org:" + organization}
	if state != "" {
		args = append(args, "--state", state)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org issues simple: %w", err)
	}

	return output, nil
}

// ListOrgPRsSimple lists pull requests across organization repositories (MCP compatible)
func (m *GitHubAdminModule) ListOrgPRsSimple(ctx context.Context, organization string, state string) (string, error) {
	args := []string{"gh", "pr", "list", "--search", "org:" + organization}
	if state != "" {
		args = append(args, "--state", state)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list org PRs simple: %w", err)
	}

	return output, nil
}

// GetVersionSimple returns GitHub CLI version (MCP compatible)
func (m *GitHubAdminModule) GetVersionSimple(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "github-cli"}).
		WithExec([]string{"gh", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub CLI version simple: %w", err)
	}

	return output, nil
}
