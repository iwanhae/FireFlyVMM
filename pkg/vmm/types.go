package vmm

import (
	"path"
)

const (
	ROOTFS_FILE_NAME = "rootfs.ext4"
	ROOTFS_DIR_NAME  = "rootfs"

	VMLINUX_FILENAME = "vmlinux.bin"
	VMLINUX_DIR_NAME = "vmlinux"

	METADATA_FILENAME = "metadata.yaml"
)

func (v *VirtualMachineManager) getVMLinuxPath(vmlinux string) string {
	return path.Join(v.TemplateDir, VMLINUX_DIR_NAME, vmlinux, VMLINUX_FILENAME)
}
func (v *VirtualMachineManager) getRootFSPath(VMID string) string {
	return path.Join(v.DataDir, VMID, ROOTFS_FILE_NAME)
}
func (v *VirtualMachineManager) getVMMetaPath(VMID string) string {
	return path.Join(v.DataDir, VMID, METADATA_FILENAME)
}

type VMStatus string

const (
	VMStatus_Running   VMStatus = "Running"
	VMStatus_Stopped   VMStatus = "Stopped"
	VMStatus_Preparing VMStatus = "Preparing"
)
