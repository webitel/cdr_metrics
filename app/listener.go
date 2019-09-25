package app

import (
	"fmt"
	"github.com/webitel/cdr_metrics/metrics"
	"github.com/webitel/wlog"
)

type Listener struct {
	app       *App
	cdrMetric *metrics.MetricCdr
	stop      chan struct{}
	stopped   chan struct{}
}

func NewListener(app *App) *Listener {
	return &Listener{
		app:       app,
		stop:      make(chan struct{}),
		stopped:   make(chan struct{}),
		cdrMetric: metrics.NewCdr(app.Config().Space, app.gateway),
	}
}

func (l *Listener) Stop() {
	close(l.stop)
	<-l.stopped
	wlog.Debug("stopped listener")
}

func (l *Listener) ListenEvents() {
	wlog.Debug("listen events")
	defer close(l.stopped)
	for {

		select {
		case m := <-l.app.mq.ConsumeCdr():
			wlog.Debug(fmt.Sprintf("receive call uuid=%s user_agent=%s", m.Uuid(), m.UserAgent()))
			l.cdrMetric.Push(m)
		case <-l.stop:
			return
		}
	}
}
