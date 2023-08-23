# Golang HTTP Boilerplate Project
This is a comprehensive Golang HTTP boilerplate project that provides a foundation for building robust and scalable HTTP services. It incorporates various features such as hot reload, route versioning, dependency injection, support for MySQL cluster and MongoDB, RabbitMQ integration, a configurable HTTP server, configuration management, logging, error handling, and standardized request/response formats. Additionally, it includes middleware for CORS, panic recovery, and request logging.

## Features
- Hot Reload: Utilize tools like air to achieve hot reload during development, ensuring faster iteration and testing.

- Route Versioning: Implement route versioning to manage different API versions effectively.

- Dependency Injection (DI): Manage dependencies using DI patterns for flexible and maintainable code.

- Interface Usage: Utilize interfaces to decouple components and improve code extensibility.

- Database Support:

	- MySQL Cluster: Implement MySQL cluster support for handling large-scale data storage and retrieval.
	- MongoDB: Integrate MongoDB for NoSQL data storage and querying.
- RabbitMQ Integration: Integrate RabbitMQ for asynchronous messaging and event-driven architecture.

- Configurable HTTP Server: Customize the behavior of the HTTP server through configuration settings.

- Configuration Management: Manage environment-based JSON configuration in a clean and efficient manner.

- Logging: Use ZeroLog to implement logging to record application activities and troubleshoot issues effectively.

- Error Handling: Develop a robust error handling strategy to ensure graceful degradation and clear error messages.

- Request and Response Global Format: Define consistent formats for HTTP request and response payloads.

- Middleware:

	- CORS Middleware: Handle Cross-Origin Resource Sharing to ensure secure API communication.
	- Panic Recovery Middleware: Recover from unexpected panics and prevent application crashes.
	- Request Logging Middleware: Log incoming requests and outgoing responses for analysis and auditing.

## Requirement

Before you proceed, please make sure your machine has met the following requirements:

| Dependency |   Version   |
| ---------- | :---------: |
| Go       | >= v1.17 |


## Getting Started
- Clone this repository:
```sh
	git clone https://github.com/hi10drasingh/golang-base.git
	cd golang-base 
```
- Install dependencies:
```sh
  go mod tidy
```
- Configure your application settings in **./config/local.config.json**
- Run application locally with [air](https://github.com/cosmtrek/air)
```sh
  air
```
- Now you can add new versioned routes in **./internal/app/handlers**
