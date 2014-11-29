package main

import (
	"flag"
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm"
	"github.com/mefellows/packer-dsc/provisioner/dsc"
	"os"
)

func main() {
	var (
		hostname string
		user     string
		pass     string
		cmd      string
		port     int
	)

	flag.StringVar(&hostname, "hostname", "localhost", "winrm host")
	flag.StringVar(&user, "username", "vagrant", "winrm admin username")
	flag.StringVar(&pass, "password", "vagrant", "winrm admin password")
	flag.IntVar(&port, "port", 5985, "winrm port")
	flag.Parse()

	cmd = flag.Arg(0)
	client := winrm.NewClient(&winrm.Endpoint{hostname, port}, user, pass)
	err := client.RunWithInput(cmd, os.Stdout, os.Stderr, os.Stdin)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
