/*
Copyright 2021 Sugun.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	//"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	shipyardv1beta1 "shipyard.dev/secretWhisperer/api/v1beta1"
)

// WhispererReconciler reconciles a Whisperer object
type WhispererReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=shipyard.shipyard.dev,resources=whisperers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=shipyard.shipyard.dev,resources=whisperers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=shipyard.shipyard.dev,resources=whisperers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets;events,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Whisperer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *WhispererReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Whisperer")

	// your logic here
	ws := &shipyardv1beta1.Whisperer{}
	if err := r.Client.Get(ctx, req.NamespacedName, ws); err != nil {
		log.Error(err, "failed to get whisperer resource")
		// Ignore NotFound errors as they will be retried automatically if the
		// resource is created in future.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cert_secret := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      ws.Spec.SecretName,
		Namespace: req.Namespace,
	}, cert_secret)
	if err != nil {
		log.Error(err, "failed to get secret resource in the namespace")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	istio_secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ws.Spec.SecretName, //name should match that of the original secert
			Namespace: "istio-system",
		},
		Data: cert_secret.Data,
		Type: cert_secret.Type,
	}

	_, err = controllerutil.CreateOrPatch(ctx, r.Client, istio_secret, func() error {
		if istio_secret.Labels == nil {
			istio_secret.Labels = make(map[string]string)
		}
		istio_secret.Labels["SourceNamespace"] = req.Namespace

		return controllerutil.SetControllerReference(ws, istio_secret, r.Scheme)
	})
	if err != nil {
		log.Error(err, "unable to create or update secret")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WhispererReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipyardv1beta1.Whisperer{}).
		Complete(r)
}
