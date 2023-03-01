package vmm

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (v *VirtualMachineManager) loop(ctx context.Context) {
	for event := range v.events {
		log.Info().Str("event", event.MetaString()).Msg("new event received")
		err := func(event Message) error {
			switch evt := event.(type) {
			case *CreateVMMessage:
				return v.CreateVM(ctx, evt)
			case *StartVMMessage:
				return v.StartVM(ctx, evt)
			case *StopVMMessage:
				return v.StopVM(ctx, evt)
			case *DeleteVMMessage:
				return v.DeleteVM(ctx, evt)
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
