package vmm

import (
	"context"
	"fmt"
)

var _ Message = &StopVMMessage{}

type StopVMMessage struct {
	VMID      string
	FroceStop bool
}

func (*StopVMMessage) MetaString() string {
	return "Stop VM"
}

func (v *VirtualMachineManager) stopVM(ctx context.Context, msg *StopVMMessage) error {
	vm, err := v.getVMRef(msg.VMID)
	if err != nil {
		return err
	}
	if !vm.IsRunning(ctx) {
		return fmt.Errorf("VM is already terminated")
	}
	if msg.FroceStop {
		return vm.m.StopVMM()
	} else {
		return vm.m.Shutdown(ctx)
	}
}
