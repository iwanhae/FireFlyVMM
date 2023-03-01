/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwanhae/monolithcloud/pkg/server"
	"github.com/iwanhae/monolithcloud/pkg/vmm"
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
		KernelArgs:     vmm.DefaultKernelArgs,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		e := make(chan os.Signal, 1)
		signal.Notify(e, syscall.SIGINT)
		<-e
		cancel()
	}()
	go vmManager.Start(ctx)
	go func() {
		<-ctx.Done()
		vmManager.Stop()
	}()

	for i := 1; i <= 3; i++ {
		vmManager.Request(&vmm.StartVMMessage{VMID: fmt.Sprintf("%05d", i)})
	}

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
