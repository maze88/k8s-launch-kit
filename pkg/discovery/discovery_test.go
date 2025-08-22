package discovery

import (
	"context"
	"testing"

	netop "github.com/Mellanox/network-operator/api/v1alpha1"
	nicop "github.com/Mellanox/nic-configuration-operator/api/v1alpha1"
	"github.com/nvidia/k8s-launch-kit/pkg/config"
	"github.com/nvidia/k8s-launch-kit/pkg/netophelper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestEnsureNicClusterPolicy_CreateWhenNoneExists(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = netop.AddToScheme(scheme)
	_ = nicop.AddToScheme(scheme)

	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	ctx := context.Background()

	defaults := &config.NetworkOperatorConfig{
		Repository:       "nvcr.io/nvidia",
		ComponentVersion: "v1.2.3",
	}

	if err := netophelper.EnsureNicClusterPolicy(ctx, c, defaults); err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}

	// Verify one created
	list := &netop.NicClusterPolicyList{}
	if err := c.List(ctx, list); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(list.Items))
	}
}

func TestEnsureNicClusterPolicy_FailsWhenExists(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = netop.AddToScheme(scheme)
	_ = nicop.AddToScheme(scheme)

	existing := &netop.NicClusterPolicy{ObjectMeta: metav1.ObjectMeta{Name: "existing"}}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	ctx := context.Background()

	defaults := &config.NetworkOperatorConfig{Repository: "repo", ComponentVersion: "v"}

	if err := netophelper.EnsureNicClusterPolicy(ctx, c, defaults); err == nil {
		t.Fatalf("expected error when policy exists, got nil")
	}
}

func TestBuildClusterConfigFromNicDevices(t *testing.T) {
	devices := []nicop.NicDevice{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "dev1"},
			Status: nicop.NicDeviceStatus{
				Ports: []nicop.NicDevicePortSpec{
					{PCI: "0000:08:00.0", NetworkInterface: "eth0", RdmaInterface: "mlx5_0"},
					{PCI: "0000:08:00.1", NetworkInterface: "eth1", RdmaInterface: "mlx5_1"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "dev2"},
			Status: nicop.NicDeviceStatus{
				Ports: []nicop.NicDevicePortSpec{
					{PCI: "0000:08:00.0", NetworkInterface: "eth0", RdmaInterface: "mlx5_0"}, // dup
					{PCI: "0000:3b:00.0", NetworkInterface: "ib0", RdmaInterface: "mlx5_2"},
				},
			},
		},
	}

	cfg := buildClusterConfigFromNicDevices(devices)

	wantPCI := []string{"0000:08:00.0", "0000:08:00.1", "0000:3b:00.0"}
	wantNet := []string{"eth0", "eth1", "ib0"}
	wantRdma := []string{"mlx5_0", "mlx5_1", "mlx5_2"}

	if len(cfg.NvidiaNICs.PF.PciAddresses) != len(wantPCI) {
		t.Fatalf("unexpected pci count: %v", cfg.NvidiaNICs.PF.PciAddresses)
	}
	for i, v := range wantPCI {
		if cfg.NvidiaNICs.PF.PciAddresses[i] != v {
			t.Fatalf("unexpected pci[%d]=%q", i, cfg.NvidiaNICs.PF.PciAddresses[i])
		}
	}
	if len(cfg.NvidiaNICs.PF.NetworkInterfaces) != len(wantNet) {
		t.Fatalf("unexpected net if count: %v", cfg.NvidiaNICs.PF.NetworkInterfaces)
	}
	for i, v := range wantNet {
		if cfg.NvidiaNICs.PF.NetworkInterfaces[i] != v {
			t.Fatalf("unexpected net[%d]=%q", i, cfg.NvidiaNICs.PF.NetworkInterfaces[i])
		}
	}
	if len(cfg.NvidiaNICs.PF.RdmaDevices) != len(wantRdma) {
		t.Fatalf("unexpected rdma count: %v", cfg.NvidiaNICs.PF.RdmaDevices)
	}
	for i, v := range wantRdma {
		if cfg.NvidiaNICs.PF.RdmaDevices[i] != v {
			t.Fatalf("unexpected rdma[%d]=%q", i, cfg.NvidiaNICs.PF.RdmaDevices[i])
		}
	}
}
