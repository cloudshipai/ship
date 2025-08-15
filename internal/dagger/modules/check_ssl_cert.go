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

// NewCheckSSLCertModule creates a new SSL certificate checker module
func NewCheckSSLCertModule(client *dagger.Client) *CheckSSLCertModule {
	return &CheckSSLCertModule{
		client: client,
		name:   "check-ssl-cert",
	}
}

// CheckCertificate checks SSL certificate for a host
func (m *CheckSSLCertModule) CheckCertificate(ctx context.Context, host string, port int) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "openssl", "curl"}).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf("echo | openssl s_client -servername %s -connect %s:%d 2>/dev/null | openssl x509 -noout -dates -subject -issuer", host, host, port),
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
