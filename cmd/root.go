/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/iwanhae/monolithcloud/pkg/server"
	"github.com/iwanhae/monolithcloud/pkg/vmm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "monolithcloud",
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
		CNINetworkName: "fcnet",
		KernelArgs:     vmm.DefaultKernelArgs,
	}
	ctx := context.Background()
	go vmManager.Start(ctx)
	for i := 1; i <= 3; i++ {
		vmManager.Request(&vmm.StartVMMessage{VMID: fmt.Sprintf("%05d", i)})
	}

	h := server.NewServer(server.ServerOpts{})
	if err := http.ListenAndServe(":8000", h); err != nil {
		return err
	}
	return nil
}
