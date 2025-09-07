package message

import (
	"user-service/config"

	"github.com/labstack/gommon/log"
)

func PublishMessage(email, message, notif_type string) error {
	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		log.Errorf("[PublishMessage-1] Failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[PublishMessage-2] Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		notif_type,
		true,
		false,
		false,
		false,
		nil,
	)
	
	if err != nil {
		log.Errorf("[PublishMessage-3] Failed to declare a queue: %v", err)
		return err
	}

}
