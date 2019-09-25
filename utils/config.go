package utils

import (
	"flag"
	"github.com/webitel/cdr_metrics/model"
)

var (
	hostGateway     = flag.String("host", "http://localhost", "Host to push gateway")
	usernameGateway = flag.String("username", "", "Auth name")
	passwordGateway = flag.String("password", "", "Password")
	amqp            = flag.String("amqp", "", "Host to AMQP")
	namespace       = flag.String("space", "", "Space")
)

func LoadConfig() (model.Config, error) {
	flag.Parse()
	conf := model.Config{
		Space: *namespace,
		MessageQuery: model.MessageQueryConfig{
			DataSource: *amqp,
		},
		Gateway: model.GatewayConfig{
			Host:     *hostGateway,
			Username: *usernameGateway,
			Password: *passwordGateway,
		},
	}

	return conf, nil
}
