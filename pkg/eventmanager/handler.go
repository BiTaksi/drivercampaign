package eventmanager

import (
	"context"

	"github.com/BiTaksi/drivercampaign/pkg/event"
)

type EventHandler interface {
	Handle(ctx context.Context) error
}

type IEventHandlerFactory interface {
	Make(e event.Event, topic string) (EventHandler, error)
}
