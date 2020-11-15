package kubecli

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/kubectl/pkg/metricsutil"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetPodMetrics(m *metricsapi.PodMetrics) v1.ResourceList {
	podMetrics := make(v1.ResourceList)
	for _, res := range metricsutil.MeasuredResources {
		podMetrics[res], _ = resource.ParseQuantity("0")
	}

	for _, c := range m.Containers {
		for _, res := range metricsutil.MeasuredResources {
			quantity := podMetrics[res]
			quantity.Add(c.Usage[res])
			podMetrics[res] = quantity
		}
	}
	return podMetrics
}

func GetAllResourceUsages(metrics *metricsutil.ResourceMetricsInfo) map[v1.ResourceName]int64 {
	result := make(map[v1.ResourceName]int64, 0)
	for _, res := range metricsutil.MeasuredResources {
		quantity := metrics.Metrics[res]
		usage := GetSingleResourceUsage(res, quantity)
		result[res] = usage
		if available, found := metrics.Available[res]; found {
			fraction := float64(quantity.MilliValue()) / float64(available.MilliValue()) * 100
			result[res] = int64(fraction)
		}
	}
	return result
}

func GetSingleResourceUsage(resourceType v1.ResourceName, quantity resource.Quantity) int64 {
	switch resourceType {
	case v1.ResourceCPU:
		return quantity.MilliValue()
	case v1.ResourceMemory:
		return quantity.Value() / (1024 * 1024)
	default:
		return quantity.Value()
	}
}

func (cli *KubeCLI) GetPodRawMetrics(namespace, name string, allNamespaces bool, selector labels.Selector) (*metricsapi.PodMetricsList, error) {
	if selector == nil {
		selector = labels.Everything()
	}

	var err error
	config, err := cli.factory.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ns := metav1.NamespaceAll
	if !allNamespaces {
		ns = namespace
	}
	versionedMetrics := &metricsv1beta1api.PodMetricsList{}
	if name != "" {
		m, err := metricsClient.MetricsV1beta1().PodMetricses(ns).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		versionedMetrics.Items = []metricsv1beta1api.PodMetrics{*m}
	} else {
		versionedMetrics, err = metricsClient.MetricsV1beta1().PodMetricses(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return nil, err
		}
	}
	metrics := &metricsapi.PodMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_PodMetricsList_To_metrics_PodMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (cli *KubeCLI) GetPodMetrics(namespace, name string, allNamespaces bool, selector labels.Selector) ([]map[v1.ResourceName]int64, error) {
	metrics, err := cli.GetPodRawMetrics(namespace, name, allNamespaces, selector)
	if err != nil {
		return nil, err
	}

	result := make([]map[v1.ResourceName]int64, 0)
	for _, metric := range metrics.Items {
		podMetrics := GetPodMetrics(&metric)
		metricsInfo := &metricsutil.ResourceMetricsInfo{
			Name:      metric.Name,
			Metrics:   podMetrics,
			Available: v1.ResourceList{},
		}
		result = append(result, GetAllResourceUsages(metricsInfo))
	}
	return result, nil
}
