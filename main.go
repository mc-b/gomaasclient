package main

import (
	"time"

	gomaasclient "github.com/ionutbalutoiu/gomaasclient/client"
	gomaasentity "github.com/ionutbalutoiu/gomaasclient/entity"
)

func main() {

	client, _ := gomaasclient.GetClient("http://10.6.37.8:5240/MAAS", "nC83nVyLDWKF8zxvPq:VnYukvR9Yh3w9jRDUe:FdU7sFWJu8DHjAeRumRNmZfasNBDqgXa", "2.0")

	// List MAAS machines
	machines, _ := client.Machines.Get()

	// List MAAS VM hosts
	vmHosts, _ := client.VMHosts.Get()

	for _, m := range machines {
		print(m.FQDN + "\n")
	}

	var last int
	for _, host := range vmHosts {
		print(host.Host.SystemID + "\n")
		last = host.ID
	}
	params := &gomaasentity.VMHostMachineParams{
		Cores:    2,
		Hostname: "order-62",
		Memory:   4096,
		Storage:  "16",
	}
	machine, err := client.VMHost.Compose(last, params)
	if err != nil {
		print(err)
		return
	}

	print(machine.SystemID)

	for {
		m, err := client.Machine.Get(machine.SystemID)
		if err != nil {
			return
		}
		if m.StatusName == "Ready" {
			break
		}
		print(".")
		time.Sleep(10 * time.Second)
	}

	mparams := &gomaasentity.MachineParams{
		Pool: "webshop",
		Zone: "10-6-37-0",
	}

	mpower := make(map[string]string)
	mpower[machine.PowerType] = "virsh"

	client.Machine.Update(machine.SystemID, mparams, mpower)

	userdata := `#cloud-config
packages:
  - nginx
`

	deploy := &gomaasentity.MachineDeployParams{
		UserData: userdata,
	}

	client.Machine.Deploy(machine.SystemID, deploy)

}
