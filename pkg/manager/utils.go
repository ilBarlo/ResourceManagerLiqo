package manager

import (
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func isPeeringEstablished(peeringCondition discoveryv1alpha1.PeeringConditionStatusType) bool {
	return peeringCondition == discoveryv1alpha1.PeeringConditionStatusEstablished
}

func aggregatePodsMetrics(podsMetrics []metricsv1beta1.PodMetrics) (cpuUsage, memUsage resource.Quantity) {
	cpuUsage = *resource.NewMilliQuantity(0, "BinarySI")
	memUsage = *resource.NewMilliQuantity(0, "BinarySI")

	for i := range podsMetrics {
		for _, container := range podsMetrics[i].Containers {
			if container.Usage.Cpu() != nil {
				cpuUsage.Add(*container.Usage.Cpu())
			}
			if container.Usage.Memory() != nil {
				memUsage.Add(*container.Usage.Memory())
			}
		}
	}

	return cpuUsage, memUsage
}

func peeringConditionsToMap(
	peeringConditions []discoveryv1alpha1.PeeringCondition,
) map[discoveryv1alpha1.PeeringConditionType]discoveryv1alpha1.PeeringCondition {
	cpm := make(map[discoveryv1alpha1.PeeringConditionType]discoveryv1alpha1.PeeringCondition)
	for i := range peeringConditions {
		cpm[peeringConditions[i].Type] = peeringConditions[i]
	}
	return cpm
}

func podMetricListToMap(pods []metricsv1beta1.PodMetrics) map[string]*metricsv1beta1.PodMetrics {
	podsMap := make(map[string]*metricsv1beta1.PodMetrics)
	for i := range pods {
		podsMap[pods[i].Name] = &pods[i]
	}
	return podsMap
}

func statusOrDefault(
	sourceMap map[discoveryv1alpha1.PeeringConditionType]discoveryv1alpha1.PeeringCondition,
	key discoveryv1alpha1.PeeringConditionType,
) discoveryv1alpha1.PeeringConditionStatusType {
	value, found := sourceMap[key]
	if !found {
		return "None"
	}
	return value.Status
}
