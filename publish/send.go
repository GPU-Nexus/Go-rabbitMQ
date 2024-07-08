package main

import (
	//"context"
	"go_rabbitMQ/messages"
	"log"
	"os"
	"bufio"

	//"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	//send queues
	GO_HELLO := "GO_hello"
	MATRIX_QUEUE := "matrix_queue"
	CUDA_QUEUE := "cuda_queue"
	GO_HELLO_REPLY := "GO_hello_reply"

	EXCHANGE_NAME := "GO_exchange"

	//reply queues
	// GO_HELLO_REPLY = "GO_hello_reply"
	// MATRIX_QUEUE_REPLY = "matrix_queue_reply"
	// CUDA_QUEUE_REPLY = "cuda_queue_reply"

	//routing keys
	GO_ROUTING_KEY := "GO_routing_key"
	MATRIX_ROUTING_KEY := "matrix_multiplication"
	CUDA_ROUTING_KEY := "cuda_program"

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

	//----------------------------------------------------------EXCHANGE DECLARATION--------------------------------------------

	// Declare an exchange

	err = ch.ExchangeDeclare(
		EXCHANGE_NAME, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	//----------------------------------------------------------QUEUE DECLARATION--------------------------------------------

	//queue for hello world
	q, err := ch.QueueDeclare(
		GO_HELLO, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//queue for matrix data
	q2, err := ch.QueueDeclare(
		MATRIX_QUEUE, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare the matrix_queue")

	//queue for cuda program
	q3, err := ch.QueueDeclare(
		CUDA_QUEUE, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare the cuda_queue")

	//----------------------------------------------------------QUEUE BINDING--------------------------------------------

	// message binding
	err = ch.QueueBind(
		q.Name,         // queue name
		GO_ROUTING_KEY, // routing key
		EXCHANGE_NAME,  // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	//matrix binding
	err = ch.QueueBind(
		q2.Name,            // queue name
		MATRIX_ROUTING_KEY, // routing key
		EXCHANGE_NAME,      // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	//cuda binding
	err = ch.QueueBind(
		q3.Name,          // queue name
		CUDA_ROUTING_KEY, // routing key
		EXCHANGE_NAME,    // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	//----------------------------------------------------------------------------------------------SENDING_MESSAGE------------------------------------------------------------

	content := bufio.NewReader(os.Stdin)
	id := bufio.NewReader(os.Stdin)
	log.Printf("Enter the Id: ")
	inputId,_ := id.ReadString('\n')
	log.Printf("Enter the message: ")
	inputContent, _ := content.ReadString('\n')
	


	//Create a new message
	msg := &messages.MyMessage{
		Id:      inputId,
		Content: inputContent,
	}

	//Marshal the message to Protobuf
	body1, err := proto.Marshal(msg)
	failOnError(err, "Failed to marshal message")

	err = ch.Publish(
		EXCHANGE_NAME,  // exchange
		GO_ROUTING_KEY, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        body1,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body1)

	//sendMatrixData(exchangeName, routingKey_matrix, ch)
	sendMatrixData(EXCHANGE_NAME, MATRIX_ROUTING_KEY, ch)

	//sendCudaProgram(ch, exchangeName)
	sendCuda(EXCHANGE_NAME, CUDA_ROUTING_KEY, ch)

	//getting reply for hello world
	gettingReplyForHelloWorld(ch, GO_HELLO_REPLY)
}

func sendMatrixData(exchangeName string, routingKey string, ch *amqp.Channel) {
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
	err = ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        append(matABytes, matBBytes...), // Combine both matrices
		})
	failOnError(err, "Failed to publish matrices to RabbitMQ")
	log.Printf(" [x] Sent matrices")
}

func sendCuda(exchangeName string, routingKey string, ch *amqp.Channel) {


	// Create a new message
	msg := &messages.CudaProgram{
		Code: "cuda code",
	}

	// Marshal the message to Protobuf
	body, err := proto.Marshal(msg)
	failOnError(err, "Failed to marshal message")

	// Publish a message to the exchange with the custom routing key
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
	log.Printf(" [x] Sent %s\n", body)
}


func gettingReplyForHelloWorld(ch *amqp.Channel, hello_reply_queue string){

	q, err := ch.QueueDeclare(
		hello_reply_queue, // name
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
			//log.Printf("Received a message: %s", d.Body)
			// Here you can unmarshal the message from protobuf if needed
			var msg messages.MyMessageReply
			err := proto.Unmarshal(d.Body, &msg)
			failOnError(err, "Failed to unmarshal message")

			// Print the message
			log.Printf("Received message: %s", msg.Reply)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {} // Keep the program running
}