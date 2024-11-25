# RealCode

RealCode is a code compilation service that supports both Golang and Python. It is built using a microservices architecture to ensure scalability and maintainability. The project leverages Golang for service implementation and Apache Kafka as a messaging queue to handle communication between services

## Features

- **Golang and Python Compilation**: Supports compiling and executing code written in Golang and Python
- **Microservices Architecture**: Modular design with separate services for different functionalities
- **Apache Kafka**: Used as a messaging queue to manage communication between the server and compilation services
- **Scalable**: Easily add more services to handle increased load without significant changes to the existing infrastructure

## Architecture

The architecture of RealCode consists of the following components:

1. **API Gateway**: The entry point for users to submit their code.
2. **Kafka Message Queue**: Handles the distribution of code compilation requests to the respective services.
3. **Golang Compilation Service**: Compiles and executes Golang code.
4. **Python Compilation Service**: Compiles and executes Python code.
5. **Result Aggregator**: Collects and sends the results back to the user.

## Getting Started

### Prerequisites

- Golang (latest version)
- Python (latest version)
- Apache Kafka
- Docker (for containerized deployment)

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a new Pull Request.


