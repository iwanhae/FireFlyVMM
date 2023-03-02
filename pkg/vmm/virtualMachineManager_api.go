package vmm

import (
	"context"
)

func (v *VirtualMachineManager) ListVMs(ctx context.Context) ([]VirtualMachine, error) {
	v.vmMapLock.Lock()
	defer v.vmMapLock.Unlock()

	tmp := []VirtualMachine{}
	for _, val := range v.vmMap {
		tmp = append(tmp, *val)
	}
	return tmp, nil
}
