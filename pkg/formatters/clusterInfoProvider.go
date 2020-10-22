package formatters

import (
	"fmt"
	"strings"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const AzureProviderID = "azure:///"
const GCPProviderID = "gce://"

type ClusterInfo struct {
	ClusterName     string
	CloudProvider   string
	CloudIdentifier string //Either Resource group and Subscription Id for AKS, GCP ProjectId for GCP clusters
}

func GenerateClusterInfo(config config.Config) *ClusterInfo {
	kubeConfig, err := rest.InClusterConfig()
	fmt.Printf("Connecting to the cluster with incluster config")
	if err != nil {
		fmt.Printf("Error connecting the cluster...")
		fmt.Printf(err.Error())
		return &ClusterInfo{} //This is optional data so default should work
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	// get the first node
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error retrieving node details")
		fmt.Print(err.Error())
		return &ClusterInfo{}
	}
	firstNode := nodes.Items[0]
	specId := firstNode.Spec.ProviderID

	if strings.HasPrefix(specId, AzureProviderID) {
		return ExtractClusterInfoForAKS(firstNode)
	} else if strings.HasPrefix(specId, GCPProviderID) {
		return ExtractClusterInfoForGCP(firstNode)
	}

	return &ClusterInfo{}
}

func ExtractClusterInfoForGCP(node v1.Node) *ClusterInfo {
	specId := node.Spec.ProviderID
	formattedString := specId[len(GCPProviderID):]
	splitBySlash := strings.Split(formattedString, "/")
	// an example of the provider ID for GCP is
	// gce://dbce-c360-isl-dev-f03b/us-east4-b/gke-gke-cluster-gateway-pool-7b945b80-kqw5
	// Can be inferred that its segregated as gce://projectID/location/node-name
	if len(splitBySlash) >= 2 {
		return &ClusterInfo{
			ClusterName:     node.GetLabels()["cluster_name"],
			CloudIdentifier: splitBySlash[0],
			CloudProvider:   "GCP",
		}
	}

	return &ClusterInfo{}
}

func ExtractClusterInfoForAKS(node v1.Node) *ClusterInfo {
	specId := node.Spec.ProviderID
	formattedString := specId[len(AzureProviderID):]
	splitBySlash := strings.Split(formattedString, "/")
	// an example of the providerID for Azure is
	// azure:///subscriptions/32287844-5d4c-4a4b-b15d-8df79f438561/resourceGroups/mc_az-rg-k8s-nas-propensity-dev-01_aks-da-nas-dev_eastus2/providers/Microsoft.Compute/virtualMachineScaleSets/aks-gatewaypool-41380962-vmss/virtualMachines/0
	// Can be inferred that its segregated as azure://subscriptions/subId/resourceGroups/rgName
	// while the resource group here is structured as mc_{rg-name}_{cluster}_{location}
	resourceGroup := splitBySlash[3]
	splitResourceGroupByUnderscore := strings.Split(resourceGroup, "_")
	return &ClusterInfo{
		CloudProvider:   "Azure",
		CloudIdentifier: "ResourceGroup: " + splitResourceGroupByUnderscore[1] + " Subscription: " + splitBySlash[1],
		ClusterName:     splitResourceGroupByUnderscore[2],
	}
}
