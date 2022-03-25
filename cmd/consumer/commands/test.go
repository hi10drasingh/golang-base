package commands

import (
	"github.com/droomlab/drm-coupon/internal/app/dependency"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func HandleDRMTesting(deps *dependency.Dependency) drmrmq.Handler {
	return func(msg amqp.Delivery) (action drmrmq.Action) {
		return drmrmq.NackRequeue
	}
}
