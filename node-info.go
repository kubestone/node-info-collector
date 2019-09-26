/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeInfo struct {
	Capacity    corev1.ResourceList
	Allocatable corev1.ResourceList
	Usage       corev1.ResourceList
	SystemInfo  corev1.NodeSystemInfo
}

func getNodeInfos(clientset *kubernetes.Clientset,
	metricsClientset *metricsv.Clientset) (map[string]*NodeInfo, error) {
	nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeInfos := make(map[string]*NodeInfo, len(nodeList.Items))
	for _, node := range nodeList.Items {
		nodeInfos[node.ObjectMeta.Name] = &NodeInfo{
			Capacity:    node.Status.Capacity,
			Allocatable: node.Status.Allocatable,
			SystemInfo:  node.Status.NodeInfo,
		}
	}

	nodeMetricsList, err := metricsClientset.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
	for _, nodeMetrics := range nodeMetricsList.Items {
		nodeInfos[nodeMetrics.ObjectMeta.Name].Usage = nodeMetrics.Usage
	}
	return nodeInfos, nil
}
