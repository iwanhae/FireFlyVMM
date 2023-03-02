package server

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/iwanhae/monolithcloud/pkg/vmm"
	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
	"github.com/labstack/echo/v4"
	"github.com/nxadm/tail"
)

type Server struct {
	vmm *vmm.VirtualMachineManager
}

var _ ServerInterface = &Server{}

// CreateVM implements ServerInterface
func (s *Server) CreateVM(c echo.Context) error {
	req := VMSpec{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	cfg := cloudinit.NewDefaultCloudConfig()
	cfg.Users = []cloudinit.UserCoinfig{
		{Name: "root", SSHAuthorizedKeys: req.SshAuthorizedKeys},
	}

	if err := s.vmm.Request(&vmm.CreateVMMessage{
		Name:            req.Name,
		MemSizeMib:      int64(req.MemorySizeMb),
		VcpuCount:       int64(req.VcpuCount),
		StorageGB:       int64(req.StorageGb),
		RootfsTemplate:  req.Rootfs,
		VmlinuxTemplate: req.Vmlinux,
		KernelArgs:      vmm.DefaultKernelArgs,
		CloudConfig:     cfg,
	}); err != nil {
		return err
	}

	return c.String(http.StatusAccepted, "accepted")
}

// GetVM implements ServerInterface
func (s *Server) GetVM(c echo.Context, vmId int) error {
	vm, err := s.vmm.GetVM(fmt.Sprintf("%05d", vmId))
	if err != nil {
		return err
	}
	return c.JSON(200, vm)
}

// ListVMs implements ServerInterface
func (s *Server) ListVMs(c echo.Context) error {
	ctx := c.Request().Context()
	vmList, err := s.vmm.ListVMs(ctx)
	if err != nil {
		return err
	}
	return c.JSON(200, vmList)
}

// SetVMStatus implements ServerInterface
func (s *Server) SetVMStatus(c echo.Context, vmId int) error {
	var status DesiredVMStatus
	if b, err := io.ReadAll(c.Request().Body); err != nil {
		return err
	} else {
		status = DesiredVMStatus(b)
	}

	switch status {
	case Stop:
		s.vmm.Request(&vmm.StopVMMessage{VMID: fmt.Sprintf("%05d", vmId), FroceStop: false})
	case ForceStop:
		s.vmm.Request(&vmm.StopVMMessage{VMID: fmt.Sprintf("%05d", vmId), FroceStop: true})
	case Run:
		s.vmm.Request(&vmm.StartVMMessage{VMID: fmt.Sprintf("%05d", vmId)})
	default:
		return fmt.Errorf("unhandled status: %q", status)
	}
	return c.String(http.StatusAccepted, "accepted")
}

// GetVMStatus implements ServerInterface
func (s *Server) GetVMStatus(c echo.Context, vmId int) error {
	vm, err := s.vmm.GetVM(fmt.Sprintf("%05d", vmId))
	if err != nil {
		return err
	}
	if vm.IsRunning(c.Request().Context()) {
		return c.String(200, string(Running))
	} else {
		return c.String(200, string(Stoped))
	}
}

// DeleteVM implements ServerInterface
func (s *Server) DeleteVM(c echo.Context, vmId int) error {
	ctx := c.Request().Context()
	vm, err := s.vmm.GetVM(fmt.Sprintf("%05d", vmId))
	if err != nil {
		return err
	}
	if vm.IsRunning(ctx) {
		return c.JSON(500, Error{Code: "-1", Message: "can not delete running vm"})
	}
	if err := s.vmm.Request(&vmm.DeleteVMMessage{VMID: vm.ID}); err != nil {
		return err
	}
	return c.String(http.StatusAccepted, "accepted")
}

// GetVMLog implements ServerInterface
func (s *Server) GetVMLog(c echo.Context, vmId int) error {
	_, err := s.vmm.GetVM(fmt.Sprintf("%05d", vmId))
	if err != nil {
		return err
	}
	t, err := tail.TailFile(
		path.Join(s.vmm.LogDir, fmt.Sprintf("%05d_stdout.log", vmId)), // I'm lazy person. this is hard coded
		tail.Config{Follow: true, ReOpen: true, Poll: true})
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
out:
	for {
		select {
		case <-c.Request().Context().Done():
			break out
		case line := <-t.Lines:
			fmt.Fprintf(c.Response(), "%s\n", line.Text)
			c.Response().Flush()
		}
	}
	t.Cleanup()
	return nil
}
