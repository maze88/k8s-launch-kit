# l8k - Network Operator CLI

l8k is a CLI tool for deploying and managing NVIDIA Network Operator on Kubernetes. This tool helps you deploy network profiles, generate deployment files, and configure cluster settings for optimal network performance with SR-IOV, RDMA, and other networking technologies.

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
This tool helps you deploy network profiles, generate deployment files, and configure
cluster settings for optimal network performance with SR-IOV, RDMA, and other
networking technologies.

Usage:
  l8k [flags]
  l8k [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number

Flags:
      --config string                    Use the provided config file. If not provided, the tool will try to discover the cluster config.
      --deploy                           Deploy the files after generating
      --discover-cluster-config string   Deploy a thin profile of the Network Operator to discover the cluster configuration and save it as file to the specified path
  -h, --help                             help for l8k
      --kubeconfig string                Specify the path to kubeconfig for the K8s cluster
      --log-level string                 Log level (debug, info, warn, error) (default "info")
      --profile string                   Select the network profile to deploy (hostdevice, sriov-rdma, macvlan-rdma)
      --save-deployment-files string     Specify the path to directory to save the generated deployment files to
      --use-cluster-config string        Specify the path to the cluster config. Skips the discovery stage before the deployment

Use "l8k [command] --help" for more information about a command.
```

<!-- END HELP -->

## Examples

### Deploy with SR-IOV RDMA profile

```bash
l8k --profile sriov-rdma --deploy
```

### Generate deployment files without deploying

```bash
l8k --profile hostdevice --save-deployment-files ./deployments
```

### Discover cluster configuration

```bash
l8k --discover-cluster-config ./cluster-config.yaml
```

### Enable debug logging

```bash
l8k --log-level debug --profile macvlan-rdma
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
