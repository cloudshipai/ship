# GitHub Admin MCP Tool

GitHub administration tools using the GitHub CLI (gh) for repository and organization management.

## Description

GitHub Admin provides comprehensive GitHub management capabilities:
- Repository creation and configuration
- Organization and team management
- Security settings and policies
- Workflow and action management
- API interactions and automation

## MCP Functions

### `github_admin_create_repo`
Create new GitHub repository using real gh CLI.

**Parameters:**
- `name` (required): Repository name
- `description`: Repository description
- `public`: Make repository public
- `private`: Make repository private
- `team`: Team to grant access
- `template`: Template repository to use
- `clone`: Clone after creation

**CLI Command:** `gh repo create <name> [--description <desc>] [--public|--private] [--team <team>] [--template <template>] [--clone]`

### `github_admin_manage_secrets`
Manage repository secrets using real gh CLI.

**Parameters:**
- `repo` (required): Repository name (owner/repo)
- `secret_name` (required): Secret name
- `secret_value`: Secret value (for set operation)
- `operation`: Operation (set, delete, list)

**CLI Command:** `gh secret <operation> <secret_name> [--repo <repo>] [--body <value>]`

### `github_admin_configure_settings`
Configure repository settings using real gh CLI.

**Parameters:**
- `repo` (required): Repository name (owner/repo)
- `default_branch`: Set default branch
- `description`: Update description
- `homepage`: Set homepage URL
- `topics`: Comma-separated topics
- `visibility`: Repository visibility (public, private)

**CLI Command:** `gh repo edit <repo> [--default-branch <branch>] [--description <desc>] [--homepage <url>] [--add-topic <topics>] [--visibility <visibility>]`

### `github_admin_manage_collaborators`
Manage repository collaborators using real gh CLI.

**Parameters:**
- `repo` (required): Repository name (owner/repo)
- `username` (required): GitHub username
- `permission`: Permission level (read, write, admin)
- `operation`: Operation (add, remove)

**CLI Command:** `gh api repos/<repo>/collaborators/<username> [--method PUT/DELETE] [-f permission=<permission>]`

### `github_admin_workflow_management`
Manage GitHub Actions workflows using real gh CLI.

**Parameters:**
- `repo` (required): Repository name (owner/repo)
- `workflow`: Workflow file name or ID
- `operation`: Operation (list, view, run, disable, enable)

**CLI Command:** `gh workflow <operation> [<workflow>] --repo <repo>`

### `github_admin_org_management`
Manage organization settings using real gh CLI.

**Parameters:**
- `org` (required): Organization name
- `operation`: Operation type (list-members, list-repos, list-teams)

**CLI Command:** `gh api orgs/<org>/<operation>`

## Common Use Cases

1. **Repository Management**: Create and configure repositories
2. **Security Management**: Handle secrets and permissions
3. **Team Collaboration**: Manage collaborators and teams
4. **Workflow Automation**: Control GitHub Actions
5. **Organization Administration**: Manage org-level settings

## Integration with Ship CLI

All GitHub Admin tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Automate repository management
- Configure security settings
- Manage team permissions
- Control GitHub Actions workflows

The tools use containerized execution via Dagger for consistent, isolated GitHub operations.