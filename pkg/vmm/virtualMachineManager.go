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
	SocketDir      string
	TemplateDir    string
	DataDir        string
	LogDir         string
	CNINetworkName string
	KernelArgs     string

	events chan Message
	vmMap  map[ /*VMID */ string]*VirtualMachine
}

type MessageHandler func(ctx context.Context, m Message) error

type Message interface {
	MetaString() string
}

func (v *VirtualMachineManager) Request(m Message) error {
	if v.events == nil {
		v.events = make(chan Message, 10)
	}
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

func (v *VirtualMachineManager) Start(ctx context.Context) {
	log.Info().Msg("VMM is started")
	defer log.Info().Msg("VMM is terminated")

	if v.vmMap == nil {
		v.vmMap = make(map[string]*VirtualMachine)
	}

	for event := range v.events {
		log.Info().Str("event", event.MetaString()).Msg("new event received")
		err := func(event Message) error {
			switch evt := event.(type) {
			case *CreateVMMessage:
				return v.CreateVM(ctx, evt)
			case *StartVMMessage:
				return v.StartVM(ctx, evt)
			default:
				return fmt.Errorf("unhandled message")
			}
		}(event)

		if err != nil {
			log.Error().Err(err).Send()
		} else {
			log.Info().Msg("event handled")
		}
	}
}

func (v *VirtualMachineManager) Stop() {
	close(v.events)
}
