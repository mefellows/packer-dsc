package main

import (
	"github.com/mefellows/packer-dsc/provisioner/dsc"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {

	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterProvisioner(new(dsc.Provisioner))
	server.Serve()
}
