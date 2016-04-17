# Packer DSC Provisioner Examples

## AWS AMI

Update the vars as required:

```
cd examples
PACKER_LOG=1 PACKER_LOG_PATH=packer.log packer build \
  -var source_ami=ami-75ebe61f \
  -var subnet_id=subnet-7ed32427 \
  -var vpc_id=vpc-89388dec \
  -var region=us-east-1 \
  ./packer.aws.json
```

## Virtualbox OVF
Update the path below to an existing OVF file:

```
cd examples
PACKER_LOG=1 PACKER_LOG_PATH=packer.log packer build \
  -debug \
  -var ovf_source_path=/Users/mfellows/Downloads/output-virtualbox-iso/packer-virtualbox-iso-1417096689.ovf \
  ./packer.json
```
