package vmm

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func (v *VirtualMachineManager) loop(ctx context.Context) {
	defer log.Info().Msg("VMM is terminated")
	for event := range v.events {
		log.Info().Str("event", event.MetaString()).Msg("new event received")
		err := func(event Message) error {
			switch evt := event.(type) {
			case *CreateVMMessage:
				return v.createVM(ctx, evt)
			case *StartVMMessage:
				return v.startVM(ctx, evt)
			case *StopVMMessage:
				return v.stopVM(ctx, evt)
			case *DeleteVMMessage:
				return v.deleteVM(ctx, evt)
			default:
				return fmt.Errorf("unhandled message")
			}
		}(event)

		if err != nil {
			log.Error().Err(err).Send()
		} else {
			log.Info().Msg("event handled")
		}
	}
}

func (v *VirtualMachineManager) Start(ctx context.Context) error {
	log.Info().Msg("VMM is started")

	if err := v.prepareDirectories(ctx); err != nil {
		return fmt.Errorf("fail to create directories: %w", err)
	}

	if v.vmMap == nil {
		v.vmMap = make(map[string]*VirtualMachine)
	}

	if v.events != nil {
		return fmt.Errorf("already started")
	}

	entries, err := os.ReadDir(v.DataDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		_, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		if _, err := v.getVMRef(name); err != nil {
			return fmt.Errorf("VM metadata is corrupted: %s", name)
		} else {
			log.Debug().Str("VMID", name).Msg("founded")
		}
	}

	v.events = make(chan Message, 10)

	go v.loop(ctx)

	return nil
}

func (v *VirtualMachineManager) Stop(ctx context.Context) error {
	close(v.events)
	v.events = nil
	v.vmMapLock.Lock()
	defer v.vmMapLock.Unlock()

	for id, vm := range v.vmMap {
		if !vm.IsRunning(ctx) {
			continue
		}
		if err := vm.m.Shutdown(ctx); err != nil {
			return err
		}
		log.Info().Str("VMID", id).Msg("try stopping VM")
	}
	for id, vm := range v.vmMap {
		if vm.m != nil {
			vm.m.Wait(ctx)
		}
		log.Info().Str("VMID", id).Msg("VM is stopped")
	}
	return nil
}

func (v *VirtualMachineManager) Kill(ctx context.Context) error {
	close(v.events)
	v.events = nil
	v.vmMapLock.Lock()
	defer v.vmMapLock.Unlock()

	for id, vm := range v.vmMap {
		if !vm.IsRunning(ctx) {
			continue
		}
		if err := vm.m.StopVMM(); err != nil {
			return err
		}
		log.Info().Str("VMID", id).Msg("force stopping VM")
	}
	return nil
}

func (v *VirtualMachineManager) IsRunning(ctx context.Context) bool {
	if v.events != nil {
		return true
	}
	v.vmMapLock.Lock()
	defer v.vmMapLock.Unlock()
	for _, vm := range v.vmMap {
		if vm.IsRunning(ctx) {
			return true
		}
	}
	return false
}
