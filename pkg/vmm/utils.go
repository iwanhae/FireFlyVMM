package vmm

import (
	"context"
	"os"
)

func (v *VirtualMachineManager) prepareDirectories(ctx context.Context) error {
	if err := os.MkdirAll(v.SocketDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.TemplateDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.DataDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.LogDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}
