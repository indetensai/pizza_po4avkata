[Русская версия](https://github.com/indetensai/pizza_po4avkata/blob/main/README_ru.md)

# PIZZA_PO4AVKATA🍕
**PIZZA_PO4AVKATA** is a golang backend project that I created with the goal of learning how to use gRPC and design microservices. The project simulates a pizza delivery service where users can browse menus, place orders, and track their deliveries. It uses Fiber, a fast and lightweight web framework, to handle HTTP requests and responses. The project uses MongoDB Atlas as the database solution, which offers a scalable and secure cloud service for MongoDB. It uses sessions for secure authentication and authorization, to safeguard the endpoints and authenticate the users. It uses gRPC, a high-performance RPC framework, to communicate between microservices. It uses Docker and docker-compose for easy deployment and portability. The project enables users to order pizzas online. It comprises of three microservices: user service, menu service, and order service. It also has an API gateway that accepts requests and routes them to the suitable service. **PIZZA_PO4AVKATA** is a fun and educational project that demonstrates the benefits of gRPC and microservices in building modern and reliable backend applications. It is also a great way to show off your golang skills and impress potential employers with your passion for pizza. 

## How to run
Before you continue, ensure you have met the following requirements:
- You have installed the latest version of Go.
- You have installed MongoDB and created a database for the project(alternative choice is to use MongoDB Atlas).
- You have created .env files, with mongodb database url defined.
- You have installed grpc.

To run **PIZZA_PO4AVKATA**, follow these steps:
1. Clone this repository: `git clone https://github.com/indetensai/pizza_po4avkata.git`
2. Change into the project directory: `cd pizza_po4avkata`
3. Install the dependencies: `go mod download`
4. Fill `.env` files with the required environment variables.
5. Build the executable: `go build -o pizza_po4avkata cmd/pizza_po4avkata/main.go`

## How to run (docker-compose)
Before you continue, ensure you have met the following requirements:
- You have installed the latest version of docker(-desktop).

To run **PIZZA_PO4AVKATA** using docker-compose, follow these steps:
1. Clone this repository: `git clone https://github.com/indetensai/pizza_po4avkata.git`
2. Change into the project directory: `cd pizza_po4avkata`
3. Run `docker compose up`

## Usage
To run **PIZZA_PO4AVKATA**, follow these steps:
1. Start the executable: `./pizza_po4avkata`
2. To interact with the API, you can use any HTTP client of your choice.

The **PIZZA_PO4AVKATA** API has the following endpoints:
- `POST /user/register`: Create a new user account.
- `POST /user/login`: Login with an existing user account and get a session.
- `GET /menu`: Get a menu.
- `POST /order`: Make an order with specified dishes. Requires authentication.
- `GET /admin/orders`: Get a list of orders and their status. Requires admin privileges.
- `POST /admin/pizza`: Create a dish in menu. Requires admin privileges.