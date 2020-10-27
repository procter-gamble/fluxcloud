package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractClusterInfoForGCP(t *testing.T) {

	labels := make(map[string]string)
	labels["cluster_name"] = "gke_cluster"
	fakeNode := v1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind: "Node",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.NodeSpec{
			ProviderID: "gce://dbce-c360-isl-dev-f03b/us-east4-b/gke-gke-cluster-gateway-pool-7b945b80-kqw5",
		},
	}

	clusterInfo := ExtractClusterInfoForGCP(fakeNode)

	assert.Equal(t, "gke_cluster", clusterInfo.ClusterName)
	assert.Equal(t, "dbce-c360-isl-dev-f03b", clusterInfo.CloudIdentifier)
	assert.Equal(t, "GCP", clusterInfo.CloudProvider)
}

func TestExtractClusterInfoForAKS(t *testing.T) {
	labels := make(map[string]string)
	labels["cluster_name"] = "gke_cluster"
	fakeNode := v1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind: "Node",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.NodeSpec{
			ProviderID: "azure:///subscriptions/32287844-5d4c-4a4b-b15d-8df79f438561/resourceGroups/mc_az-rg-k8s-nas-propensity-dev-01_aks-da-nas-dev_eastus2/providers/Microsoft.Compute/virtualMachineScaleSets/aks-gatewaypool-41380962-vmss/virtualMachines/0",
		},
	}

	clusterInfo := ExtractClusterInfoForAKS(fakeNode)

	assert.Equal(t, "aks-da-nas-dev", clusterInfo.ClusterName)
	assert.Equal(t, "ResourceGroup: az-rg-k8s-nas-propensity-dev-01 Subscription: 32287844-5d4c-4a4b-b15d-8df79f438561", clusterInfo.CloudIdentifier)
	assert.Equal(t, "Azure", clusterInfo.CloudProvider)
}
