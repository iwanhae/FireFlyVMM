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

	"github.com/iwanhae/monolithcloud/pkg/server"
	"github.com/iwanhae/monolithcloud/pkg/vmm"
	"github.com/rs/zerolog/log"
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

	vmManager.Start(ctx)

	h := server.NewServer(server.ServerOpts{VMM: &vmManager})
	server := http.Server{Addr: ":9000", Handler: h}
	go func() {
		ctx := context.Background() // Terminating Context
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		switch s := <-c; {
		case s == syscall.SIGTERM || s == os.Interrupt:
			if err := vmManager.Stop(ctx); err != nil {
				log.Error().Err(err).Msg("fail to stop VMManager")
			}
			if !vmManager.IsRunning(ctx) {
				if err := server.Shutdown(ctx); err != nil {
					log.Error().Err(err).Msg("fail to stop Server")
				}
			}
		case s == syscall.SIGQUIT:
			if err := vmManager.Kill(ctx); err != nil {
				log.Error().Err(err).Msg("fail to stop VMManager")
			}
			if err := server.Close(); err != nil {
				log.Error().Err(err).Msg("fail to stop Server")
			}
		}
	}()
	log.Info().Msgf("listen on %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
