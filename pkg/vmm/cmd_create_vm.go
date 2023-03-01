package vmm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
)

var _ Message = &CreateVMMessage{}

type CreateVMMessage struct {
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

func (v *VirtualMachineManager) CreateVM(ctx context.Context, msg *CreateVMMessage) error {
	id, err := v.generateVMID()
	if err != nil {
		return fmt.Errorf("failed to generate a new VMID: %w", err)
	}

	// Defaulting
	msg.CloudConfig.Hostname = msg.Name
	if msg.KernelArgs == "" {
		msg.KernelArgs = DefaultKernelArgs
	}

	// Create Base Dir for VM
	dataPath := path.Join(v.DataDir, id)
	if err := os.Mkdir(dataPath, 0755); err != nil {
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

func (v *VirtualMachineManager) generateVMID() (string, error) {
	entries, err := os.ReadDir(v.DataDir)
	if err != nil {
		return "", err
	}
	max := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		v, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if max < v {
			max = v
		}
	}
	return fmt.Sprintf("%05d", max+1), nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
