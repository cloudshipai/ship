# History Scrub

Git history cleanup and sensitive data removal using industry-standard tools.

## Description

History Scrub provides comprehensive git repository history cleanup capabilities using proven tools like BFG Repo-Cleaner, git-filter-repo, and native git commands. It enables safe removal of sensitive data, large files, and unwanted content from git history while maintaining repository integrity. These tools are essential for security incidents, repository optimization, and compliance requirements.

## MCP Tools

### BFG Repo-Cleaner Tools
- **`history_scrub_bfg_remove_large_files`** - Remove large files from git history using BFG
- **`history_scrub_bfg_replace_text`** - Replace sensitive text in git history using BFG

### Git Filter-Repo Tools  
- **`history_scrub_filter_repo_remove_path`** - Remove files/paths using git-filter-repo

### Git Native Tools
- **`history_scrub_search_history`** - Search git history for sensitive patterns
- **`history_scrub_backup_repository`** - Create backup before cleanup operations
- **`history_scrub_git_cleanup`** - Run git cleanup after history rewriting

## Real CLI Commands Used

### BFG Repo-Cleaner Commands
- `bfg --strip-blobs-bigger-than <size> <repo.git>` - Remove large files
- `bfg --replace-text <replacements-file> <repo.git>` - Replace sensitive text
- `bfg --delete-files <pattern> <repo.git>` - Delete specific files
- `bfg --delete-folders <pattern> <repo.git>` - Delete specific folders

### Git Filter-Repo Commands
- `git filter-repo --path <path>` - Keep only specified paths
- `git filter-repo --path <path> --invert-paths` - Remove specified paths
- `git filter-repo --replace-text <file>` - Replace text patterns
- `git filter-repo --strip-blobs-bigger-than <size>` - Remove large files

### Git Native Commands
- `git log -S"<pattern>" --oneline --all` - Search history for patterns
- `git clone --bare <source> <backup>` - Create repository backup
- `git reflog expire --expire=now --all` - Clean reflog
- `git gc --prune=now --aggressive` - Garbage collection cleanup

## Use Cases

### Security Incident Response
- **Credential Leaks**: Remove accidentally committed passwords, API keys, tokens
- **Sensitive Data**: Clean up personally identifiable information (PII)
- **Configuration Files**: Remove files containing sensitive settings
- **Database Dumps**: Remove accidentally committed database exports

### Repository Optimization
- **Large Files**: Remove oversized binaries, datasets, or media files
- **Build Artifacts**: Clean up accidentally committed build outputs
- **Dependencies**: Remove large vendored dependencies
- **History Reduction**: Reduce repository size for faster clones

### Compliance and Governance
- **Data Protection**: Ensure compliance with privacy regulations
- **Corporate Policies**: Enforce data handling policies
- **Audit Requirements**: Clean history for compliance audits
- **Legal Requirements**: Remove content for legal compliance

### Repository Migration
- **Path Restructuring**: Reorganize directory structures
- **Project Extraction**: Extract subdirectories to new repositories
- **History Sanitization**: Clean history before open-sourcing
- **Vendor Transitions**: Prepare repositories for vendor changes

## Tool Capabilities

### BFG Repo-Cleaner Features
- **10-720x faster** than git-filter-branch
- **Simple interface** for common cleanup tasks
- **Safe operation** - protects HEAD commit by default
- **Comprehensive cleanup** of refs, commits, and blobs

### Git Filter-Repo Features
- **Official replacement** for git-filter-branch
- **Versatile rewriting** capabilities
- **Performance optimized** for large repositories
- **Advanced filtering** options and callbacks

### Git Native Features
- **Built-in commands** - no additional tools required
- **Pattern searching** with git log -S and -G
- **Reflog management** for complete cleanup
- **Garbage collection** for space reclamation

## Configuration Examples

### BFG Text Replacement File
```text
# replacements.txt format: old==>new
password123==>***REMOVED***
api_key_abc123==>***REMOVED***
secret_token_xyz==>***REMOVED***
database_password==>***REMOVED***
```

### Git Filter-Repo Usage
```bash
# Remove specific files
git filter-repo --path secrets.txt --invert-paths

# Keep only specific directory
git filter-repo --path src/

# Replace text patterns
git filter-repo --replace-text replacements.txt

# Remove large files
git filter-repo --strip-blobs-bigger-than 1M
```

