package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CheckSSLCertModule runs SSL certificate checking
type CheckSSLCertModule struct {
	client *dagger.Client
	name   string
}

const checkSSLCertBinary = "/usr/local/bin/check_ssl_cert"

// NewCheckSSLCertModule creates a new SSL certificate checker module
func NewCheckSSLCertModule(client *dagger.Client) *CheckSSLCertModule {
	return &CheckSSLCertModule{
		client: client,
		name:   "check-ssl-cert",
	}
}

// CheckCertificate checks SSL certificate for a host using check_ssl_cert
func (m *CheckSSLCertModule) CheckCertificate(ctx context.Context, host string, port int) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "openssl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"}).
		WithExec([]string{
			checkSSLCertBinary,
			"-H", host,
			"-p", fmt.Sprintf("%d", port),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check SSL certificate: %w", err)
	}

	return output, nil
}

// CheckCertificateExpiry checks certificate expiry
func (m *CheckSSLCertModule) CheckCertificateExpiry(ctx context.Context, host string, port int, warningDays int) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "openssl", "curl"}).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`
				cert_date=$(echo | openssl s_client -servername %s -connect %s:%d 2>/dev/null | openssl x509 -noout -enddate | cut -d= -f2)
				cert_epoch=$(date -d "$cert_date" +%%s)
				now_epoch=$(date +%%s)
				days_remaining=$(( (cert_epoch - now_epoch) / 86400 ))
				echo "Days remaining: $days_remaining"
				if [ $days_remaining -lt %d ]; then
					echo "WARNING: Certificate expires in $days_remaining days"
				fi
			`, host, host, port, warningDays),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check certificate expiry: %w", err)
	}

	return output, nil
}

// ValidateCertificateChain validates certificate chain
func (m *CheckSSLCertModule) ValidateCertificateChain(ctx context.Context, host string, port int) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "openssl", "ca-certificates"}).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf("echo | openssl s_client -servername %s -connect %s:%d -verify_return_error 2>/dev/null && echo 'Certificate chain is valid' || echo 'Certificate chain validation failed'", host, host, port),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate certificate chain: %w", err)
	}

	return output, nil
}

// CheckCertificateFromFile checks SSL certificate from a local file
func (m *CheckSSLCertModule) CheckCertificateFromFile(ctx context.Context, filePath string, warningDays int, criticalDays int, allowSelfSigned bool) (string, error) {
	args := []string{checkSSLCertBinary, "-f", "/cert"}
	if warningDays > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", warningDays))
	}
	if criticalDays > 0 {
		args = append(args, "-c", fmt.Sprintf("%d", criticalDays))
	}
	if allowSelfSigned {
		args = append(args, "-s")
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "openssl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"}).
		WithFile("/cert", m.client.Host().File(filePath)).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check SSL certificate from file: %w", err)
	}

	return output, nil
}

// CheckCertificateWithAdvancedOptions checks SSL certificate with advanced options
func (m *CheckSSLCertModule) CheckCertificateWithAdvancedOptions(ctx context.Context, host string, port int, protocol string, warningDays int, criticalDays int, allowSelfSigned bool, rootCert string, checkChain bool, ignoreAuth bool, timeout int, debug bool) (string, error) {
	args := []string{checkSSLCertBinary, "-H", host, "-p", fmt.Sprintf("%d", port)}
	
	if protocol != "" {
		args = append(args, "-P", protocol)
	}
	if warningDays > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", warningDays))
	}
	if criticalDays > 0 {
		args = append(args, "-c", fmt.Sprintf("%d", criticalDays))
	}
	if allowSelfSigned {
		args = append(args, "-s")
	}
	if rootCert != "" {
		args = append(args, "-r", "/rootcert")
	}
	if checkChain {
		args = append(args, "--check-chain")
	}
	if ignoreAuth {
		args = append(args, "-A")
	}
	if timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", timeout))
	}
	if debug {
		args = append(args, "-d")
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "openssl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"})
	
	if rootCert != "" {
		container = container.WithFile("/rootcert", m.client.Host().File(rootCert))
	}
	
	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check SSL certificate with advanced options: %w", err)
	}

	return output, nil
}

// CheckCertificateFingerprint checks SSL certificate fingerprint
func (m *CheckSSLCertModule) CheckCertificateFingerprint(ctx context.Context, host string, port int, expectedFingerprint string) (string, error) {
	args := []string{checkSSLCertBinary, "-H", host, "-p", fmt.Sprintf("%d", port), "--fingerprint", expectedFingerprint}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "openssl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check SSL certificate fingerprint: %w", err)
	}

	return output, nil
}

// CheckCertificateComprehensive performs comprehensive SSL certificate check with all options
func (m *CheckSSLCertModule) CheckCertificateComprehensive(ctx context.Context, host string, port int, timeout int, debug bool) (string, error) {
	args := []string{checkSSLCertBinary, "-H", host, "-p", fmt.Sprintf("%d", port), "--all"}
	
	if timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", timeout))
	}
	if debug {
		args = append(args, "-d")
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "openssl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run comprehensive SSL certificate check: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of check_ssl_cert
func (m *CheckSSLCertModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "bash"}).
		WithExec([]string{"sh", "-c", "curl -L https://raw.githubusercontent.com/matteocorti/check_ssl_cert/master/check_ssl_cert -o /usr/local/bin/check_ssl_cert && chmod +x /usr/local/bin/check_ssl_cert"}).
		WithExec([]string{checkSSLCertBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get check_ssl_cert version: %w", err)
	}

	return output, nil
}
