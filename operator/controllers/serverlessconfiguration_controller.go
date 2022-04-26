/*
Copyright 2022.

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

	serverlessApi "github.com/Tomasz-Smelcerz-SAP/kyma-operator-serverless/k8s-api/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ServerlessConfigurationReconciler reconciles a ServerlessConfiguration object
type ServerlessConfigurationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kyma.kyma-project.io,resources=serverlessconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kyma.kyma-project.io,resources=serverlessconfigurations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kyma.kyma-project.io,resources=serverlessconfigurations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServerlessConfiguration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ServerlessConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling:", "object", req.String())

	configMapName := req.Name + "-config"
	configMapObjKey := client.ObjectKey{
		Name:      configMapName,
		Namespace: req.Namespace,
	}

	obj := serverlessApi.ServerlessConfiguration{}
	err := r.Client.Get(ctx, req.NamespacedName, &obj)

	if apierrors.IsNotFound(err) {
		//object is deleted
		logger.Info("Object is deleted:", "object", req.NamespacedName)

		//try to delete related ConfigMap object
		cm := corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: req.Namespace,
			},
		}

		err = r.Client.Delete(ctx, &cm)
		if apierrors.IsNotFound(err) {
			//ConfigMap does not exist. Success.
			return ctrl.Result{}, nil
		}
		if err != nil {
			logger.Error(err, "Error deleting ConfigMap", "object:", configMapObjKey)
			return ctrl.Result{}, err
		}
		logger.Info("Successfully deleted ConfigMap:", "object:", configMapObjKey)

		return ctrl.Result{}, nil
	}

	if err != nil {
		logger.Error(err, "Error reading ServerlessConfiguration", "object", req.NamespacedName)
		return ctrl.Result{}, err
	}

	githubAuthKey := obj.Spec.GithubRepository.AuthKey
	githubUrl := obj.Spec.GithubRepository.URL

	data := map[string]string{}
	data["githubAuthKey"] = githubAuthKey
	data["githubUrl"] = githubUrl

	cm := corev1.ConfigMap{}
	err = r.Client.Get(ctx, configMapObjKey, &cm)

	if apierrors.IsNotFound(err) {
		//ConfigMap does not exist yet - must be created
		logger.Info("ConfigMap does not exist:", "object", configMapObjKey)

		cm = corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: req.Namespace,
				Labels:    setLabels(nil),
			},
			Data: data,
		}

		err := r.Client.Create(ctx, &cm)
		if err != nil {
			logger.Error(err, "Error creating ConfigMap", "object:", configMapObjKey)
			return ctrl.Result{}, err
		}
		logger.Info("Successfully created ConfigMap:", "object:", configMapObjKey)

		return ctrl.Result{}, nil
	}

	if err != nil {
		logger.Error(err, "Error reading ConfigMap:", "object:", configMapObjKey)
		return ctrl.Result{}, err
	}

	//ConfigMap exists, it must be updated

	cm.Data = data
	cm.Labels = setLabels(cm.Labels)

	err = r.Client.Update(ctx, &cm)
	if err != nil {
		logger.Error(err, "Error updating ConfigMap", "object:", configMapObjKey)
		return ctrl.Result{}, err
	}
	logger.Info("Successfully updated ConfigMap", "object:", configMapObjKey)

	logger.Info("Successfully reconciled ServerlessConfiguration:", "object:", req.String())

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServerlessConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&serverlessApi.ServerlessConfiguration{}).
		Complete(r)
}

func setLabels(target map[string]string) map[string]string {
	if target == nil {
		target = map[string]string{}
	}

	//see: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	target["app.kubernetes.io/part-of"] = "kyma-project.io_kyma"
	target["app.kubernetes.io/component"] = "serverless"
	target["app.kubernetes.io/created-by"] = "serverless-controller"
	target["kyma-project.io/watched-by"] = "mothership"

	return target
}
