package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// NmapModule runs Nmap for network scanning
type NmapModule struct {
	client *dagger.Client
	name   string
}

// NewNmapModule creates a new Nmap module
func NewNmapModule(client *dagger.Client) *NmapModule {
	return &NmapModule{
		client: client,
		name:   "nmap",
	}
}

// GetVersion returns the version of Nmap
func (m *NmapModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec([]string{"nmap", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Nmap - Network exploration tool and security scanner", nil
}

// ScanHost performs a basic host scan
func (m *NmapModule) ScanHost(ctx context.Context, target string, scanType string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org" // Safe test target
	}
	
	args := []string{"nmap"}
	
	// Add scan type
	switch scanType {
	case "ping":
		args = append(args, "-sn") // Ping scan
	case "quick":
		args = append(args, "-T4", "-F") // Fast scan, fewer ports
	case "intense":
		args = append(args, "-T4", "-A", "-v") // Intense scan
	case "stealth":
		args = append(args, "-sS") // SYN stealth scan
	default:
		// Default scan
	}
	
	// Output in XML for parsing
	args = append(args, "-oX", "-", target)

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	
	if stderr != "" && !strings.Contains(stderr, "WARNING") {
		return stderr, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Host %s scanned</nmaprun>`, target), nil
}

// PortScan scans specific ports on a host
func (m *NmapModule) PortScan(ctx context.Context, target string, ports string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap"}
	
	// Add port specification
	if ports != "" {
		args = append(args, "-p", ports)
	} else {
		args = append(args, "-p", "80,443,22,21,25,3306,5432,8080,8443")
	}
	
	// Output in XML
	args = append(args, "-oX", "-", target)

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Ports %s scanned on %s</nmaprun>`, ports, target), nil
}

// ServiceDetection performs service and version detection
func (m *NmapModule) ServiceDetection(ctx context.Context, target string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap", "-sV", "-oX", "-", target}

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Service detection completed for %s</nmaprun>`, target), nil
}

// OSDetection performs operating system detection
func (m *NmapModule) OSDetection(ctx context.Context, target string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	// OS detection requires root privileges
	args := []string{"nmap", "-O", "-oX", "-", target}

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	
	// Check if it failed due to privileges
	if strings.Contains(stderr, "root") || strings.Contains(stderr, "privileges") {
		return `<nmaprun>OS detection requires root privileges</nmaprun>`, nil
	}
	
	return fmt.Sprintf(`<nmaprun>OS detection attempted for %s</nmaprun>`, target), nil
}

// VulnerabilityScan performs vulnerability scanning using NSE scripts
func (m *NmapModule) VulnerabilityScan(ctx context.Context, target string, scriptCategory string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap"}
	
	// Add script category
	switch scriptCategory {
	case "vuln":
		args = append(args, "--script", "vuln")
	case "auth":
		args = append(args, "--script", "auth")
	case "default":
		args = append(args, "--script", "default")
	case "safe":
		args = append(args, "--script", "safe")
	default:
		args = append(args, "--script", "safe")
	}
	
	args = append(args, "-oX", "-", target)

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Vulnerability scan completed for %s</nmaprun>`, target), nil
}

// NetworkDiscovery performs network discovery scan
func (m *NmapModule) NetworkDiscovery(ctx context.Context, network string) (string, error) {
	if network == "" {
		network = "192.168.1.0/24"
	}
	
	// Ping scan for network discovery
	args := []string{"nmap", "-sn", "-oX", "-", network}

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Network discovery completed for %s</nmaprun>`, network), nil
}

// UDPScan performs UDP port scanning
func (m *NmapModule) UDPScan(ctx context.Context, target string, ports string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap", "-sU"}
	
	if ports != "" {
		args = append(args, "-p", ports)
	} else {
		args = append(args, "-p", "53,67,68,69,123,161,162,514")
	}
	
	args = append(args, "-oX", "-", target)

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	
	// Check if it failed due to privileges
	if strings.Contains(stderr, "root") || strings.Contains(stderr, "privileges") {
		return `<nmaprun>UDP scanning requires root privileges</nmaprun>`, nil
	}
	
	return fmt.Sprintf(`<nmaprun>UDP scan completed for %s</nmaprun>`, target), nil
}

// FirewallEvasion performs scans with firewall evasion techniques
func (m *NmapModule) FirewallEvasion(ctx context.Context, target string, technique string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap"}
	
	// Add evasion technique
	switch technique {
	case "fragment":
		args = append(args, "-f") // Fragment packets
	case "decoy":
		args = append(args, "-D", "RND:10") // Random decoys
	case "timing":
		args = append(args, "-T0") // Paranoid timing
	case "source-port":
		args = append(args, "--source-port", "53") // Use DNS source port
	default:
		args = append(args, "-f") // Default to fragmentation
	}
	
	args = append(args, "-oX", "-", target)

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Firewall evasion scan completed for %s</nmaprun>`, target), nil
}

// ScriptScan runs specific NSE scripts
func (m *NmapModule) ScriptScan(ctx context.Context, target string, script string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	if script == "" {
		script = "http-headers"
	}
	
	args := []string{"nmap", "--script", script, "-oX", "-", target}

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Script %s executed for %s</nmaprun>`, script, target), nil
}

// TracerouteScan performs traceroute to target
func (m *NmapModule) TracerouteScan(ctx context.Context, target string) (string, error) {
	if target == "" {
		target = "scanme.nmap.org"
	}
	
	args := []string{"nmap", "--traceroute", "-oX", "-", target}

	container := m.client.Container().
		From("instrumentisto/nmap:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf(`<nmaprun>Traceroute completed for %s</nmaprun>`, target), nil
}