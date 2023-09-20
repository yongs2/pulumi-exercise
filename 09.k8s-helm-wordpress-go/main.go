package main

import (
	"log"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// gofmt -w -s *.go; pulumi up -y
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		namespace := "default"
		chartName := "wordpress"
		chartVersion := "15.0.5"
		storageClass := "openebs-hostpath"
		wordpress, err := helmv3.NewChart(ctx, chartName, helmv3.ChartArgs{
			Chart:     pulumi.String(chartName),
			Version:   pulumi.String(chartVersion),
			Namespace: pulumi.String(namespace),
			FetchArgs: helmv3.FetchArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
			Values: pulumi.Map{
				"global": pulumi.Map{
					"storageClass": pulumi.String(storageClass),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the public IP for WordPress.
		// If namespace is default, do not add namespace in the second argument. [source](https://github.com/pulumi/pulumi-kubernetes/blob/master/sdk/go/kubernetes/helm/v3/chart.go)
		frontendIP := wordpress.GetResource("v1/Service", namespace+"/"+chartName, namespace).ApplyT(func(r interface{}) (pulumi.StringPtrOutput, error) {
			log.Printf("interface=[%T][%v]", r, r)
			if r == nil {
				return pulumi.String("").ToStringPtrOutput(), nil
			}
			svc := r.(*corev1.Service)
			return svc.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip(), nil
		})
		ctx.Export("frontendIP", frontendIP)
		return nil
	})
}
