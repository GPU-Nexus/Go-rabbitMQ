package main

import (
	"go_rabbitMQ/messages"
	"log"
	"os"
	"strings"
	"google.golang.org/protobuf/proto"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	//consuming  queues
	GO_HELLO := "GO_hello"
	//MATRIX_QUEUE := "matrix_queue"
	//CUDA_QUEUE := "cuda_queue"

	EXCHANGE_NAME := "GO_exchange"

	//reply queues
	GO_HELLO_REPLY := "GO_hello_reply"
	// MATRIX_QUEUE_REPLY = "matrix_queue_reply"
	// CUDA_QUEUE_REPLY = "cuda_queue_reply"

	//routing keys
	GO_REPLY_ROUTING_KEY := "GO_Reply_routing_key"

	// Load environment variables from .env file
	err := godotenv.Load()
	failOnError(err, "Error loading .env file")

	// Fetch RabbitMQ password from environment variables
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	if rabbitMQPassword == "" {
		log.Panic("RABBITMQ_PASSWORD not found in environment variables")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:" + rabbitMQPassword + "@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		GO_HELLO, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			// Here you can unmarshal the message from protobuf if needed
			// For example, if you are expecting a protobuf message
			var msg messages.MyMessage
			err := proto.Unmarshal(d.Body, &msg)
			failOnError(err, "Failed to unmarshal message")

			// Print the message
			log.Printf("Received message: %s", msg.Content)
			log.Printf("Received message ID: %s", msg.Id)
			
			msgIdTrimmed := strings.TrimSpace(msg.Id)
			//sending reply message
			if(msgIdTrimmed == "123"){
				log.Printf("Sending reply message")
				sendingHeloWorldReply(ch, EXCHANGE_NAME, GO_HELLO_REPLY, GO_REPLY_ROUTING_KEY, msg.Content)
				
			}else {
				log.Printf("Message ID did not match: %s", msgIdTrimmed)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {} // Keep the program running

}

func sendingHeloWorldReply(ch *amqp.Channel, exchangeName string, replyQName string, routingKey string, messageReceived string) {

	// Declare an exchange
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	//queue for hello world reply message
	qReply, err := ch.QueueDeclare(
		replyQName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		qReply.Name,  // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)

	failOnError(err, "Failed to bind a queue")

	//----------------------------------------------------------SENDING REPLY MESSAGE--------------------------------------------

	// create reply message
	msg := &messages.MyMessageReply{
		Id:      "124",
		Reply: "Hi I am the consumer. this is the message you sent: " + messageReceived + " thanks for sending message.",
	}

	// Marshal the message into protobuf
	body, err := proto.Marshal(msg)
	failOnError(err, "Failed to marshal message")

	err = ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        body,
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", msg.Reply)

	
}
