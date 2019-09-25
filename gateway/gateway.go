package gateway

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Collector prometheus.Collector

const JOB_NAME = "cdr_monitoring"
const GROUPING_KEY_INSTANCE = "instance"

type Gateway struct {
	host     string
	username string
	password string
	pusher   *push.Pusher
}

func NewGateway(host, username, password string) *Gateway {
	pusher := push.New(host, JOB_NAME)

	if username != "" {
		pusher.BasicAuth(username, password)
	}

	return &Gateway{
		host:     host,
		username: username,
		password: password,
		pusher:   pusher,
	}
}

func (g *Gateway) Pusher() *push.Pusher {
	pusher := push.New(g.host, JOB_NAME)

	if g.username != "" {
		pusher.BasicAuth(g.username, g.password)
	}

	return pusher
}

func (g *Gateway) Push(instance string, metrics ...Collector) error {

	pusher := g.Pusher()
	for _, m := range metrics {
		pusher.Collector(m)
	}

	pusher.Grouping(GROUPING_KEY_INSTANCE, instance)
	//pusher.Grouping(GROUPING_KEY_IP, g.ip)
	//pusher.Grouping(GROUPING_KEY_OS, utils.OsVersion())
	return pusher.Add()
}
