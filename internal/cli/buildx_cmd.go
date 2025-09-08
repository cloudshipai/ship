package cli

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/spf13/cobra"
)

var buildxCmd = &cobra.Command{
	Use:   "buildx",
	Short: "Docker BuildX multi-platform image building",
	Long: `Docker BuildX tools for building and publishing multi-platform container images.

BuildX provides advanced features like:
- Multi-platform image building (linux/amd64, linux/arm64, etc.)
- Build cache optimization
- Registry publishing with authentication
- Development environments with BuildX pre-installed

Examples:
  # Build a multi-platform image
  ship buildx build --src-dir ./myapp --tag myapp:latest --platform linux/amd64,linux/arm64

  # Build and publish to a registry
  ship buildx publish --src-dir ./myapp --tag ghcr.io/user/myapp:latest --username user --password token

  # Set up a development environment
  ship buildx dev --src-dir ./myapp`,
}

var buildxBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an OCI image using Docker BuildX",
	Long: `Build a container image using Docker BuildX with multi-platform support.

This command builds a container image from a Dockerfile using Docker BuildX,
which supports building for multiple architectures simultaneously.

Examples:
  # Build for default platform (linux/amd64)
  ship buildx build --src-dir ./myapp --tag myapp:latest

  # Build for multiple platforms
  ship buildx build --src-dir ./myapp --tag myapp:latest --platform linux/amd64,linux/arm64

  # Build with custom Dockerfile location
  ship buildx build --src-dir ./myapp --tag myapp:latest --dockerfile-path docker/Dockerfile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srcDir, _ := cmd.Flags().GetString("src-dir")
		tag, _ := cmd.Flags().GetString("tag")
		platform, _ := cmd.Flags().GetString("platform")
		dockerfilePath, _ := cmd.Flags().GetString("dockerfile-path")

		if srcDir == "" {
			return fmt.Errorf("--src-dir is required")
		}
		if tag == "" {
			return fmt.Errorf("--tag is required")
		}

		ctx := context.Background()
		engine, err := dagger.NewEngine(ctx)
		if err != nil {
			return fmt.Errorf("failed to create dagger engine: %w", err)
		}
		defer engine.Close()

		buildxModule := modules.NewBuildXModule(engine.GetClient())
		result, err := buildxModule.Build(ctx, srcDir, tag, platform, dockerfilePath)
		if err != nil {
			return fmt.Errorf("buildx build failed: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

var buildxPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Build and publish an OCI image to a container registry",
	Long: `Build a container image with Docker BuildX and publish it to a container registry.

This command builds the image and pushes it directly to the specified registry,
supporting multi-platform builds and various authentication methods.

Examples:
  # Publish to Docker Hub
  ship buildx publish --src-dir ./myapp --tag myapp:latest --username user --password pass

  # Publish to GitHub Container Registry
  ship buildx publish --src-dir ./myapp --tag ghcr.io/user/myapp:latest --username user --password token --registry ghcr.io

  # Multi-platform publish
  ship buildx publish --src-dir ./myapp --tag myapp:latest --platform linux/amd64,linux/arm64 --username user --password pass`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srcDir, _ := cmd.Flags().GetString("src-dir")
		tag, _ := cmd.Flags().GetString("tag")
		platform, _ := cmd.Flags().GetString("platform")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		registry, _ := cmd.Flags().GetString("registry")
		dockerfilePath, _ := cmd.Flags().GetString("dockerfile-path")

		if srcDir == "" {
			return fmt.Errorf("--src-dir is required")
		}
		if tag == "" {
			return fmt.Errorf("--tag is required")
		}
		if username == "" {
			return fmt.Errorf("--username is required")
		}
		if password == "" {
			return fmt.Errorf("--password is required")
		}

		ctx := context.Background()
		engine, err := dagger.NewEngine(ctx)
		if err != nil {
			return fmt.Errorf("failed to create dagger engine: %w", err)
		}
		defer engine.Close()

		buildxModule := modules.NewBuildXModule(engine.GetClient())
		result, err := buildxModule.Publish(ctx, srcDir, tag, platform, username, password, registry, dockerfilePath)
		if err != nil {
			return fmt.Errorf("buildx publish failed: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

var buildxDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Set up a development environment with Docker BuildX",
	Long: `Create a development environment container with Docker BuildX pre-installed.

This command sets up a containerized development environment with BuildX ready to use,
optionally mounting your source code for interactive development.

Examples:
  # Basic dev environment
  ship buildx dev

  # Dev environment with source code mounted
  ship buildx dev --src-dir ./myapp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srcDir, _ := cmd.Flags().GetString("src-dir")

		ctx := context.Background()
		engine, err := dagger.NewEngine(ctx)
		if err != nil {
			return fmt.Errorf("failed to create dagger engine: %w", err)
		}
		defer engine.Close()

		buildxModule := modules.NewBuildXModule(engine.GetClient())
		result, err := buildxModule.Dev(ctx, srcDir)
		if err != nil {
			return fmt.Errorf("buildx dev setup failed: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

var buildxVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Docker BuildX version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		engine, err := dagger.NewEngine(ctx)
		if err != nil {
			return fmt.Errorf("failed to create dagger engine: %w", err)
		}
		defer engine.Close()

		buildxModule := modules.NewBuildXModule(engine.GetClient())
		result, err := buildxModule.GetVersion(ctx)
		if err != nil {
			return fmt.Errorf("failed to get buildx version: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

func init() {
	// Add subcommands to buildx
	buildxCmd.AddCommand(buildxBuildCmd)
	buildxCmd.AddCommand(buildxPublishCmd)
	buildxCmd.AddCommand(buildxDevCmd)
	buildxCmd.AddCommand(buildxVersionCmd)

	// Build command flags
	buildxBuildCmd.Flags().String("src-dir", "", "Source directory path containing Dockerfile (required)")
	buildxBuildCmd.Flags().String("tag", "", "Image tag (e.g., myapp:latest) (required)")
	buildxBuildCmd.Flags().String("platform", "linux/amd64", "Target platform(s) (can be comma-separated)")
	buildxBuildCmd.Flags().String("dockerfile-path", ".", "Path to Dockerfile relative to src-dir")

	// Publish command flags
	buildxPublishCmd.Flags().String("src-dir", "", "Source directory path containing Dockerfile (required)")
	buildxPublishCmd.Flags().String("tag", "", "Image tag to publish (required)")
	buildxPublishCmd.Flags().String("platform", "linux/amd64", "Target platform(s) (can be comma-separated)")
	buildxPublishCmd.Flags().String("username", "", "Registry username (required)")
	buildxPublishCmd.Flags().String("password", "", "Registry password or token (required)")
	buildxPublishCmd.Flags().String("registry", "docker.io", "Container registry URL")
	buildxPublishCmd.Flags().String("dockerfile-path", ".", "Path to Dockerfile relative to src-dir")

	// Dev command flags
	buildxDevCmd.Flags().String("src-dir", "", "Source directory to mount (optional)")

	// Add buildx to root command
	rootCmd.AddCommand(buildxCmd)
}