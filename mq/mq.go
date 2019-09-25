package mq

import (
	"github.com/webitel/cdr_metrics/model"
)

type MQ interface {
	ConsumeCdr() <-chan *model.Cdr
	Start() error
	Close()
}
