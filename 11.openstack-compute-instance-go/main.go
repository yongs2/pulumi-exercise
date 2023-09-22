package main

import (
	"io/ioutil"
	"log"

	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/compute"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// gofmt -w -s *.go; pulumi up -y
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read publicKey
		publicKey, err := ioutil.ReadFile("/root/.ssh/id_rsa.pub")
		if err != nil {
			return err
		}
		// Create keypair, openstack keypair list; openstack keypair show basic
		keypair, err := compute.NewKeypair(ctx, "basic", &compute.KeypairArgs{
			Name:      pulumi.String("basic"),
			PublicKey: pulumi.String(string(publicKey)),
		})
		ctx.Export("keypairName", keypair.Name)

		// [Get Network](https://www.pulumi.com/registry/packages/openstack/api-docs/networking/getnetwork/)
		// openstack network show management --format=json
		network, err := networking.LookupNetwork(ctx, &networking.LookupNetworkArgs{
			Name: pulumi.StringRef("management"),
		}, nil)
		if err != nil {
			return err
		}
		ctx.Export("network.Id", pulumi.String(network.Id))
		ctx.Export("network.Name", pulumi.String(*network.Name))

		// [Get Subnet](https://www.pulumi.com/registry/packages/openstack/api-docs/networking/getsubnet/)
		// openstack subnet show management --format=json
		subnet, err := networking.LookupSubnet(ctx, &networking.LookupSubnetArgs{
			Name: pulumi.StringRef("management"),
		}, nil)
		if err != nil {
			return err
		}
		ctx.Export("subnet.Id", pulumi.String(subnet.Id))

		// Create Port
		// openstack port create --network management --fixed-ip subnet=management,ip-address=192.168.5.48 basic
		newPort, err := networking.NewPort(ctx, "basic", &networking.PortArgs{
			Name:         pulumi.String("basic"),
			NetworkId:    pulumi.String(network.Id),
			AdminStateUp: pulumi.Bool(true),
			FixedIps: networking.PortFixedIpArray{
				networking.PortFixedIpArgs{
					SubnetId:  pulumi.String(subnet.Id),
					IpAddress: pulumi.String("192.168.5.48"),
				},
			},
		})
		if err != nil {
			return err
		}
		_ = newPort.Name.ApplyT(func(v interface{}) string {
			log.Printf("newPort.Name=[%v]", v.(string))
			return v.(string)
		}).(pulumi.StringOutput)
		ctx.Export("newPort.Name", newPort.Name)

		// Create Falvor
		flavor, err := compute.NewFlavor(ctx, "C1M2D20", &compute.FlavorArgs{
			Name:     pulumi.String("C1M2D20"),
			Ram:      pulumi.Int(2048),
			Vcpus:    pulumi.Int(1),
			Disk:     pulumi.Int(20),
			IsPublic: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}
		ctx.Export("Flavor.Id", flavor.FlavorId)
		ctx.Export("Flavor.Name", flavor.Name)

		// Create an OpenStack resource (Compute Instance)
		// openstack server create --flavor "C1M2D20" --image ubuntu-20.04.5 --key-name packstack --port basic basic
		instance, err := compute.NewInstance(ctx, "basic", &compute.InstanceArgs{
			Name:      pulumi.String("basic"),
			ImageName: pulumi.String("ubuntu-20.04.5"),
			FlavorId:  flavor.FlavorId,
			KeyPair:   keypair.Name,
			Networks: compute.InstanceNetworkArray{
				compute.InstanceNetworkArgs{
					Port: newPort.ID(), // ID() in CustomResourceState (pulimi/sdk/go/pulumi/resource.go)
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the IP of the instance
		ctx.Export("instanceIP", instance.AccessIpV4)
		return nil
	})
}
