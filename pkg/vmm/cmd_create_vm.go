package vmm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
)

var _ Message = &CreateVMMessage{}

type CreateVMMessage struct {
	ID *string

	Name       string
	MemSizeMib int64
	VcpuCount  int64
	StorageGB  int64

	CloudConfig cloudinit.CloudConfig

	RootfsTemplate  string
	VmlinuxTemplate string

	KernelArgs string
}

func (*CreateVMMessage) MetaString() string {
	return "Create VM"
}

func (v *VirtualMachineManager) createVM(ctx context.Context, msg *CreateVMMessage) error {
	// Defaulting
	msg.CloudConfig.Hostname = msg.Name
	if msg.KernelArgs == "" {
		msg.KernelArgs = DefaultKernelArgs
	}

	// ID Validation
	id := ""
	var err error
	if msg.ID == nil {
		id, err = v.generateVMID()
		if err != nil {
			return fmt.Errorf("failed to generate a new VMID: %w", err)
		}
	} else {
		id = *msg.ID
	}
	if vm, _ := v.loadVM(id); vm != nil {
		return fmt.Errorf("VM with specified id (%s) already exists", id)
	}

	// Create Base Dir for VM
	dataPath := path.Join(v.DataDir, id)
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("failed to generate a base dir: %w", err)
	}

	// COPY rootfs to datadir
	rootfsPath := path.Join(v.DataDir, id, ROOTFS_FILE_NAME)
	if _, err := copy(
		path.Join(v.TemplateDir, ROOTFS_DIR_NAME, msg.RootfsTemplate, ROOTFS_FILE_NAME),
		rootfsPath,
	); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	// Resize rootfs
	if err := os.Truncate(rootfsPath, msg.StorageGB*1024 /*GB*/ *1024 /*MB*/ *1024 /*KB*/); err != nil {
		return fmt.Errorf("fail to resize rootfs: %w", err)
	}

	// Create Cloud Config
	vm := &VirtualMachine{
		ID:          id,
		Name:        msg.Name,
		MemSizeMib:  msg.MemSizeMib,
		VcpuCount:   msg.VcpuCount,
		VMLinux:     msg.VmlinuxTemplate,
		CloudConfig: msg.CloudConfig,
		KernelArgs:  msg.KernelArgs,
		vmMetaPath:  v.getVMMetaPath(id),
	}
	if err := vm.save(ctx); err != nil {
		return fmt.Errorf("failed to update vm meta: %w", err)
	}
	return nil
}
