# l8k - Network Operator CLI

l8k is a CLI tool for deploying and managing NVIDIA Network Operator on Kubernetes. The tool operates in 3 distinct phases to provide flexible deployment workflows for optimal network performance with SR-IOV, RDMA, and other networking technologies.

## How l8k Works - 3 Phase Operation

### Phase 1: Discover Cluster Configuration
Deploy a minimal Network Operator profile to automatically discover your cluster's network capabilities and hardware configuration. This phase can be skipped if you provide your own configuration file.

### Phase 2: Generate Deployment Files
Based on the discovered or provided configuration, generate a complete set of YAML deployment files tailored to your selected network profile. Files can be saved to disk for review or version control.

### Phase 3: Deploy to Cluster
Apply the generated deployment files directly to your Kubernetes cluster. This phase requires cluster access and can be skipped if you only want to generate files.

## Features

- Deploy different network profiles (hostdevice, sriov-rdma, macvlan-rdma)
- Generate and save deployment files
- Automatic deployment capability
- Cluster configuration discovery
- Structured logging with configurable levels
- Docker support for containerized deployments

## Installation

### Build from source

```bash
git clone <repository-url>
cd launch-kubernetes
make build
```

The binary will be available at `build/l8k`.

### Docker

Build the Docker image:

```bash
make docker-build
```

## Usage

<!-- BEGIN HELP -->
<!-- This section is automatically updated by running 'make update-readme' -->

```
l8k is a CLI tool for deploying and managing NVIDIA Network Operator on Kubernetes.

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
network performance with SR-IOV, RDMA, and other networking technologies.

Usage:
  l8k [flags]
  l8k [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number

Flags:
      --deploy                           Phase 3: Deploy the generated files to the Kubernetes cluster
      --discover-cluster-config string   Phase 1: Deploy a thin Network Operator profile to discover cluster capabilities and save configuration to the specified path
  -h, --help                             help for l8k
      --kubeconfig string                Phase 3: Path to kubeconfig file for cluster deployment (required when using --deploy)
      --log-level string                 Log level (debug, info, warn, error) (default "info")
      --profile string                   Phase 2: Select the network profile to generate deployment files for (host-device-rdma, macvlan-rdma, sriov-rdma, test-profile)
      --save-deployment-files string     Phase 2: Save generated deployment files to the specified directory
      --user-config string               Phase 1: Use provided cluster configuration file instead of auto-discovery (skips Phase 1)

Use "l8k [command] --help" for more information about a command.
```

<!-- END HELP -->

## Usage Examples

### Complete 3-Phase Workflow

Discover cluster config, generate files, and deploy:

```bash
# All phases: discover → generate → deploy
l8k --discover-cluster-config ./cluster-config.yaml \
    --profile sriov-rdma \
    --save-deployment-files ./deployments \
    --deploy --kubeconfig ~/.kube/config
```

### Phase 1 Only: Discover Cluster Configuration

```bash
# Phase 1: Just discover cluster capabilities
l8k --discover-cluster-config ./my-cluster-config.yaml
```

### Skip Phase 1: Use Existing Configuration  

Generate and deploy with pre-existing config:

```bash
# Phase 2 & 3: Use existing config → generate → deploy
l8k --user-config ./existing-config.yaml \
    --profile host-device-rdma \
    --deploy --kubeconfig ~/.kube/config
```

### Phase 2 Only: Generate Deployment Files

```bash
# Phase 2: Generate files without deploying
l8k --user-config ./config.yaml \
    --profile macvlan-rdma \
    --save-deployment-files ./deployments
```

### Advanced Usage

Enable debug logging and save to custom directory:

```bash
l8k --user-config ./config.yaml \
    --profile sriov-rdma \
    --save-deployment-files /opt/deployments \
    --log-level debug
```

## Development

### Building

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make clean        # Clean build artifacts
```

### Testing

```bash
make test         # Run tests
make coverage     # Run tests with coverage
```

### Linting

```bash
make lint         # Run linter
make lint-check   # Install and run linter
```

### Docker

```bash
make docker-build # Build Docker image
make docker-run   # Run Docker container
```

## Contributing

1. Ensure you have Go 1.21+ installed
2. Run `make dev-setup` to install dependencies and run initial checks
3. Make your changes
4. Run `make lint` and `make test` to ensure code quality
5. Submit a pull request

## License

This project is licensed under the terms specified in the repository.
