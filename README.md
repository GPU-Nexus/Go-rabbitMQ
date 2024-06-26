# run the rabbitmq docker image or use the deployed rabbitmq in gcp, change the url accordingly
## use to install dependencies
```
go mod tidy
```
## go to the \publish and run to publish the hello worl message
```
go run send.go
``` 

## go to the \consume and run to consume that message
```
go run consumer.go
``` 
