package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

func TestHostDeviceRdmaTemplate(t *testing.T) {
	// Define test config matching l8k-config.yaml structure
	config := Config{
		NetworkOperator: NetworkOperatorConfig{
			Version:          "v25.7.0",
			ComponentVersion: "network-operator-v25.7.0",
			Repository:       "nvcr.io/nvstaging/mellanox",
			Namespace:        "network-operator",
		},
		NvIpam: NvIpamConfig{
			PoolName:     "nv-ipam-pool",
			Subnet:       "192.168.2.0/24",
			Gateway:      "192.168.2.1",
			SubnetOffset: "0.0.1.0",
		},
		Hostdev: HostdevConfig{
			ResourceName: "hostdev-resource",
			NetworkName:  "hostdev-network",
		},
	}

	// Read the template file
	templatePath := filepath.Join("..", "..", "profiles", "host-device-rdma", "10-nicclusterpolicy.yaml")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	// Parse and execute the template
	tmpl, err := template.New("nicclusterpolicy").Parse(string(templateContent))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Verify that template variables were replaced
	result := buf.String()

	// Check that no template variables remain unreplaced
	if bytes.Contains([]byte(result), []byte("{{.")) {
		t.Errorf("Template contains unreplaced variables")
	}

	// Check that our test values are present
	if !bytes.Contains([]byte(result), []byte(config.NetworkOperator.Repository)) {
		t.Errorf("Network Operator repository not found in rendered template")
	}

	if !bytes.Contains([]byte(result), []byte(config.Hostdev.ResourceName)) {
		t.Errorf("Hostdev resource name not found in rendered template")
	}

	t.Logf("Template rendered successfully:\n%s", result)
}
