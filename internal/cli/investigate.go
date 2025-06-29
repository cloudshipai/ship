package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	investigateEnv      string
	investigateProvider string
	investigateAI       bool
)

var investigateCmd = &cobra.Command{
	Use:   "investigate",
	Short: "Run infrastructure investigation using Steampipe",
	Long:  `Run automated cloud infrastructure investigations using Steampipe and Dagger`,
	RunE:  runInvestigate,
}

func init() {
	rootCmd.AddCommand(investigateCmd)

	investigateCmd.Flags().StringVar(&investigateEnv, "env", "prod", "Environment to investigate")
	investigateCmd.Flags().StringVar(&investigateProvider, "provider", "", "Cloud provider (aws, cloudflare, heroku)")
	investigateCmd.Flags().BoolVar(&investigateAI, "ai", false, "Use AI for goal mapping")

	investigateCmd.MarkFlagRequired("provider")
}

func runInvestigate(cmd *cobra.Command, args []string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create pipeline
	pipeline := dagger.NewPipeline(engine, fmt.Sprintf("investigate-%s", investigateProvider))

	// For now, use hardcoded queries
	// In production, these would come from the goals API
	queries := getDefaultQueries(investigateProvider)

	fmt.Printf("Running %s investigation...\n", investigateProvider)

	var results map[string]string

	switch investigateProvider {
	case "aws":
		results, err = pipeline.InvestigateAWS(queries)
	case "cloudflare":
		results, err = pipeline.InvestigateCloudflare(queries)
	case "heroku":
		results, err = pipeline.InvestigateHeroku(queries)
	default:
		return fmt.Errorf("unsupported provider: %s", investigateProvider)
	}

	if err != nil {
		return fmt.Errorf("investigation failed: %w", err)
	}

	// Display results
	green := color.New(color.FgGreen)
	green.Printf("\n✓ Investigation completed!\n\n")

	fmt.Printf("Results:\n")
	for queryID, result := range results {
		fmt.Printf("\n=== %s ===\n", queryID)
		if len(result) > 500 {
			fmt.Printf("%s... (truncated)\n", result[:500])
		} else {
			fmt.Println(result)
		}
	}

	fmt.Printf("\nNext steps:\n")
	fmt.Println("1. Results will be automatically pushed to Cloudship")
	fmt.Println("2. Check the Cloudship UI for detailed analysis")
	fmt.Println("3. Review cost optimization and security recommendations")

	return nil
}

func getDefaultQueries(provider string) []string {
	switch provider {
	case "aws":
		return []string{
			"SELECT region, count(*) as instance_count FROM aws_ec2_instance GROUP BY region",
			"SELECT instance_id, instance_type, instance_state FROM aws_ec2_instance WHERE instance_state = 'running' LIMIT 10",
		}
	case "cloudflare":
		return []string{
			"SELECT name, status FROM cloudflare_zone LIMIT 10",
			"SELECT zone_name, type, name FROM cloudflare_dns_record LIMIT 10",
		}
	case "heroku":
		return []string{
			"SELECT name, region FROM heroku_app LIMIT 10",
			"SELECT app_name, dyno_type FROM heroku_dyno LIMIT 10",
		}
	default:
		return []string{}
	}
}
