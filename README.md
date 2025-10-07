# SamsDataStore TCP Server ⚡️

## Overview
SamsDataStore Clone MVP is a high-performance, in-memory key-value store built in Go, designed to emulate a subset of Redis functionalities. This project demonstrates robust TCP networking, efficient RESP (REdis Serialization Protocol) parsing, and concurrent data handling with Time-To-Live (TTL) support, making it an excellent foundation for distributed caching or real-time data storage solutions. 🚀

## Features
- **Concurrent TCP Server**: Manages multiple client connections simultaneously using goroutines for high throughput.
- **RESP Protocol Implementation**: Full support for parsing and encoding Redis Serialization Protocol (RESP) messages, including arrays, bulk strings, and simple strings.
- **In-Memory Data Storage**: Utilizes Go's built-in `map` for rapid key-value data access.
- **Time-To-Live (TTL) Support**: Keys can be set with an expiration time, ensuring automatic removal after a specified duration.
- **Thread-Safe Operations**: Employs `sync.RWMutex` to protect data store operations, guaranteeing data consistency in a concurrent environment.
- **Core Redis Command Support**: Implements fundamental commands such as `PING`, `ECHO`, `SET`, `GET`, `DEL`, and `QUIT`.

## Getting Started
### Installation
To get a local copy of GoRedis Clone MVP running on your machine, follow these steps.

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/samueltuoyo15/Redis-Clone-Mvp.git
    ```
2.  **Navigate to the Project Directory**:
    ```bash
    cd Redis-Clone-Mvp
    ```
3.  **Run the Server**:
    You can run the server directly or build an executable.
    
    ✨ **Direct Run (for development)**:
    ```bash
    go run main.go
    ```
    
    📦 **Build and Run (for production/deployment)**:
    ```bash
    go build -o goredis-server .
    ./goredis-server
    ```
    The server will start listening on `localhost:6378`.

4.  **Using `air` for Live Reload (Optional, for development)**:
    This project includes an `.air.toml` configuration for `air`, a live-reloading tool for Go applications.
    
    First, install `air`:
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
    Then, run `air` from the project root:
    ```bash
    air
    ```
    `air` will automatically restart the server on code changes.

### Environment Variables
No environment variables are currently required to run the GoRedis Clone MVP server, as the listening port is hardcoded to `6378`.

## Usage
Once the server is running, you can interact with it using a simple TCP client like `netcat` or a Redis client that supports the RESP protocol.

**Example using `netcat`**:

1.  **Connect to the server**:
    ```bash
    nc localhost 6378
    ```
2.  **Send commands**:
    Type Redis commands in RESP array format or inline command format (which the server also supports).
    
    **Ping**:
    ```
    PING
    ```
    (Server response: `+PONG`)
    
    **Set a key**:
    ```
    SET mykey myvalue
    ```
    (Server response: `+OK`)
    
    **Get a key**:
    ```
    GET mykey
    ```
    (Server response: `$7\r\nmyvalue`)
    
    **Set a key with a 10-second expiry**:
    ```
    SET expirekey tempvalue EX 10
    ```
    (Server response: `+OK`)
    
    **Delete a key**:
    ```
    DEL mykey
    ```
    (Server response: `:1`)
    
    **Quit the connection**:
    ```
    QUIT
    ```
    (Server response: `+OK`, and the connection closes)

## API Documentation
GoRedis Clone MVP implements a custom TCP server that communicates using the RESP protocol, mimicking Redis.

### Base URL
The server listens on `TCP localhost:6378`.

### Endpoints (RESP Commands)
#### PING
Checks if the server is alive. Can optionally take a message to echo back.

**Request**:
```
*1\r\n$4\r\nPING\r\n
```
(Inline: `PING`)

**Request with argument**:
```
*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n
```
(Inline: `PING hello`)

**Response**:
```
+PONG\r\n
```
(If no argument provided)

**Response with argument**:
```
$5\r\nhello\r\n
```
(If an argument is provided)

**Errors**:
- None specific. If multiple arguments are provided, only the first argument after `PING` is echoed.

#### ECHO `<message>`
Echoes the provided message back to the client.

**Request**:
```
*2\r\n$4\r\nECHO\r\n$11\r\nHello World\r\n
```
(Inline: `ECHO Hello World`)

**Response**:
```
$11\r\nHello World\r\n
```

**Errors**:
- `-ERR wrong number of arguments for 'ECHO' command\r\n`: If no message argument is provided.

#### SET `<key>` `<value>` [EX `<seconds>`]
Sets the string value of a key. If `EX` is provided, the key will expire after the specified number of seconds.

**Request (basic)**:
```
*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
```
(Inline: `SET foo bar`)

**Request (with expiry)**:
```
*5\r\n$3\r\nSET\r\n$4\r\nmykey\r\n$7\r\nmyvalue\r\n$2\r\nEX\r\n$2\r\n60\r\n
```
(Inline: `SET mykey myvalue EX 60`)

**Response**:
```
+OK\r\n
```

**Errors**:
- `-ERR wrong number of arguments for 'SET' command\r\n`: If key or value is missing.
- Note: If an invalid value is provided for `EX` (e.g., non-numeric), the command will proceed, but the expiry will be ignored, and the key will be set without a TTL. No explicit error is returned for an invalid `EX` argument.

#### GET `<key>`
Retrieves the string value associated with a key.

**Request**:
```
*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n
```
(Inline: `GET foo`)

**Response (key exists)**:
```
$3\r\nbar\r\n
```

**Response (key does not exist or expired)**:
```
$-1\r\n
```

**Errors**:
- `-ERR wrong number of arguments for 'GET' command\r\n`: If no key argument is provided.

#### DEL `<key>` [`<key>` ...]
Deletes one or more keys.

**Request (single key)**:
```
*2\r\n$3\r\nDEL\r\n$3\r\nfoo\r\n
```
(Inline: `DEL foo`)

**Request (multiple keys)**:
```
*3\r\n$3\r\nDEL\r\n$4\r\nkey1\r\n$4\r\nkey2\r\n
```
(Inline: `DEL key1 key2`)

**Response**:
```
:1\r\n
```
(Integer representing the number of keys successfully deleted)

**Errors**:
- `-ERR wrong number of arguments for 'DEL' command\r\n`: If no key arguments are provided.

#### QUIT
Closes the current client connection.

**Request**:
```
*1\r\n$4\r\nQUIT\r\n
```
(Inline: `QUIT`)

**Response**:
```
+OK\r\n
```
(The server sends `OK` and then gracefully closes the TCP connection.)

**Errors**:
- None specific.

## Technologies Used
| Technology         | Description                                     | Link                                                        |
| :----------------- | :---------------------------------------------- | :---------------------------------------------------------- |
| **Go**             | Core language for server implementation.        | [https://golang.org/](https://golang.org/)                  |
| **`net` package**  | TCP networking primitives for server and client. | [https://pkg.go.dev/net](https://pkg.go.dev/net)          |
| **`bufio` package**| Buffered I/O operations for efficient data reading. | [https://pkg.go.dev/bufio](https://pkg.go.dev/bufio)      |
| **`sync` package** | Primitives for concurrent programming (e.g., `RWMutex`). | [https://pkg.go.dev/sync](https://pkg.go.dev/sync)        |
| **`time` package** | Handling time-based operations for TTL.         | [https://pkg.go.dev/time](https://pkg.go.dev/time)        |
| **`strconv` package** | String conversions for RESP parsing.        | [https://pkg.go.dev/strconv](https://pkg.go.dev/strconv)  |
| **Air**            | Live-reloading tool for Go applications (development). | [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air) |

## Contributing
Contributions are welcome! If you have suggestions for improvements, new features, or bug fixes, please follow these guidelines:

*   Fork the repository.
*   Create a new branch for your feature or bug fix: `git checkout -b feature/your-feature-name` or `bugfix/issue-description`.
*   Make your changes and ensure they adhere to Go best practices.
*   Write clear, concise commit messages. 📝
*   Push your branch and submit a pull request. ⬆️
*   Describe your changes in detail in the pull request. ✍️

## Author
**Samuel Tuoyo**

*   **LinkedIn**: [https://www.linkedin.com/in/samuel-tuoyo-8568b62b6/]
*   **Twitter**: [https://x.com/TuoyoS26091]
---

[![Go Version](https://img.shields.io/badge/Go-1.25.1-00ADD8?logo=go)](https://golang.org/)
[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)
