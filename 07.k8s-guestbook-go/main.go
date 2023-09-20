package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
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

		// Redis leader Deployment + Service
		_, err := NewServiceDeployment(ctx, "redis-leader", &ServiceDeploymentArgs{
			Image: pulumi.String("redis"),
			Ports: pulumi.IntArray{pulumi.Int(6379)},
		})
		if err != nil {
			return err
		}

		// Redis replica Deployment + Service
		_, err = NewServiceDeployment(ctx, "redis-replica", &ServiceDeploymentArgs{
			Image: pulumi.String("pulumi/guestbook-redis-replica"),
			Ports: pulumi.IntArray{pulumi.Int(6379)},
		})
		if err != nil {
			return err
		}

		// Frontend Deployment + Service
		frontendReplicas := 1
		frontend, err := NewServiceDeployment(ctx, "frontend", &ServiceDeploymentArgs{
			AllocateIPAddress: true,
			Image:             pulumi.String("pulumi/guestbook-php-redis"),
			UseLoadBalancer:   pulumi.Bool(useLoadBalancer),
			Ports:             pulumi.IntArray{pulumi.Int(80)},
			Replicas:          pulumi.Int(frontendReplicas),
		})
		if err != nil {
			return err
		}

		if useLoadBalancer {
			ctx.Export("frontendIP", frontend.FrontendIP)
		} else {
			ctx.Export("frontendIP", frontend.Service.Spec.ApplyT(
				func(spec corev1.ServiceSpec) *string { return spec.ClusterIP }))
		}

		return nil
	})
}
