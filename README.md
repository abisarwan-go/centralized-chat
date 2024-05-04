# Simple Chat Application in Go

This repository contains the Go code for a simple chat application that supports basic user authentication and allows users to chat with each other in real-time over TCP.

## Overview

The application can run in two modes:
- **Server Mode:** Listens for incoming TCP connections and handles user authentication and message forwarding between users.
- **Client Mode:** Connects to the server, authenticates, and allows the user to send and receive messages.

## Features

- User authentication
- Real-time messaging
- List online users
- Clear terminal screen for better user experience

## Installation

To run this application, you need Go installed on your machine. You can download it from [Go's official site](https://golang.org/dl/).

After installing Go, clone this repository to your local machine:

```bash
git clone https://github.com/abisarwan-go/centralized-chat
cd simple-chat-app-go
```

## Usage

To start the server, run the following command in your terminal:

```bash
go run main.go
```

To connect as a client, run the same command in a different terminal window. The application will automatically detect that the server is already running and will initiate in client mode.

### As a Client

1. Enter your unique ID for authentication.
2. Choose from the available options to chat with other users or exit.

### As a Server

Just start the server. It will handle incoming connections and facilitate the authentication and communication between connected clients.

## Contributing

Contributions are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)