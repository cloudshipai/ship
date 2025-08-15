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
