# Go RPC Chat Application

A simple real-time chat system built using Go and net/rpc.
This project includes both a client and a server, allowing multiple users to join, send messages, and receive live updates using polling.

---

## Features

* RPC-based communication
* Supports multiple clients
* Real-time message updates using polling
* Broadcast system for sending messages to all users
* Simple and lightweight

---

## Project Structure

```
project/
│── client.go
│── server.go
│── README.md
```

---

## How to Run

### 1. Run the Server

Make sure to enable the server build tag:

```bash
go run -tags=server .
```

The server will start at:

```
localhost:1234
```

---

### 2. Run the Client

Open another terminal and run:

```bash
go run .
```

You should see:

```
Connected as User 0
You:
```

Open additional terminals to connect more users.

---

## Client Commands

* Typing a message sends it to all connected users
* Type `exit` to close the client

---

## How It Works

### Server:

* Tracks connected clients
* Sends broadcast messages to all clients except the sender
* Maintains a buffered channel for each client to store pending messages

### Client:

* Polls the server every 300ms for any incoming message
* Sends user input as RPC calls

---

## Example Output

```
User 1 joined
User 0: Hello
User 1: Hi
```

---

## Requirements

* Go 1.18 or newer
* Multiple terminals to simulate multiple users

---

## Contributing

Pull requests and improvements are welcome.

---

## License

MIT License
