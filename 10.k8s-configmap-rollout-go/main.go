package main

import (
	"io/ioutil"

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

		replicas := 1
		configFileName := "default.conf"
		appName := "nginx"
		appImage := "nginx:1.13.6-alpine"
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}
		configName := "nginxconfig"
		frontendName := "frontend"

		// nginx Configuration data to proxy traffic to `pulumi.github.io`. Read from `default.conf` file.
		dataFile, err := ioutil.ReadFile(configFileName)
		if err != nil {
			return err
		}

		nginxConfig, err := corev1.NewConfigMap(ctx, appName, &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   pulumi.String(configName),
			},
			Data: pulumi.StringMap{configFileName: pulumi.String(string(dataFile))},
		})
		nginxConfigName := nginxConfig.Metadata.Name()

		// Deploy 1 nginx replica, mounting the configuration data into the nginx container.
		deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   pulumi.String(appName),
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(replicas),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String(appImage),
								VolumeMounts: &corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String("nginx-configs"),
										MountPath: pulumi.String("/etc/nginx/conf.d"),
									},
								},
							},
						},
						Volumes: &corev1.VolumeArray{
							&corev1.VolumeArgs{
								Name: pulumi.String("nginx-configs"),
								ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
									Name: nginxConfigName,
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Expose proxy to the public internet
		var frontendServiceType string
		if isMinikube {
			frontendServiceType = "ClusterIP"
		} else {
			frontendServiceType = "LoadBalancer"
		}

		frontend, err := corev1.NewService(ctx, appName, &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   pulumi.String(frontendName),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type: pulumi.String(frontendServiceType),
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

		// Export the public IP
		if isMinikube {
			ctx.Export("frontendIP", frontend.Spec.ApplyT(func(spec corev1.ServiceSpec) *string {
				return spec.ClusterIP
			}))
		} else {
			ctx.Export("frontendIP", frontend.Status.ApplyT(func(status *corev1.ServiceStatus) *string {
				if len(status.LoadBalancer.Ingress) <= 0 {
					return nil
				}
				ingress := status.LoadBalancer.Ingress[0]
				if ingress.Hostname != nil {
					return ingress.Hostname
				}
				return ingress.Ip
			}))
		}

		ctx.Export("name", deployment.Metadata.Name())
		ctx.Export("svcName", frontend.Metadata.Name())
		ctx.Export("nginxConfigName", nginxConfigName)

		return nil
	})
}
