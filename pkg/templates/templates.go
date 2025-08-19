package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Config represents the l8k-config.yaml structure
type Config struct {
	NetworkOperator NetworkOperatorConfig `yaml:"networkOperator"`
	NvIpam          NvIpamConfig          `yaml:"nvIpam"`
	Sriov           SriovConfig           `yaml:"sriov"`
	Hostdev         HostdevConfig         `yaml:"hostdev"`
	RdmaShared      RdmaSharedConfig      `yaml:"rdmaShared"`
	Ipoib           IpoibConfig           `yaml:"ipoib"`
	Macvlan         MacvlanConfig         `yaml:"macvlan"`
	ClusterConfig   ClusterConfigStruct   `yaml:"clusterConfig"`
}

type NetworkOperatorConfig struct {
	Version          string `yaml:"version"`
	ComponentVersion string `yaml:"componentVersion"`
	Repository       string `yaml:"repository"`
	Namespace        string `yaml:"namespace"`
}

type NvIpamConfig struct {
	PoolName     string `yaml:"poolName"`
	Subnet       string `yaml:"subnet"`
	Gateway      string `yaml:"gateway"`
	SubnetOffset string `yaml:"subnetOffset"`
}

type SriovConfig struct {
	Mtu          int    `yaml:"mtu"`
	NumVfs       int    `yaml:"numVfs"`
	Priority     int    `yaml:"priority"`
	ResourceName string `yaml:"resourceName"`
	NetworkName  string `yaml:"networkName"`
}

type HostdevConfig struct {
	ResourceName string `yaml:"resourceName"`
	NetworkName  string `yaml:"networkName"`
}

type RdmaSharedConfig struct {
	ResourceName string `yaml:"resourceName"`
	HcaMax       int    `yaml:"hcaMax"`
}

type IpoibConfig struct {
	NetworkName                string `yaml:"networkName"`
	SeparateNetworkPerDevice   bool   `yaml:"separateNetworkPerDevice"`
	SingleNetworkForAllDevices bool   `yaml:"singleNetworkForAllDevices"`
}

type MacvlanConfig struct {
	NetworkName                string `yaml:"networkName"`
	SeparateNetworkPerDevice   bool   `yaml:"separateNetworkPerDevice"`
	SingleNetworkForAllDevices bool   `yaml:"singleNetworkForAllDevices"`
}

type ClusterConfigStruct struct {
	Nodes      NodesConfig `yaml:"nodes"`
	NvidiaNICs interface{} `yaml:"nvidiaNICs"`
	PF         PFConfig    `yaml:"pf"`
}

type NodesConfig struct {
	Capabilities CapabilitiesConfig `yaml:"capabilities"`
}

type CapabilitiesConfig struct {
	Sriov bool `yaml:"sriov"`
	Rdma  bool `yaml:"rdma"`
	Ib    bool `yaml:"ib"`
}

type PFConfig struct {
	RdmaDevices       []string `yaml:"rdmaDevices"`
	PciAddresses      []string `yaml:"pciAddresses"`
	NetworkInterfaces []string `yaml:"networkInterfaces"`
}

// ProcessTemplate processes a Go template file with the given config
func ProcessTemplate(templatePath string, config Config) (string, error) {
	// Read the template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// ProcessProfileTemplates processes all template files in a profile directory
func ProcessProfileTemplates(profileDir string, config Config) (map[string]string, error) {
	results := make(map[string]string)

	// Find all YAML files in the profile directory
	entries, err := os.ReadDir(profileDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory %s: %w", profileDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			templatePath := filepath.Join(profileDir, filename)

			processed, err := ProcessTemplate(templatePath, config)
			if err != nil {
				return nil, fmt.Errorf("failed to process template %s: %w", filename, err)
			}

			results[filename] = processed
		}
	}

	return results, nil
}
