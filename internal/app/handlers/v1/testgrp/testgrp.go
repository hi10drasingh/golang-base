package testgrp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app/dependency"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/pkg/errors"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Handlers holds local dependencies for handler.
type Handlers struct {
	Slug  string
	Log   drmlog.Logger
	SQL   *nap.DB
	NoSQL *mongo.Client
	RMQ   *drmrmq.RabbitMQ
}

// NewHandlers return new instance of handler.
func NewHandlers(deps *dependency.Dependency) *Handlers {
	return &Handlers{
		Slug:  "test",
		Log:   deps.Log,
		SQL:   deps.SQL,
		NoSQL: deps.NoSQL,
		RMQ:   deps.RMQ,
	}
}

// Hello is sample route handler.
func (h Handlers) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pub, err := h.RMQ.NewPublisher(ctx)
	if err != nil {
		return errors.Wrap(err, "Create New Publisher")
	}

	body, err := json.Marshal(struct {
		Code string            `json:"code"`
		Data map[string]string `json:"data"`
	}{
		Code: "success",
		Data: map[string]string{
			"name":    "hitendra",
			"age":     "24",
			"testing": "true",
		},
	})
	if err != nil {
		return errors.Wrap(err, "JSON Marshal Body")
	}

	err = pub.Publish(ctx, &drmrmq.PublisherConfig{
		Exchange: "drmrmq_testing_exc",
		Key:      "drmrmq_testing_key",
		Msg: amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		},
	})

	if err != nil {
		return errors.Wrap(err, "Msg Publish")
	}

	fmt.Fprint(w, "hello, world!\n")

	return nil
}
