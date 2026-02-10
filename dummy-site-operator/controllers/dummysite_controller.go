package controllers

import (
	"context"
	"io"
	"net/http"

	webv1alpha1 "github.com/stacvirus/dummy-site-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DummySiteReconciler struct {
	client.Client
}

func (r *DummySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var site webv1alpha1.DummySite
	if err := r.Get(ctx, req.NamespacedName, &site); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 1. Fetch website
	resp, err := http.Get(site.Spec.WebsiteURL)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 2. ConfigMap
	cm := &corev1.ConfigMap{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      site.Name + "-html",
			Namespace: site.Namespace,
		},
		Data: map[string]string{
			"index.html": string(body),
		},
	}
	_ = ctrl.SetControllerReference(&site, cm, r.Scheme())
	_ = r.Client.Create(ctx, cm)

	// 3. Deployment
	dep := nginxDeployment(site.Name, site.Namespace)
	_ = ctrl.SetControllerReference(&site, dep, r.Scheme())
	_ = r.Client.Create(ctx, dep)

	// 4. Service
	svc := nginxService(site.Name, site.Namespace)
	_ = ctrl.SetControllerReference(&site, svc, r.Scheme())
	_ = r.Client.Create(ctx, svc)

	return ctrl.Result{}, nil
}

func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1alpha1.DummySite{}).
		Complete(r)
}

func nginxDeployment(name, ns string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "html",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name + "-html",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func nginxService(name, ns string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt32(80),
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
