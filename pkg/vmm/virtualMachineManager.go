package vmm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
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

	events chan Message
	vmMap  map[ /*VMID */ string]*VirtualMachine
}

type MessageHandler func(ctx context.Context, m Message) error

type Message interface {
	MetaString() string
}

func (v *VirtualMachineManager) Request(m Message) error {
	v.events <- m
	return nil
}

func (v *VirtualMachineManager) loadVM(VMID string) (*VirtualMachine, error) {
	vm := newVirtualMachine(v.getVMMetaPath(VMID))

	f, err := os.OpenFile(
		path.Join(v.DataDir, VMID, METADATA_FILENAME),
		os.O_RDONLY, 0755)
	if err != nil {
		return vm, err
	}

	if err := yaml.NewDecoder(f).Decode(&vm); err != nil {
		return vm, err
	}

	v.vmMap[VMID] = vm
	return vm, nil
}

func (v *VirtualMachineManager) GetVM(VMID string) (*VirtualMachine, error) {
	vm, ok := v.vmMap[VMID]
	if ok {
		return vm, nil
	}

	vm, err := v.loadVM(VMID)
	if err != nil {
		return vm, fmt.Errorf("%q VM not found: %w", VMID, err)
	}

	return vm, nil
}

func (v *VirtualMachineManager) Start(ctx context.Context) error {
	log.Info().Msg("VMM is started")
	defer log.Info().Msg("VMM is terminated")

	if err := v.prepareDirectories(ctx); err != nil {
		return fmt.Errorf("fail to create directories: %w", err)
	}

	if v.vmMap == nil {
		v.vmMap = make(map[string]*VirtualMachine)
	}

	if v.events != nil {
		return fmt.Errorf("already started")
	}
	v.events = make(chan Message, 10)

	go v.loop(ctx)

	return nil
}

func (v *VirtualMachineManager) Stop() {
	close(v.events)
	v.events = nil
}
