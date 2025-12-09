package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// OpenVASModule provides integration with OpenVAS/Greenbone Vulnerability Management
// This module requires a running GVM stack and connects via gvm-tools
type OpenVASModule struct {
	client *dagger.Client
	name   string
}

// NewOpenVASModule creates a new OpenVAS module
func NewOpenVASModule(client *dagger.Client) *OpenVASModule {
	return &OpenVASModule{
		client: client,
		name:   "openvas",
	}
}

// getBaseContainer returns a container with gvm-tools installed
func (m *OpenVASModule) getBaseContainer() *dagger.Container {
	return m.client.Container().
		From(getImageTag("openvas", "registry.community.greenbone.net/community/gvm-tools:latest"))
}

// GetVersion returns the version of gvm-tools
func (m *OpenVASModule) GetVersion(ctx context.Context) (string, error) {
	container := m.getBaseContainer().
		WithExec([]string{"gvm-cli", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "OpenVAS/GVM Tools", nil
}

// ListScans lists all scans from a GVM instance
// Requires connection to a running GVM service
func (m *OpenVASModule) ListScans(ctx context.Context, socketPath string, username string, password string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}

	// GMP XML command to get tasks (scans)
	gmpCommand := `<get_tasks/>`

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer()

	// Mount the GVM socket if provided
	if socketPath != "" && socketPath != "/run/gvmd/gvmd.sock" {
		container = container.WithUnixSocket(socketPath, m.client.Host().UnixSocket(socketPath))
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		if strings.Contains(stderr, "Connection refused") || strings.Contains(stderr, "No such file") {
			return "", fmt.Errorf("cannot connect to GVM service. Ensure GVM is running and socket path is correct: %s", stderr)
		}
		return stderr, nil
	}

	return `{"message": "No scans found"}`, nil
}

// ListTargets lists all targets from a GVM instance
func (m *OpenVASModule) ListTargets(ctx context.Context, socketPath string, username string, password string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}

	gmpCommand := `<get_targets/>`

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "No targets found"}`, nil
}

// CreateTarget creates a new target for scanning
func (m *OpenVASModule) CreateTarget(ctx context.Context, socketPath string, username string, password string, targetName string, hosts string, portList string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if targetName == "" {
		return "", fmt.Errorf("target name is required")
	}
	if hosts == "" {
		return "", fmt.Errorf("hosts is required")
	}

	// Use default port list if not specified
	portListRef := ""
	if portList != "" {
		portListRef = fmt.Sprintf(`<port_list id="%s"/>`, portList)
	}

	gmpCommand := fmt.Sprintf(`<create_target><name>%s</name><hosts>%s</hosts>%s</create_target>`,
		targetName, hosts, portListRef)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return "", fmt.Errorf("failed to create target: %s", stderr)
	}

	return `{"message": "Target created successfully"}`, nil
}

// CreateScan creates a new scan task
func (m *OpenVASModule) CreateScan(ctx context.Context, socketPath string, username string, password string, scanName string, targetID string, configID string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if scanName == "" {
		return "", fmt.Errorf("scan name is required")
	}
	if targetID == "" {
		return "", fmt.Errorf("target ID is required")
	}

	// Use Full and fast scan config if not specified
	if configID == "" {
		configID = "daba56c8-73ec-11df-a475-002264764cea" // Full and fast
	}

	gmpCommand := fmt.Sprintf(`<create_task><name>%s</name><target id="%s"/><config id="%s"/></create_task>`,
		scanName, targetID, configID)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return "", fmt.Errorf("failed to create scan: %s", stderr)
	}

	return `{"message": "Scan created successfully"}`, nil
}

// StartScan starts an existing scan task
func (m *OpenVASModule) StartScan(ctx context.Context, socketPath string, username string, password string, taskID string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if taskID == "" {
		return "", fmt.Errorf("task ID is required")
	}

	gmpCommand := fmt.Sprintf(`<start_task task_id="%s"/>`, taskID)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "Scan started"}`, nil
}

// GetScanStatus gets the status of a scan
func (m *OpenVASModule) GetScanStatus(ctx context.Context, socketPath string, username string, password string, taskID string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if taskID == "" {
		return "", fmt.Errorf("task ID is required")
	}

	gmpCommand := fmt.Sprintf(`<get_tasks task_id="%s"/>`, taskID)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "Task not found"}`, nil
}

// GetScanResults gets the results of a completed scan
func (m *OpenVASModule) GetScanResults(ctx context.Context, socketPath string, username string, password string, reportID string, minSeverity string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if reportID == "" {
		return "", fmt.Errorf("report ID is required")
	}

	// Build filter for minimum severity if specified
	filter := ""
	if minSeverity != "" {
		filter = fmt.Sprintf(` filter="severity>%s"`, minSeverity)
	}

	gmpCommand := fmt.Sprintf(`<get_results report_id="%s"%s/>`, reportID, filter)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "No results found"}`, nil
}

// ListScanConfigs lists available scan configurations
func (m *OpenVASModule) ListScanConfigs(ctx context.Context, socketPath string, username string, password string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}

	gmpCommand := `<get_configs/>`

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "No scan configs found"}`, nil
}

// ListPortLists lists available port lists
func (m *OpenVASModule) ListPortLists(ctx context.Context, socketPath string, username string, password string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}

	gmpCommand := `<get_port_lists/>`

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "No port lists found"}`, nil
}

// StopScan stops a running scan
func (m *OpenVASModule) StopScan(ctx context.Context, socketPath string, username string, password string, taskID string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if taskID == "" {
		return "", fmt.Errorf("task ID is required")
	}

	gmpCommand := fmt.Sprintf(`<stop_task task_id="%s"/>`, taskID)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "Scan stopped"}`, nil
}

// DeleteScan deletes a scan task
func (m *OpenVASModule) DeleteScan(ctx context.Context, socketPath string, username string, password string, taskID string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}
	if taskID == "" {
		return "", fmt.Errorf("task ID is required")
	}

	gmpCommand := fmt.Sprintf(`<delete_task task_id="%s"/>`, taskID)

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return `{"message": "Scan deleted"}`, nil
}

// GetSystemStatus gets the overall system status of the GVM instance
func (m *OpenVASModule) GetSystemStatus(ctx context.Context, socketPath string, username string, password string) (string, error) {
	if socketPath == "" {
		socketPath = "/run/gvmd/gvmd.sock"
	}

	gmpCommand := `<get_version/>`

	args := []string{
		"gvm-cli", "socket",
		"--socketpath", socketPath,
		"--gmp-username", username,
		"--gmp-password", password,
		"--xml", gmpCommand,
	}

	container := m.getBaseContainer().
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return "", fmt.Errorf("cannot connect to GVM: %s", stderr)
	}

	return `{"message": "System status unavailable"}`, nil
}
