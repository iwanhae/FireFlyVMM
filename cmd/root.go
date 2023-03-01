/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iwanhae/monolithcloud/pkg/server"
	"github.com/iwanhae/monolithcloud/pkg/vmm"
	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "firefly",
	RunE: rootRunE,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rootRunE(cmd *cobra.Command, args []string) error {
	vmManager := vmm.VirtualMachineManager{
		SocketDir:      "./_sockets",
		DataDir:        "./_data",
		LogDir:         "./_log",
		TemplateDir:    "./templates",
		CNINetworkName: "fcnet-bridge",
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		e := make(chan os.Signal, 1)
		signal.Notify(e, syscall.SIGINT)
		<-e
		cancel()
	}()
	vmManager.Start(ctx)
	go func() {
		<-ctx.Done()
		vmManager.Stop()
	}()

	vmManager.Request(&vmm.CreateVMMessage{
		Name:            "Hello",
		MemSizeMib:      2048,
		VcpuCount:       2,
		StorageGB:       50,
		RootfsTemplate:  "ubuntu_2204",
		VmlinuxTemplate: "5.10.156",
		KernelArgs:      vmm.DefaultKernelArgs,
		CloudConfig:     cloudinit.NewDefaultCloudConfig(),
	})
	vmManager.Request(&vmm.StartVMMessage{
		VMID: "00001",
	})
	time.Sleep(10 * time.Second)
	vmManager.Request(&vmm.StopVMMessage{
		VMID: "00001",
	})
	time.Sleep(10 * time.Second)
	vmManager.Request(&vmm.DeleteVMMessage{
		VMID: "00001",
	})

	h := server.NewServer(server.ServerOpts{})
	s := http.Server{Addr: ":9000", Handler: h}
	go func() {
		<-ctx.Done()
		s.Close()
	}()
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
