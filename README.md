# run the rabbitmq docker image or use the deployed rabbitmq in gcp, change the url accordingly

## first make the ssh connections by running the .bat file

## use to install dependencies
```
go mod tidy
```

## to compile the .proto file
```
protoc --go_out=. --go_opt=paths=source_relative messages/messages.proto
```

## go to the \publish and run to publish the hello world message
```
cd publish
go run send.go
``` 

## go to the \consume and run to consume that message
```
cd consume
go run consumer.go
``` 


