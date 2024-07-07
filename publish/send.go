package main

import (
	//"context"
	"go_rabbitMQ/messages"
	"log"
	"os"

	//"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Define the RabbitMQ connection string with the custom username and password
	rabbirMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	conn, err := amqp.Dial("amqp://guest:" + rabbirMQPassword + "@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare an exchange
	exchangeName := "GO_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare a queue
	queueName := "GO_hello"
	queueName2 := "matrix_queue"

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//queue for matrix data
	q2, err := ch.QueueDeclare(
		queueName2, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare the matrix_queue")

	//Create a new message
    msg := &messages.MyMessage{
        Id:      "123",
        Content: "Hello, RabbitMQ!",
    }

	//Marshal the message to Protobuf
    body1, err := proto.Marshal(msg)
    failOnError(err, "Failed to marshal message")

	// Bind the queue to the exchange with a routing key
	routingKey := "GO_routing_key"
	routingKey_matrix := "matrix_multiplication"
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	//matrix binding
	err = ch.QueueBind(
		q2.Name,       // queue name
		routingKey_matrix,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	// Prepare matrix data (example)
    matA := &messages.Matrix{
        Data: []float32{1, 2, 3, 4, 5, 6},
        Rows: 2,
        Cols: 3,
    }

    matB := &messages.Matrix{
        Data: []float32{7, 8, 9, 10, 11, 12},
        Rows: 3,
        Cols: 2,
    }

	// Serialize matrix data to Protobuf
    matABytes, err := proto.Marshal(matA)
    failOnError(err, "Failed to marshal Matrix A")
    matBBytes, err := proto.Marshal(matB)
    failOnError(err, "Failed to marshal Matrix B")

	// Publish a message to the exchange with the custom routing key
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// Publish matrix data to RabbitMQ
    err = ch.Publish(
        exchangeName, // exchange
        routingKey_matrix,   // routing key
        false,        // mandatory
        false,        // immediate
        amqp.Publishing{
            ContentType: "application/protobuf",
            Body:        append(matABytes, matBBytes...), // Combine both matrices
        })
    failOnError(err, "Failed to publish matrices to RabbitMQ") 
	log.Printf(" [x] Sent matrices")

	
	err = ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        body1,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body1)
}
