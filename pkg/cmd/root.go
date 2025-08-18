package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log"

	applog "github.com/nvidia/k8s-launch-kit/pkg/log"
)

var (
	// Used for flags.
	cfgFile               string
	logLevel              string
	profile               string
	saveDeploymentFiles   string
	deploy                bool
	kubeconfig            string
	useClusterConfig      string
	discoverClusterConfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "l8k",
	Short: "Network Operator deployment and configuration tool",
	Long: `l8k is a CLI tool for deploying and managing NVIDIA Network Operator on Kubernetes.
This tool helps you deploy network profiles, generate deployment files, and configure
cluster settings for optimal network performance with SR-IOV, RDMA, and other
networking technologies.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.Log.WithName("l8k")

		// Set log level if provided
		if logLevel != "" {
			if err := applog.SetLogLevel(logLevel); err != nil {
				logger.Error(err, "Failed to set log level", "level", logLevel)
				os.Exit(1)
			}
		}

		// Validate profile if provided
		if profile != "" {
			validProfiles := []string{"hostdevice", "sriov-rdma", "macvlan-rdma"}
			valid := false
			for _, p := range validProfiles {
				if profile == p {
					valid = true
					break
				}
			}
			if !valid {
				logger.Error(fmt.Errorf("invalid profile"), "Invalid profile provided", "profile", profile, "validProfiles", validProfiles)
				os.Exit(1)
			}
			logger.Info("Using profile", "profile", profile)
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

		if useClusterConfig != "" {
			logger.Info("Using cluster config", "path", useClusterConfig)
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

	// Network profile flag
	rootCmd.Flags().StringVar(&profile, "profile", "", "Select the network profile to deploy (hostdevice, sriov-rdma, macvlan-rdma)")

	// Deployment files flag
	rootCmd.Flags().StringVar(&saveDeploymentFiles, "save-deployment-files", "", "Specify the path to directory to save the generated deployment files to")

	// Deploy flag
	rootCmd.Flags().BoolVar(&deploy, "deploy", false, "Deploy the files after generating")

	// Kubeconfig flag
	rootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Specify the path to kubeconfig for the K8s cluster")

	// Use cluster config flag
	rootCmd.Flags().StringVar(&useClusterConfig, "use-cluster-config", "", "Specify the path to the cluster config. Skips the discovery stage before the deployment")

	// Discover cluster config flag
	rootCmd.Flags().StringVar(&discoverClusterConfig, "discover-cluster-config", "", "Deploy a thin profile of the Network Operator to discover the cluster configuration and save it as file to the specified path")

	// Log level flag
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

	// Config flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Use the provided config file. If not provided, the tool will try to discover the cluster config.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Initialize logging
	applog.InitLog()

	// Implementation for config initialization
	// This can be expanded later to read from config files
}
