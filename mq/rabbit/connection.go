package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/webitel/cdr_metrics/model"
	"github.com/webitel/cdr_metrics/mq"
	"github.com/webitel/wlog"
	"os"
	"time"
)

const (
	MAX_ATTEMPTS_CONNECT = 100
	RECONNECT_SEC        = 5
	QUEUE_LENGTH         = 500
)

const (
	EXIT_DECLARE_EXCHANGE = 110
	EXIT_DECLARE_QUEUE    = 111
	EXIT_BIND             = 112
)

type connection struct {
	connectionAttempts int
	connection         *amqp.Connection
	channel            *amqp.Channel
	queueName          string
	stopping           bool

	dataSource string
	cdrEvent   chan *model.Cdr
}

func NewConnection(dataSource string) mq.MQ {
	return &connection{
		dataSource: dataSource,
		cdrEvent:   make(chan *model.Cdr, QUEUE_LENGTH),
		queueName:  fmt.Sprintf("cdr_metrics.%s", model.NewId()),
	}
}

func (c *connection) Start() error {
	c.initConnection()
	return nil
}

func (c *connection) Close() {
	wlog.Debug("AMQP receive stop client")
	c.stopping = true
	if c.channel != nil {
		c.channel.Close()
		wlog.Debug("close AMQP channel")
	}

	if c.connection != nil {
		c.connection.Close()
		wlog.Debug("close AMQP connection")
	}
}

func (c *connection) ConsumeCdr() <-chan *model.Cdr {
	return c.cdrEvent
}

func (a *connection) initConnection() {
	var err error

	if a.connectionAttempts >= MAX_ATTEMPTS_CONNECT {
		wlog.Critical(fmt.Sprintf("Failed to open AMQP connection..."))
		time.Sleep(time.Second)
		os.Exit(1)
	}
	a.connectionAttempts++
	a.connection, err = amqp.Dial(a.dataSource)
	if err != nil {
		wlog.Critical(fmt.Sprintf("Failed to open AMQP connection %s to err:%v", a.dataSource, err.Error()))
		time.Sleep(time.Second * RECONNECT_SEC)
		a.initConnection()
	} else {
		a.connectionAttempts = 0
		a.channel, err = a.connection.Channel()
		if err != nil {
			wlog.Critical(fmt.Sprintf("Failed to open AMQP channel to err:%v", err.Error()))
			time.Sleep(time.Second)
			os.Exit(1)
		} else {
			a.initQueues()
		}
	}
}

func (c *connection) initQueues() {
	queue, err := c.channel.QueueDeclare(c.queueName,
		false, true, true, true, nil)
	if err != nil {
		wlog.Critical(fmt.Sprintf("Failed to declare AMQP queue %v to err:%v", c.queueName, err.Error()))
		time.Sleep(time.Second)
		os.Exit(EXIT_DECLARE_QUEUE)
	}

	wlog.Debug(fmt.Sprintf("Success declare queue %v, connected consumers %v", queue.Name, queue.Consumers))
	c.subscribe()
}

func (c *connection) subscribe() {
	err := c.channel.QueueBind(c.queueName, "leg.#", model.CDR_EXCHANGE, false, nil)
	if err != nil {
		wlog.Critical(fmt.Sprintf("Error binding queue %s to %s: %s", c.queueName, model.CDR_EXCHANGE, err.Error()))
		time.Sleep(time.Second)
		os.Exit(EXIT_BIND)
	}

	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		wlog.Critical(fmt.Sprintf("Error create consume for queue %s: %s", c.queueName, err.Error()))
		time.Sleep(time.Second)
		os.Exit(EXIT_BIND)
	}

	go func() {
		defer wlog.Debug(fmt.Sprintf("stop listening queue"))
		for m := range msgs {
			c.cdrEvent <- model.CdrFromJson(m.Body)
		}

		if !c.stopping {
			c.initConnection()
		}
	}()
}
