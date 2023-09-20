package main

import (
	"log"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// gofmt -w -s *.go; pulumi up -y
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		isMinikube := config.GetBool(ctx, "isMinikube")
		appName := "nginx"
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}
		deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String("nginx"),
							}},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		feType := "LoadBalancer"
		if isMinikube {
			feType = "ClusterIP"
		}
		log.Printf("deployment=[%T][%v]", deployment.ToDeploymentOutput(), deployment.ToDeploymentOutput())
		log.Printf("deployment.Spec=[%T][%v]", deployment.ToDeploymentOutput().Spec(), deployment.ToDeploymentOutput().Spec())
		template := deployment.Spec.ApplyT(func(v appsv1.DeploymentSpec) *corev1.PodTemplateSpec {
			return &v.Template
		}).(corev1.PodTemplateSpecPtrOutput)
		log.Printf("template=[%T][%v]", template, template)

		meta := template.ApplyT(func(v *corev1.PodTemplateSpec) *metav1.ObjectMeta {
			return v.Metadata
		}).(metav1.ObjectMetaPtrOutput)
		log.Printf("meta=[%T][%v]", meta, meta)

		frontend, err := corev1.NewService(ctx, appName, &corev1.ServiceArgs{
			Metadata: meta,
			Spec: &corev1.ServiceSpecArgs{
				Type: pulumi.String(feType),
				Ports: &corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
					},
				},
				Selector: appLabels,
			},
		})

		var ip pulumi.StringOutput

		if isMinikube {
			ip = frontend.Spec.ApplyT(func(val corev1.ServiceSpec) string {
				log.Printf("ClusterIP=[%v]", val.ClusterIP)
				if val.ClusterIP != nil {
					return *val.ClusterIP
				}
				return ""
			}).(pulumi.StringOutput)
		} else {
			ip = frontend.Status.ApplyT(func(val *corev1.ServiceStatus) string {
				log.Printf("LoadBalancer.Ingress=[%v]", val.LoadBalancer.Ingress)
				if len(val.LoadBalancer.Ingress) <= 0 {
					return ""
				}
				if val.LoadBalancer.Ingress[0].Ip != nil {
					return *val.LoadBalancer.Ingress[0].Ip
				}
				return *val.LoadBalancer.Ingress[0].Hostname
			}).(pulumi.StringOutput)
		}

		ctx.Export("ip", ip)
		ctx.Export("name", deployment.Metadata.Name())

		return nil
	})
}
