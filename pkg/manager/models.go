package manager

import (
	"encoding/json"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ArchitectureResources represents the map used to see the available resources
type ArchitectureResources map[OSArchitecture]ResourcesAvailable

// OSArchitecture represents the key used for ArchitectureResources
type OSArchitecture struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

// NodeInfo represents the information of a certain node of a cluster which is returned to the client
type NodeInfo struct {
	Name         string `json:"name"`
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

// ResourceMetrics represents the metrics of a cluster which is returned to the client.
type ResourcesAvailable struct {
	AvailableCpus   float64 `json:"availableCpus"`
	AvailableMemory float64 `json:"availableMemory"`
	TotalMemory     float64 `json:"totalMemory"`
	TotalCpus       float64 `json:"totalCpus"`
}

// ClusterDto represents the data of a cluster which is returned to the client.
type ClusterDto struct {
	clusterID         string
	Name              string                                       `json:"name"`
	NodeInfo          NodeInfo                                     `json:"node"`
	Networking        discoveryv1alpha1.PeeringConditionStatusType `json:"networking"`
	Authentication    discoveryv1alpha1.PeeringConditionStatusType `json:"authentication"`
	OutgoingPeering   discoveryv1alpha1.PeeringConditionStatusType `json:"outgoingPeering"`
	OutgoingResources *ResourceMetrics                             `json:"outgoingResources"`
	IncomingPeering   discoveryv1alpha1.PeeringConditionStatusType `json:"incomingPeering"`
	IncomingResources *ResourceMetrics                             `json:"incomingResources"`
	Age               string                                       `json:"age,omitempty"`
}

// ResourceMetrics represents the metrics of a cluster which is returned to the client.
type ResourceMetrics struct {
	UsedCpus    float64 `json:"usedCpus"`
	UsedMemory  float64 `json:"usedMemory"`
	TotalMemory float64 `json:"totalMemory"`
	TotalCpus   float64 `json:"totalCpus"`
}

// ErrorResponse is returned to the client in case of error.
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int16  `json:"status"`
}

func fromForeignCluster(fc *discoveryv1alpha1.ForeignCluster) *ClusterDto {
	pc := peeringConditionsToMap(fc.Status.PeeringConditions)

	clusterDto := &ClusterDto{
		Name:            fc.Name,
		clusterID:       fc.Spec.ClusterIdentity.ClusterID,
		OutgoingPeering: statusOrDefault(pc, discoveryv1alpha1.OutgoingPeeringCondition),
		IncomingPeering: statusOrDefault(pc, discoveryv1alpha1.IncomingPeeringCondition),
		Networking:      statusOrDefault(pc, discoveryv1alpha1.NetworkStatusCondition),
		Authentication:  statusOrDefault(pc, discoveryv1alpha1.AuthenticationStatusCondition),
		NodeInfo:        NodeInfo{"unknown", "unknwon", "unknown"},
	}

	auth, found := pc[discoveryv1alpha1.AuthenticationStatusCondition]
	if found {
		clusterDto.Age = auth.LastTransitionTime.Time.String()
	}

	return clusterDto
}

func newResourceMetrics(cpuUsage, memUsage resource.Quantity, totalResources corev1.ResourceList) *ResourceMetrics {
	return &ResourceMetrics{
		UsedCpus:    cpuUsage.AsApproximateFloat64(),
		TotalCpus:   totalResources.Cpu().AsApproximateFloat64(),
		UsedMemory:  memUsage.AsApproximateFloat64(),
		TotalMemory: totalResources.Memory().AsApproximateFloat64(),
	}
}

func (ar ArchitectureResources) MarshalJSON() ([]byte, error) {
	// Creiamo una mappa temporanea con i dati serializzabili in JSON

	tempMap := make(map[string]ResourcesAvailable)
	for key, value := range ar {
		k := key.Architecture + " & " + key.OS
		tempMap[k] = value
	}
	// Serializziamo la mappa temporanea in JSON
	return json.Marshal(tempMap)
}
