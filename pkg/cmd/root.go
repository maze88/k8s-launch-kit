package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log"

	applog "github.com/nvidia/k8s-launch-kit/pkg/log"
	"github.com/nvidia/k8s-launch-kit/pkg/profiles"
)

var (
	logLevel              string
	profile               string
	saveDeploymentFiles   string
	deploy                bool
	kubeconfig            string
	userConfig            string
	discoverClusterConfig string
	logger                = log.Log.WithName("l8k")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "l8k",
	Short: "Network Operator deployment and configuration tool",
	Long: `l8k is a CLI tool for deploying and managing NVIDIA Network Operator on Kubernetes.

The tool operates in 3 phases:

1. DISCOVER CLUSTER CONFIG: Deploy a thin profile of the Network Operator to discover 
   the cluster configuration and capabilities. This phase can be skipped if you provide 
   your own configuration with --user-config.

2. GENERATE DEPLOYMENT FILES: Based on the discovered or provided configuration, 
   generate a complete set of YAML deployment files for the selected network profile. 
   Files can be saved to disk using --save-deployment-files.

3. DEPLOY TO CLUSTER: Apply the generated deployment files to your Kubernetes cluster. 
   This phase requires --kubeconfig and can be skipped if --deploy is not specified.

This tool helps you deploy network profiles and configure cluster settings for optimal 
network performance with SR-IOV, RDMA, and other networking technologies.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set log level if provided
		if logLevel != "" {
			if err := applog.SetLogLevel(logLevel); err != nil {
				logger.Error(err, "Failed to set log level", "level", logLevel)
				os.Exit(1)
			}
		}

		if err := validateFlags(); err != nil {
			logger.Error(err, "Invalid command line arguments")
			os.Exit(1)
		}

		if saveDeploymentFiles != "" {
			logger.Info("Deployment files will be saved", "directory", saveDeploymentFiles)
		}

		if deploy {
			logger.Info("Deploy flag is enabled")
		}

		if kubeconfig != "" {
			logger.Info("Using kubeconfig", "path", kubeconfig)
		}

		if userConfig != "" {
			logger.Info("Using user config", "path", userConfig)
		}

		if discoverClusterConfig != "" {
			logger.Info("Will discover cluster config", "outputPath", discoverClusterConfig)
		}

		logger.Info("Welcome to Network Operator CLI!")
		fmt.Println("Use --help to see available options.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Phase 1: Cluster discovery flags
	rootCmd.Flags().StringVar(&discoverClusterConfig, "discover-cluster-config", "", "Phase 1: Deploy a thin Network Operator profile to discover cluster capabilities and save configuration to the specified path")
	rootCmd.Flags().StringVar(&userConfig, "user-config", "", "Phase 1: Use provided cluster configuration file instead of auto-discovery (skips Phase 1)")

	// Phase 2: Deployment generation flags
	rootCmd.Flags().StringVar(&profile, "profile", "", "Phase 2: Select the network profile to generate deployment files for ("+profiles.GetProfilesString()+")")
	rootCmd.Flags().StringVar(&saveDeploymentFiles, "save-deployment-files", "", "Phase 2: Save generated deployment files to the specified directory")

	// Phase 3: Cluster deployment flags
	rootCmd.Flags().BoolVar(&deploy, "deploy", false, "Phase 3: Deploy the generated files to the Kubernetes cluster")
	rootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Phase 3: Path to kubeconfig file for cluster deployment (required when using --deploy)")
	// Log level flag
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
}

// validateFlags validates the command line flag combinations
func validateFlags() error {
	// Rule 1: Either user-config or discover-cluster-config should be provided
	if userConfig == "" && discoverClusterConfig == "" {
		return fmt.Errorf("either --user-config or --discover-cluster-config must be provided")
	}

	// Rule 1b: Both user-config and discover-cluster-config cannot be provided together
	if userConfig != "" && discoverClusterConfig != "" {
		return fmt.Errorf("--user-config and --discover-cluster-config cannot be used together")
	}

	// Rule 2: if profile is selected, either save-deployment-files or deploy options should be provided
	if profile != "" && saveDeploymentFiles == "" && !deploy {
		return fmt.Errorf("when --profile is specified, either --save-deployment-files or --deploy must be provided")
	}

	// Rule 3: save-deployment-files or deploy can't work without profile
	if profile == "" && (saveDeploymentFiles != "" || deploy) {
		return fmt.Errorf("--save-deployment-files and --deploy require --profile to be specified")
	}

	// Rule 4: if deploy is provided, kubeconfig should be too
	if deploy && kubeconfig == "" {
		return fmt.Errorf("--deploy requires --kubeconfig to be specified")
	}

	if profile != "" {
		if !profiles.IsValidProfile(profile) {
			logger.Error(fmt.Errorf("invalid profile"), "Invalid profile provided", "profile", profile, "availableProfiles", profiles.AvailableProfiles)
			os.Exit(1)
		}
		logger.Info("Using profile", "profile", profile)
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Initialize logging
	applog.InitLog()

	// Implementation for config initialization
	// This can be expanded later to read from config files
}
