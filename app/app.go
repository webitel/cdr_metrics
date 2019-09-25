package app

import (
	"github.com/webitel/cdr_metrics/gateway"
	"github.com/webitel/cdr_metrics/model"
	"github.com/webitel/cdr_metrics/mq"
	"github.com/webitel/cdr_metrics/mq/rabbit"
	"github.com/webitel/cdr_metrics/utils"
	"github.com/webitel/wlog"
	"sync"
)

type App struct {
	Log       *wlog.Logger
	config    model.Config
	mq        mq.MQ
	startOnce sync.Once
	stop      chan struct{}
	listener  *Listener
	gateway   *gateway.Gateway
}

func NewApp(options ...string) (*App, error) {
	app := &App{}

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}
	app.config = config

	app.Log = wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  wlog.LevelDebug,
	})

	wlog.RedirectStdLog(app.Log)
	wlog.InitGlobalLogger(app.Log)

	app.mq = rabbit.NewConnection(app.Config().MessageQuery.DataSource)
	if err = app.mq.Start(); err != nil {
		return nil, err
	}

	gateway := gateway.NewGateway(app.Config().Gateway.Host, app.Config().Gateway.Username, app.Config().Gateway.Password)
	app.gateway = gateway

	app.listener = NewListener(app)

	app.startOnce.Do(func() {
		go app.listener.ListenEvents()
	})

	return app, nil
}

func (app *App) Shutdown() {
	app.mq.Close()
	app.listener.Stop()
}

func (app *App) Config() model.Config {
	return app.config
}
