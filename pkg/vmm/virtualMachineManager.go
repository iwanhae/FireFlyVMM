package vmm

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	DefaultKernelArgs = "ro console=ttyS0 noapic reboot=k panic=1 pci=off network-config=disabled"
)

type VirtualMachineManager struct {
	SocketDir   string
	TemplateDir string
	DataDir     string
	LogDir      string

	CNINetworkName string

	events    chan Message
	vmMap     map[ /*VMID */ string]*VirtualMachine
	vmMapLock sync.Mutex
}

type MessageHandler func(ctx context.Context, m Message) error

type Message interface {
	MetaString() string
}

func (v *VirtualMachineManager) Request(m Message) error {
	v.events <- m
	return nil
}

func (v *VirtualMachineManager) GetVM(VMID string) (VirtualMachine, error) {
	vm, err := v.getVMRef(VMID)
	if err != nil {
		return VirtualMachine{}, err
	}

	return *vm, nil
}

func (v *VirtualMachineManager) getVMRef(VMID string) (*VirtualMachine, error) {
	v.vmMapLock.Lock()
	defer v.vmMapLock.Unlock()
	vm, ok := v.vmMap[VMID]
	if ok {
		return vm, nil
	}

	vm, err := v.loadVM(VMID)
	if err != nil {
		return nil, fmt.Errorf("%q VM not found: %w", VMID, err)
	} else {
		v.vmMap[VMID] = vm
	}

	return vm, nil
}

func (v *VirtualMachineManager) loadVM(VMID string) (*VirtualMachine, error) {
	vm := newVirtualMachine(v.getVMMetaPath(VMID))

	f, err := os.OpenFile(
		path.Join(v.DataDir, VMID, METADATA_FILENAME),
		os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(&vm); err != nil {
		return nil, err
	}

	return vm, nil
}
