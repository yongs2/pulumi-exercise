package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	ansible "openstack-ansible-kubespray-k8s-go/ansible"
	k8scommon "openstack-ansible-kubespray-k8s-go/k8s/common"
	k8smaster "openstack-ansible-kubespray-k8s-go/k8s/master"
	osinsta "openstack-ansible-kubespray-k8s-go/osinsta"
)

// gofmt -w -s *.go; pulumi up -y
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read PrivateKey
		privateKey, err := ioutil.ReadFile("/root/.ssh/id_rsa")
		if err != nil {
			log.Printf("ReadFile.Err[%v]", err)
			return err
		}
		// Read publicKey
		publicKey, err := ioutil.ReadFile("/root/.ssh/id_rsa.pub")
		if err != nil {
			log.Printf("ReadFile.Err[%v]", err)
			return err
		}
		// Create keypair, openstack keypair list; openstack keypair show basic
		keypair, err := compute.NewKeypair(ctx, "basic", &compute.KeypairArgs{
			Name:      pulumi.String("basic"),
			PublicKey: pulumi.String(string(publicKey)),
		})
		ctx.Export("keypairName", keypair.Name)
		log.Printf("Created keypair")

		// Create Falvor
		flavor, err := compute.NewFlavor(ctx, "C1M2D20", &compute.FlavorArgs{
			Name:     pulumi.String("C1M2D20"),
			Ram:      pulumi.Int(2048),
			Vcpus:    pulumi.Int(1),
			Disk:     pulumi.Int(20),
			IsPublic: pulumi.Bool(true),
		}, pulumi.DependsOn([]pulumi.Resource{keypair}))
		if err != nil {
			log.Printf("NewFlavor.Err[%v]", err)
			return err
		}
		ctx.Export("Flavor.Id", flavor.FlavorId)
		ctx.Export("Flavor.Name", flavor.Name)
		log.Printf("Created flavor")

		nCntMaster := 1
		nCntWorker := 1
		ipAddresses := []string{
			"192.168.5.48", // k8s-master-001
			"192.168.5.59", // k8s-master-002
			"192.168.5.60", // k8s-master-003
			"192.168.5.63", // k8s-worker-001
		}
		metallbIpRange := "192.168.5.64-192.168.5.66"
		kubeSprayVersion := "v2.21.0"
		nIdxIpAddresses := 0
		imageName := "ubuntu-22.04-LTS"
		loginUserName := "ubuntu"
		openstackNetworkName := "management"
		osinstaDefaultArgs := osinsta.OpenStackInstanceArgs{
			Name:        pulumi.String(""),
			KeyPairName: keypair.Name,
			PrivateKey:  string(privateKey[:]),
			FlavorId:    flavor.ID(),
			Image:       pulumi.String(imageName),
			NetworkName: openstackNetworkName,
			FixedIpV4:   pulumi.String(""),
		}

		var resources []pulumi.Resource
		// Create instances for k8s-master
		k8sMaster := make([]*osinsta.OpenStackInstance, nCntMaster)
		k8sScript := make([]*remote.Command, nCntMaster+nCntWorker)
		nK8sScript := 0
		for i := 0; i < nCntMaster; i++ {
			name := fmt.Sprintf("k8s-master-%02d", i+1)
			osinstaArgs := osinstaDefaultArgs
			osinstaArgs.Name = pulumi.String(name)
			osinstaArgs.FixedIpV4 = pulumi.String(ipAddresses[nIdxIpAddresses])
			k8sMaster[i], err = osinsta.NewOpenStackInstance(ctx, name, &osinstaArgs, pulumi.DependsOn([]pulumi.Resource{flavor}))
			if err != nil {
				log.Printf("Created NewOpenStackInstance[%v].Err[%v]", name, err)
				return err
			}
			ctx.Export(name, k8sMaster[i].InstanceIP)
			nIdxIpAddresses += 1

			// Run install k8s script
			k8sScript[nK8sScript], err = k8scommon.RunCommand(ctx, name, k8scommon.ClusterVariable{
				PrivateKey: string(privateKey[:]),
				Host:       k8sMaster[i].InstanceIP,
				UserName:   loginUserName,
			}, pulumi.DependsOn([]pulumi.Resource{k8sMaster[i]}))
			if err != nil {
				log.Printf("Created RunCommand[%v].Err[%v]", name, err)
				return err
			}
			resources = append(resources, k8sScript[nK8sScript])
			nK8sScript += 1
		}

		// Create instances for k8s-worker
		k8sWorker := make([]*osinsta.OpenStackInstance, nCntWorker)
		for i := 0; i < nCntWorker; i++ {
			name := fmt.Sprintf("k8s-worker-%02d", i+1)
			osinstaArgs := osinstaDefaultArgs
			osinstaArgs.Name = pulumi.String(name)
			osinstaArgs.FixedIpV4 = pulumi.String(ipAddresses[nIdxIpAddresses])
			k8sWorker[i], err = osinsta.NewOpenStackInstance(ctx, name, &osinstaArgs, pulumi.DependsOn([]pulumi.Resource{flavor}))
			if err != nil {
				log.Printf("Created NewOpenStackInstance[%v].Err[%v]", name, err)
				return err
			}
			ctx.Export(name, k8sWorker[i].InstanceIP)
			nIdxIpAddresses += 1

			// Run install k8s script
			k8sScript[nK8sScript], err = k8scommon.RunCommand(ctx, name, k8scommon.ClusterVariable{
				PrivateKey: string(privateKey[:]),
				Host:       k8sWorker[i].InstanceIP,
				UserName:   loginUserName,
			}, pulumi.DependsOn([]pulumi.Resource{k8sWorker[i]}))
			if err != nil {
				log.Printf("Created RunCommand[%v].Err[%v]", name, err)
				return err
			}
			resources = append(resources, k8sScript[nK8sScript])
			nK8sScript += 1
		}

		// Create instance for ansible
		var ansible1 *osinsta.OpenStackInstance
		var scriptAnsible *remote.Command
		for i := 0; i < 1; i++ {
			name := "ansible"
			osinstaArgs := osinstaDefaultArgs
			osinstaArgs.Name = pulumi.String(name)
			osinstaArgs.FixedIpV4 = pulumi.String(ipAddresses[nIdxIpAddresses])
			ansible1, err = osinsta.NewOpenStackInstance(ctx, name, &osinstaArgs, pulumi.DependsOn(append(resources, flavor)))
			if err != nil {
				log.Printf("Created NewOpenStackInstance[%v].Err[%v]", name, err)
				return err
			}
			nIdxIpAddresses += 1
			// Export the IP of the instance
			ctx.Export("AnsibleIP", ansible1.InstanceIP)

			// Run kubespray script
			scriptAnsible, err = ansible.RunCommand(ctx, name, ansible.ClusterVariable{
				PrivateKey:       string(privateKey[:]),
				Host:             ansible1.InstanceIP,
				UserName:         loginUserName,
				HttpProxy:        "",
				MasterIps:        strings.Join(ipAddresses[0:nCntMaster], ","),
				WorkerIps:        strings.Join(ipAddresses[nCntMaster:nCntMaster+nCntWorker], ","),
				KubeControlHosts: len(ipAddresses[0:nCntMaster]),
				MetallbIpRange:   metallbIpRange,
				KubesprayVersion: kubeSprayVersion,
			}, pulumi.DependsOn([]pulumi.Resource{ansible1}))
			if err != nil {
				log.Printf("Created RunCommand.Err[%v]", err)
				return err
			}
		}

		// Run k8s master script after all ntels script have run
		for i := 0; i < nCntMaster; i++ {
			name := fmt.Sprintf("k8s-master-%02d", i+1)
			_, err = k8smaster.RunCommand(ctx, name, k8smaster.ClusterVariable{
				PrivateKey: string(privateKey[:]),
				Host:       k8sMaster[i].InstanceIP,
				UserName:   loginUserName,
			}, pulumi.DependsOn([]pulumi.Resource{scriptAnsible}))
			if err != nil {
				log.Printf("Created RunCommand[%v].Err[%v]", name, err)
				return err
			}
		}

		return nil
	})
}
