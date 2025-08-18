# Container Registry MCP Tool

Container Registry tools provide Docker CLI operations for managing container images and registries.

## Description

Container Registry operations include:
- Pushing and pulling container images
- Tagging and managing image versions
- Registry authentication
- Image inspection and metadata management
- Multi-registry support

## MCP Functions

### `container_registry_login`
Login to container registry using Docker CLI.

**Parameters:**
- `registry` (required): Registry URL
- `username` (required): Username for authentication
- `password` (required): Password for authentication

**CLI Command:** `docker login <registry> -u <username> -p <password>`

### `container_registry_push`
Push container image to registry using Docker CLI.

**Parameters:**
- `image` (required): Image name with tag
- `registry`: Registry URL (if not in image name)

**CLI Command:** `docker push <image>`

### `container_registry_pull`
Pull container image from registry using Docker CLI.

**Parameters:**
- `image` (required): Image name with tag
- `platform`: Target platform (e.g., linux/amd64)

**CLI Command:** `docker pull <image> [--platform <platform>]`

### `container_registry_tag`
Tag container image using Docker CLI.

**Parameters:**
- `source_image` (required): Source image name
- `target_image` (required): Target image name with new tag

**CLI Command:** `docker tag <source_image> <target_image>`

### `container_registry_list_tags`
List tags for image in registry using Docker CLI and registry API.

**Parameters:**
- `image` (required): Image name
- `registry`: Registry URL

**CLI Command:** Uses Docker CLI with registry API calls

### `container_registry_delete`
Delete image from registry using Docker CLI.

**Parameters:**
- `image` (required): Image name with tag
- `force`: Force deletion

**CLI Command:** `docker image rm <image> [--force]`

## Common Use Cases

1. **Image Management**: Push, pull, and manage container images
2. **Version Control**: Tag and version container images
3. **Multi-Registry Support**: Work with different registries
4. **CI/CD Integration**: Automate image deployment
5. **Image Distribution**: Share images across teams

## Integration with Ship CLI

All Container Registry tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Manage container images across registries
- Automate image deployment workflows
- Handle registry authentication
- Perform image operations

The tools use containerized execution via Dagger for consistent, isolated registry operations.