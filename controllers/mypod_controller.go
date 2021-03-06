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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	mypodv1 "github.com/harryyann/mypod-operator/api/v1"
)

// MyPodReconciler reconciles a MyPod object
type MyPodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=harryyann.github.io,resources=mypods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=harryyann.github.io,resources=mypods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=harryyann.github.io,resources=mypods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyPod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MyPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 按照name和namespace获取MyPod对象
	myPod := mypodv1.MyPod{}
	err := r.Get(ctx, req.NamespacedName, &myPod)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// 如果mypod对象没有的话，就把对应的pod删了
			pod := v1.Pod{}
			pod.Name = myPod.Name
			pod.Namespace = myPod.Namespace
			err = r.Delete(ctx, &pod)
		}
		return ctrl.Result{}, err
	}

	pod := v1.Pod{}
	err = r.Get(ctx, req.NamespacedName, &pod)
	if client.IgnoreNotFound(err) != nil {
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// 没有这个pod就创建一个
		pod.Name = myPod.Name
		pod.Namespace = myPod.Namespace
		pod.Annotations = myPod.Spec.PodAnnotations
		pod.Labels = myPod.Spec.PodLabels
		pod.Labels["creator"] = "mypod-controller"
		pod.Spec = myPod.Spec.PodSpec
		err = r.Create(ctx, &pod, &client.CreateOptions{})
		if err != nil {
			logger.Error(err, "create pod error")
			return ctrl.Result{}, err
		}
		logger.Info("create pod " + pod.Name + " success")
		goto UpdateStatus
	}

	// 更新我们就先把pod删了，再重建
	err = r.Delete(ctx, &pod)
	if err != nil {
		return ctrl.Result{}, err
	}
	pod = v1.Pod{}
	pod.Name = myPod.Name
	pod.Namespace = myPod.Namespace
	pod.Annotations = myPod.Spec.PodAnnotations
	pod.Labels = myPod.Spec.PodLabels
	pod.Labels["creator"] = "mypod-controller"
	pod.Spec = myPod.Spec.PodSpec
	err = r.Create(ctx, &pod, &client.CreateOptions{})
	if err != nil {
		logger.Error(err, "create pod error")
		return ctrl.Result{}, err
	}
	logger.Info("update pod " + pod.Name + " success")

UpdateStatus:
	// 在创建这个pod之后持续的监控这个pod，拿到想要的状态信息，填到status中，这里简化写，就让它睡5秒得了
	time.Sleep(time.Second * 5)
	err = r.Get(ctx, req.NamespacedName, &pod)
	myPod.Status.PodPhase = string(pod.Status.Phase)
	myPod.Status.NodeIp = pod.Status.HostIP
	myPod.Status.PodIp = pod.Status.PodIP
	err = r.Status().Update(ctx, &myPod)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mypodv1.MyPod{}).
		Complete(r)
}
