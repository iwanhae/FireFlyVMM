package vmm

import (
	"context"
	"fmt"
	"os"
	"path"
)

var _ Message = &DeleteVMMessage{}

type DeleteVMMessage struct {
	VMID string
}

func (*DeleteVMMessage) MetaString() string {
	return "Delete VM"
}

func (v *VirtualMachineManager) deleteVM(ctx context.Context, msg *DeleteVMMessage) error {
	vm, err := v.getVMRef(msg.VMID)
	if err != nil {
		return nil
	}
	if vm.IsRunning(ctx) {
		return fmt.Errorf("can not delete on running vm")
	}
	if err := os.RemoveAll(
		path.Join(v.DataDir, msg.VMID),
	); err != nil {
		return err
	}
	delete(v.vmMap, msg.VMID)
	return nil
}
