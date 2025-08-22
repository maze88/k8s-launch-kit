package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/nvidia/k8s-launch-kit/pkg/kubeclient"
	"github.com/nvidia/k8s-launch-kit/pkg/netophelper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	yaml "sigs.k8s.io/yaml"
)

// Apply reads Kubernetes manifests from dirPath and applies them to the cluster.
// If a NicClusterPolicy is present, it is applied first and the function waits
// for it to become ready before applying the remaining manifests.
func Apply(ctx context.Context, kubeconfigPath, dirPath string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	c, err := kubeclient.New(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}

	// List files in directory (non-recursive) and sort
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	filePaths := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dirPath, e.Name()))
	}
	sort.Strings(filePaths)

	// Identify NicClusterPolicy file(s)
	var ncpFile string
	for _, p := range filePaths {
		b, rErr := os.ReadFile(p)
		if rErr != nil {
			return rErr
		}
		if containsNicClusterPolicyKind(b) {
			if ncpFile != "" {
				return fmt.Errorf("multiple NicClusterPolicy manifests found; only one is allowed")
			}
			ncpFile = p
		}
	}

	// Apply NicClusterPolicy file first if present
	if ncpFile != "" {
		content, rErr := os.ReadFile(ncpFile)
		if rErr != nil {
			return rErr
		}
		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(content, obj); err != nil {
			return fmt.Errorf("failed to decode NicClusterPolicy: %w", err)
		}
		// Ensure GVK set for server-side apply
		apiv, kind := obj.GetAPIVersion(), obj.GetKind()
		if apiv != "" && kind != "" {
			gv, err := schema.ParseGroupVersion(apiv)
			if err == nil {
				obj.SetGroupVersionKind(gv.WithKind(kind))
			}
		}
		if err := applyUnstructured(ctx, c, obj); err != nil {
			return err
		}
		if err := netophelper.WaitNicClusterPolicyReady(ctx, c, obj.GetName()); err != nil {
			return err
		}
	}

	// Apply remaining files in order
	for _, p := range filePaths {
		if p == ncpFile {
			continue
		}
		content, rErr := os.ReadFile(p)
		if rErr != nil {
			return rErr
		}
		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(content, obj); err != nil {
			return fmt.Errorf("failed to decode manifest %s: %w", p, err)
		}
		// Ensure GVK set for server-side apply
		apiv, kind := obj.GetAPIVersion(), obj.GetKind()
		if apiv != "" && kind != "" {
			gv, err := schema.ParseGroupVersion(apiv)
			if err == nil {
				obj.SetGroupVersionKind(gv.WithKind(kind))
			}
		}
		if err := applyUnstructured(ctx, c, obj); err != nil {
			return err
		}
	}

	return nil
}

func containsNicClusterPolicyKind(b []byte) bool {
	// Very small YAML sniffing: unmarshal just Kind
	type metaOnly struct {
		Kind string `yaml:"kind"`
	}
	var mo metaOnly
	if err := yaml.Unmarshal(b, &mo); err != nil {
		return false
	}
	return mo.Kind == "NicClusterPolicy"
}

func applyUnstructured(ctx context.Context, c client.Client, obj *unstructured.Unstructured) error {
	// kubectl-style server-side apply
	return c.Patch(ctx, obj, client.Apply, client.FieldOwner("l8k"), client.ForceOwnership)
}
