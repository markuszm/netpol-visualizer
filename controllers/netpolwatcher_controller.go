/*


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

package controllers

import (
	"context"
	"fmt"
	"github.com/markuszm/netpol-visualizer/database"
	"github.com/markuszm/netpol-visualizer/model"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NetPolWatcherReconciler reconciles a NetPolWatcher object
type NetPolWatcherReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Database database.Client
}

var Everything = model.Pod{
	Name:      "EVERYTHING",
	Namespace: "EVERYWHERE",
}

// +kubebuilder:rbac:groups=netpol.qaware.com,resources=netpolwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=netpol.qaware.com,resources=netpolwatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies/status,verbs=get
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

func (r *NetPolWatcherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("netpolwatcher", req.NamespacedName)

	netpol := networkingv1.NetworkPolicy{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, &netpol)
	if err != nil {
		r.Log.Error(err, "Error getting network policy", "Name", req.NamespacedName.String())
	}

	selector := netpol.Spec.PodSelector

	selectedPods := []model.Pod{}
	ingressEdges := []model.UnrestrictedEdge{}
	egressEdges := []model.UnrestrictedEdge{}

	if selectorIsEmpty(selector) {
		// select all pods here from namespace
		podList := v1.PodList{}
		err := r.Client.List(ctx, &podList, client.InNamespace(req.Namespace))
		if err != nil {
			r.Log.Error(err, "Error getting pods", "Name", req.NamespacedName.String())
		}
		for _, pod := range podList.Items {
			selectedPods = append(selectedPods, model.Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			})
		}
	} else {
		podList := v1.PodList{}
		listOptions, err := labelSelectorToListOptions(selector)
		if err != nil {
			r.Log.Error(err, "Error creating pod label selector", "Name", req.NamespacedName.String())
		}
		listOptions = append(listOptions, client.InNamespace(req.Namespace))
		err = r.Client.List(ctx, &podList, listOptions...)
		if err != nil {
			r.Log.Error(err, "Error getting pods", "Name", req.NamespacedName.String())
		}
		for _, pod := range podList.Items {
			selectedPods = append(selectedPods, model.Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			})
		}
	}

	for _, pod := range selectedPods {
		r.Log.Info(fmt.Sprintf("Selected Pod: %s/%s", pod.Namespace, pod.Name), "Name", req.NamespacedName.String())
	}

	// Ingress Case
	if len(netpol.Spec.PolicyTypes) == 0 || containsPolicyType(netpol.Spec.PolicyTypes, networkingv1.PolicyTypeIngress) || len(netpol.Spec.Ingress) > 0 {
		// loop through ingress rules - if none exists then access restricted
		for _, ingressRule := range netpol.Spec.Ingress {
			var ports []int
			// get all ports if none specified then add 0 as default
			for _, port := range ingressRule.Ports {
				ports = append(ports, port.Port.IntValue())
				// if named ports panic
			}
			if len(ports) == 0 {
				ports = append(ports, 0)
			}
			// if no ingress rule specified then for all pods everything is allowed as Ingress
			if len(ingressRule.From) == 0 {
				for _, pod := range selectedPods {
					for _, port := range ports {
						ingressEdges = append(ingressEdges, model.UnrestrictedEdge{
							From: Everything,
							To:   pod,
							Port: port,
						})
					}
				}
			}
			// loop through ingress rules specifing from
			for _, peer := range ingressRule.From {
				// if pod selector specified then select pods later
				var podSelectorListOptions []client.ListOption
				if peer.PodSelector != nil {
					podSelectorListOptions, err = labelSelectorToListOptions(*peer.PodSelector)
					if err != nil {
						r.Log.Error(err, "Error creating pod label selector", "Name", req.NamespacedName.String())
					}
				}

				// resolve namespace selector, if none specified then add network policy namespace to list
				var namespaces []string
				if peer.NamespaceSelector != nil {
					listOptions, err := labelSelectorToListOptions(*peer.NamespaceSelector)
					if err != nil {
						r.Log.Error(err, "Error creating namespace label selector", "Name", req.NamespacedName.String())
					}
					namespaceList := v1.NamespaceList{}
					err = r.Client.List(ctx, &namespaceList, listOptions...)
					if err != nil {
						r.Log.Error(err, "Error getting namespaces for namespace selector", "Name", req.NamespacedName.String())
					}
					for _, item := range namespaceList.Items {
						namespaces = append(namespaces, item.Name)
					}
				} else {
					namespaces = append(namespaces, req.Namespace)
				}

				// for every namespace select all pods with pod selector or existing in the namespace if no selector specified
				for _, namespace := range namespaces {
					podList := v1.PodList{}
					podSelectorListOptions = append(podSelectorListOptions, client.InNamespace(namespace))
					err = r.Client.List(ctx, &podList, podSelectorListOptions...)
					if err != nil {
						r.Log.Error(err, "Error getting pods", "Name", req.NamespacedName.String())
					}

					// add ingress edges for each pod from each by network policy selected pod
					for _, ingressPod := range podList.Items {
						for _, pod := range selectedPods {
							for _, port := range ports {
								ingressEdges = append(ingressEdges, model.UnrestrictedEdge{
									From: model.Pod{
										Name:      ingressPod.Name,
										Namespace: ingressPod.Namespace,
									},
									To:   pod,
									Port: port,
								})
							}
						}
					}
				}
			}
		}
	}

	// Egress Case
	if containsPolicyType(netpol.Spec.PolicyTypes, networkingv1.PolicyTypeEgress) || len(netpol.Spec.Egress) > 0 {
		// loop through ingress rules - if none exists then access restricted
		for _, egressRule := range netpol.Spec.Egress {
			var ports []int
			// get all ports if none specified then add 0 as default
			for _, port := range egressRule.Ports {
				ports = append(ports, port.Port.IntValue())
				// if named ports panic
			}
			if len(ports) == 0 {
				ports = append(ports, 0)
			}
			// if no ingress rule specified then for all pods everything is allowed as Egress
			if len(egressRule.To) == 0 {
				for _, pod := range selectedPods {
					for _, port := range ports {
						egressEdges = append(egressEdges, model.UnrestrictedEdge{
							From: Everything,
							To:   pod,
							Port: port,
						})
					}
				}
			}
			// loop through ingress rules specifing from
			for _, peer := range egressRule.To {
				// if pod selector specified then select pods later
				var podSelectorListOptions []client.ListOption
				if peer.PodSelector != nil {
					podSelectorListOptions, err = labelSelectorToListOptions(*peer.PodSelector)
					if err != nil {
						r.Log.Error(err, "Error creating pod label selector", "Name", req.NamespacedName.String())
					}
				}

				// resolve namespace selector, if none specified then add network policy namespace to list
				var namespaces []string
				if peer.NamespaceSelector != nil {
					listOptions, err := labelSelectorToListOptions(*peer.NamespaceSelector)
					if err != nil {
						r.Log.Error(err, "Error creating namespace label selector", "Name", req.NamespacedName.String())
					}
					namespaceList := v1.NamespaceList{}
					err = r.Client.List(ctx, &namespaceList, listOptions...)
					if err != nil {
						r.Log.Error(err, "Error getting namespaces for namespace selector", "Name", req.NamespacedName.String())
					}
					for _, item := range namespaceList.Items {
						namespaces = append(namespaces, item.Name)
					}
				} else {
					namespaces = append(namespaces, req.Namespace)
				}

				// for every namespace select all pods with pod selector or existing in the namespace if no selector specified
				for _, namespace := range namespaces {
					podList := v1.PodList{}
					podSelectorListOptions = append(podSelectorListOptions, client.InNamespace(namespace))
					err = r.Client.List(ctx, &podList, podSelectorListOptions...)
					if err != nil {
						r.Log.Error(err, "Error getting pods", "Name", req.NamespacedName.String())
					}

					// add ingress edges for each pod from each by network policy selected pod
					for _, ingressPod := range podList.Items {
						for _, pod := range selectedPods {
							for _, port := range ports {
								egressEdges = append(egressEdges, model.UnrestrictedEdge{
									From: pod,
									To: model.Pod{
										Name:      ingressPod.Name,
										Namespace: ingressPod.Namespace,
									},
									Port: port,
								})
							}
						}
					}
				}
			}
		}
	}

	for _, edge := range ingressEdges {
		r.Log.Info(fmt.Sprintf("Create Ingress Edge: From: %s/%s To: %s/%s Port: %v", edge.From.Namespace, edge.From.Name, edge.To.Namespace, edge.To.Name, edge.Port), "Name", req.NamespacedName.String())
	}

	for _, edge := range egressEdges {
		r.Log.Info(fmt.Sprintf("Create Egress Edge: From: %s/%s To: %s/%s Port: %v", edge.From.Namespace, edge.From.Name, edge.To.Namespace, edge.To.Name, edge.Port), "Name", req.NamespacedName.String())
	}

	r.Log.Info("Trying to insert ingress edges")
	err = r.Database.Insert(ingressEdges)
	if err != nil {
		r.Log.Error(err, "Error inserting ingress edges", "Name", req.NamespacedName)
	}
	r.Log.Info("Trying to insert egress edges")
	err = r.Database.Insert(egressEdges)
	if err != nil {
		r.Log.Error(err, "Error inserting egress edges", "Name", req.NamespacedName)
	}

	return ctrl.Result{}, nil
}

func (r *NetPolWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.NetworkPolicy{}).
		Complete(r)
}

func selectorIsEmpty(selector metav1.LabelSelector) bool {
	matchLabelsEmpty := selector.MatchLabels == nil || len(selector.MatchLabels) == 0
	matchExpressionsEmpty := selector.MatchExpressions == nil || len(selector.MatchExpressions) == 0
	return matchLabelsEmpty && matchExpressionsEmpty
}

func labelSelectorOperatorToSelectorOperator(op metav1.LabelSelectorOperator) selection.Operator {
	switch op {
	case metav1.LabelSelectorOpIn:
		return selection.In
	case metav1.LabelSelectorOpNotIn:
		return selection.NotIn
	case metav1.LabelSelectorOpExists:
		return selection.Exists
	case metav1.LabelSelectorOpDoesNotExist:
		return selection.DoesNotExist
	}
	return selection.Equals
}

func labelSelectorToListOptions(selector metav1.LabelSelector) ([]client.ListOption, error) {
	var listOptions []client.ListOption
	labelsSelector := labels.NewSelector()
	for _, labelSelectorRequirement := range selector.MatchExpressions {
		requirement, err := labels.NewRequirement(labelSelectorRequirement.Key, labelSelectorOperatorToSelectorOperator(labelSelectorRequirement.Operator), labelSelectorRequirement.Values)
		if err == nil {
			labelsSelector = labelsSelector.Add(*requirement)
		}
	}
	listOptions = append(listOptions, client.MatchingLabelsSelector{Selector: labelsSelector})
	if selector.MatchLabels != nil && len(selector.MatchLabels) > 0 {
		listOptions = append(listOptions, client.MatchingLabels(selector.MatchLabels))
	}
	return listOptions, nil
}

func containsPolicyType(policyTypes []networkingv1.PolicyType, selector networkingv1.PolicyType) bool {
	for _, policyType := range policyTypes {
		if policyType == selector {
			return true
		}
	}
	return false
}