### BFG Usage Examples
```bash
# Remove files larger than 100MB
bfg --strip-blobs-bigger-than 100M repo.git

# Replace sensitive text
bfg --replace-text passwords.txt repo.git

# Delete specific files
bfg --delete-files "*.log" repo.git

# Delete specific folders
bfg --delete-folders temp repo.git
```

## Workflow Recommendations

### Pre-Cleanup Preparation
1. **Create Backup**: Always backup repository before cleanup
2. **Coordinate Team**: Notify team members of upcoming changes
3. **Document Changes**: Record what will be removed and why
4. **Test Locally**: Practice cleanup on local clone first

### Cleanup Process
1. **Clone Fresh**: Start with fresh clone from origin
2. **Run Cleanup Tool**: Execute BFG or git-filter-repo
3. **Verify Results**: Check that sensitive data is removed
4. **Clean References**: Run git cleanup commands
5. **Force Push**: Update remote repository (requires coordination)

### Post-Cleanup Actions
1. **Verify Cleanup**: Confirm sensitive data removal
2. **Update Team**: Instruct team to re-clone repositories
3. **Update CI/CD**: Update any hardcoded repository references
4. **Monitor**: Watch for any issues or missed data

## Safety Considerations

### Backup Strategy
- **Multiple Backups**: Create backups at different stages
- **Bare Clones**: Use bare clones for complete backup
- **Remote Backups**: Store backups in separate locations
- **Verification**: Test backup restoration procedures

### Team Coordination
- **Communication**: Notify all team members before cleanup
- **Timing**: Schedule during low-activity periods
- **Re-cloning**: Ensure team re-clones after cleanup
- **Branch Management**: Consider impact on feature branches

### Repository Impact
- **Commit Hashes**: All commit hashes will change
- **Branch References**: Update any hardcoded commit references
- **CI/CD Pipelines**: Update build systems and deployment scripts
- **Integration Points**: Check third-party integrations

## Integration Patterns

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Security Scan Before Cleanup
  run: |
    git log -S"password" --oneline || echo "No sensitive patterns found"
    
- name: Create Backup
  run: |
    git clone --bare . ../backup.git
    
- name: Run BFG Cleanup
  run: |
    bfg --replace-text secrets.txt .
    git reflog expire --expire=now --all
    git gc --prune=now --aggressive
```

### Security Scanning Integration
```bash
# Pre-cleanup security scan
git log -S"api_key" --all --oneline
git log -S"password" --all --oneline
git log -G"secret.*=" --all --oneline

# Post-cleanup verification
git log --all --grep="password" --oneline
git log --all --grep="secret" --oneline
```

### Automated Workflows
```bash
#!/bin/bash
# automated-cleanup.sh

# 1. Create backup
git clone --bare "$REPO_URL" "backup-$(date +%Y%m%d).git"

# 2. Run cleanup
cd "$REPO_DIR"
bfg --replace-text replacements.txt .

# 3. Clean up git references
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# 4. Verify cleanup
if git log -S"password" --oneline; then
    echo "WARNING: Sensitive data may still exist"
    exit 1
fi

echo "Cleanup completed successfully"
```

## Best Practices

### Planning and Preparation
- **Assess Impact**: Understand full scope of changes needed
- **Test Approach**: Practice on repository copies first
- **Document Process**: Maintain detailed cleanup procedures
- **Schedule Properly**: Plan for team downtime and coordination

### Execution Guidelines
- **Start Fresh**: Always begin with fresh repository clone
- **Incremental Approach**: Make smaller, focused changes when possible
- **Verify Results**: Thoroughly check cleanup effectiveness
- **Complete Process**: Don't skip cleanup and garbage collection steps

### Communication and Training
- **Team Education**: Train team on prevention and response
- **Clear Procedures**: Document emergency cleanup procedures
- **Regular Audits**: Periodically scan for sensitive data
- **Prevention Focus**: Implement pre-commit hooks and scanning

History Scrub provides essential capabilities for maintaining secure and optimized git repositories through proven, industry-standard tools and practices.