package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		replicas := 1
		if replicas = cfg.GetInt("replicas"); replicas < 0 {
			replicas = 1
		}
		appName := "nginx"
		appImage := "nginx:1.8"
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}
		deployment, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(replicas),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap(appLabels),
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap(appLabels),
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String(appImage),
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("name", deployment.Metadata.ApplyT(func(metadata metav1.ObjectMeta) (*string, error) {
			return metadata.Name, nil
		}).(pulumi.StringPtrOutput))
		return nil
	})
}
