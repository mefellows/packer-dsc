TEST?=./...

default: test

bin:
	@sh -c "$(CURDIR)/scripts/build.sh"

dev:
	go build -o "bin/provisioner-dsc" ./plugin/provisioner-dsc

test:
	"$(CURDIR)/scripts/test.sh"

testrace:
	go test -race $(TEST) $(TESTARGS)

deps:
	go install github.com/hashicorp/packer/cmd/mapstructure-to-hcl2
	go generate github.com/mefellows/packer-dsc/...

.PHONY: bin default dev test deps
