rabbit-up:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=password heidiks/rabbitmq-delayed-message-exchange:latest

rabbit-down:
	docker stop rabbitmq
	docker rm rabbitmq