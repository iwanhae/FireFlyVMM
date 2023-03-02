package vmm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func newVirtualMachine(vmMetaPath string) *VirtualMachine {
	return &VirtualMachine{vmMetaPath: vmMetaPath}
}

type VirtualMachine struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip"`

	Name       string `json:"name"`
	MemSizeMib int64  `json:"memory_size_mb"`
	VcpuCount  int64  `json:"vcpu_count"`
	VMLinux    string `json:"vmlinux"`

	CloudConfig cloudinit.CloudConfig `json:"cloud_config"`
	KernelArgs  string                `json:"kernel_args"`

	vmMetaPath string               `json:"-"`
	m          *firecracker.Machine `json:"-"`
	stdout     io.WriteCloser       `json:"-"`
	stderr     io.WriteCloser       `json:"-"`
	//stdin  io.ReadCloser
}

func (vm *VirtualMachine) IsRunning(ctx context.Context) bool {
	if vm.m == nil {
		return false
	}
	// TODO: chekc status
	info, err := vm.m.DescribeInstanceInfo(ctx)
	if err != nil {
		return false
	}
	if info.State == nil {
		return false
	}
	if *info.State == "Running" {
		return true
	}
	return false
}

func (vm *VirtualMachine) SocketPath(socketDir string) string {
	return path.Join(socketDir, fmt.Sprintf("%s.sock", vm.ID))
}

func (vm *VirtualMachine) save(ctx context.Context) error {
	if vm.vmMetaPath == "" {
		return fmt.Errorf("no vmMetaPath set")
	}
	f, err := os.OpenFile(
		vm.vmMetaPath,
		os.O_WRONLY|os.O_CREATE,
		0755,
	)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewEncoder(f).Encode(vm); err != nil {
		return fmt.Errorf("failed to marshal vm meta: %w", err)
	}
	return nil
}

func (vm *VirtualMachine) watch(ctx context.Context, m *firecracker.Machine) {
	vm.m = m
	for _, ni := range vm.m.Cfg.NetworkInterfaces {
		vm.IPAddress = ni.StaticConfiguration.IPConfiguration.IPAddr.IP.String()
	}
	vm.save(ctx)

	if err := m.Wait(ctx); err != nil {
		log.Error().Err(err).
			Str("id", vm.ID).
			Msg("VM stopped")
	} else {
		log.Info().
			Str("id", vm.ID).
			Msg("VM stopped")
	}
	m.StopVMM()
	vm.stdout.Close()
	vm.stderr.Close()
}
