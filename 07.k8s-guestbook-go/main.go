package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// gofmt -w -s *.go; pulumi up -y
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Initialize config
		conf := config.New(ctx, "")

		// Create only services of type `ClusterIP`
		// for clusters that don't support `LoadBalancer` services
		useLoadBalancer := conf.GetBool("useLoadBalancer")

		// Redis leader
		redisLeaderLabels := pulumi.StringMap{
			"app": pulumi.String("redis-leader"),
		}
		// Redis leader Deployment
		_, err := appsv1.NewDeployment(ctx, "redis-leader", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: redisLeaderLabels,
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: redisLeaderLabels,
				},
				Replicas: pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: redisLeaderLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("redis-leader"),
								Image: pulumi.String("redis"),
								Resources: &corev1.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(6379),
									},
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
		// Redis leader Service
		_, err = corev1.NewService(ctx, "redis-leader", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("redis-leader"),
				Labels: redisLeaderLabels,
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Port:       pulumi.Int(6379),
						TargetPort: pulumi.Int(6379),
					},
				},
				Selector: redisLeaderLabels,
			},
		})
		if err != nil {
			return err
		}

		// Redis replica
		redisReplicaLabels := pulumi.StringMap{
			"app": pulumi.String("redis-replica"),
		}
		// Redis replica Deployment
		_, err = appsv1.NewDeployment(ctx, "redis-replica", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: redisReplicaLabels,
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: redisReplicaLabels,
				},
				Replicas: pulumi.Int(2),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: redisReplicaLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("redis-replica"),
								Image: pulumi.String("pulumi/guestbook-redis-replica"),
								Resources: &corev1.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								Env: corev1.EnvVarArray{
									corev1.EnvVarArgs{
										Name:  pulumi.String("GET_HOSTS_FROM"),
										Value: pulumi.String("dns"),
									},
								},
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(6379),
									},
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
		// Redis replica Service
		_, err = corev1.NewService(ctx, "redis-replica", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("redis-replica"),
				Labels: redisReplicaLabels,
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Port: pulumi.Int(6379),
					},
				},
				Selector: redisReplicaLabels,
			},
		})
		if err != nil {
			return err
		}

		// frontend
		frontendReplicas := 1
		frontendLabels := pulumi.StringMap{
			"app": pulumi.String("frontend"),
		}
		// Frontend Deployment
		_, err = appsv1.NewDeployment(ctx, "frontend", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: frontendLabels,
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: frontendLabels,
				},
				Replicas: pulumi.Int(frontendReplicas),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: frontendLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("php-redis"),
								Image: pulumi.String("pulumi/guestbook-php-redis"),
								Resources: &corev1.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								Env: corev1.EnvVarArray{
									corev1.EnvVarArgs{
										Name:  pulumi.String("GET_HOSTS_FROM"),
										Value: pulumi.String("dns"),
									},
								},
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
									},
								},
							}},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		// Frontend Service
		var frontendServiceType string
		if useLoadBalancer {
			frontendServiceType = "LoadBalancer"
		} else {
			frontendServiceType = "ClusterIP"
		}
		frontendService, err := corev1.NewService(ctx, "frontend", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: frontendLabels,
				Name:   pulumi.String("frontend"),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type: pulumi.String(frontendServiceType),
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Port: pulumi.Int(80),
					},
				},
				Selector: frontendLabels,
			},
		})
		if err != nil {
			return err
		}

		if useLoadBalancer {
			ctx.Export("frontendIP", frontendService.Status.ApplyT(func(status *corev1.ServiceStatus) string {
				if len(status.LoadBalancer.Ingress) > 0 {
					ingress := status.LoadBalancer.Ingress[0]
					if ingress.Hostname != nil {
						return *ingress.Hostname
					}
					if ingress.Ip != nil {
						return *ingress.Ip
					}
				}
				return ""
			}))
		} else {
			ctx.Export("frontendIP", frontendService.Spec.ApplyT(func(spec corev1.ServiceSpec) string {
				if spec.ClusterIP != nil {
					return *spec.ClusterIP
				}
				return ""
			}))
		}

		return nil
	})
}
