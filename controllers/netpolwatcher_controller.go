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
	v1 "k8s.io/api/networking/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NetPolWatcherReconciler reconciles a NetPolWatcher object
type NetPolWatcherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=netpol.qaware.com,resources=netpolwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=netpol.qaware.com,resources=netpolwatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies/status,verbs=get

func (r *NetPolWatcherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("netpolwatcher", req.NamespacedName)

	r.Log.Info("Hello Netpol", "Name", req.Name, "Namespace", req.Namespace)

	netpol := v1.NetworkPolicy{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, &netpol)
	if err != nil {
		r.Log.Error(err, "Error getting network policy", "Name", req.Name)
	}
	r.Log.Info(netpol.String())

	return ctrl.Result{}, nil
}

func (r *NetPolWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.NetworkPolicy{}).
		Complete(r)
}
