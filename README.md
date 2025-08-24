# go-tcp-chat

A simple TCP-based concurrent chat server written in Go. Clients can connect via the terminal using `nc` (netcat), choose a nickname, and send public or private messages to other connected users.

## Features

- Handles multiple concurrent clients using Goroutines
- Prompts each user to choose a unique nickname on connect
- Broadcasts messages to all users
- Supports private messaging using `@nickname message` syntax
- Synchronized access to client data using `sync.Mutex`
- Standard library only (no external dependencies)

## Requirements

- Go 1.18 or higher
- `nc` (netcat) installed for connecting as a client

## Getting Started

### Clone the Repository

```bash 
git clone https://github.com/sanki92/go-tcp-chat.git
cd go-tcp-chat
go run main.go
```
#### Connect as a Client
```bash
nc localhost 8080
```
