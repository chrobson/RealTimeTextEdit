# Real-Time Collaborative Text Editor

A simple real-time collaborative text editor built with Go, WebSockets, and Redis. This application allows multiple users to edit the same text document simultaneously, with changes reflected in real-time across all connected clients.


---

## Features

- Real-time collaborative editing using WebSockets.
- Supports concurrent edits with basic conflict resolution.
- UTF-8 encoding support for international characters (including Polish characters).
- Simple and clean user interface.
- Written in Go with a minimal set of dependencies.
- Unit tests using Given-When-Then structure and Testify assertions.


### Running the Application

#### With Docker

**Build and Start Services**

   ```bash
   docker-compose up --build
   ```

   This command will build the Go application image and start both the application and Redis services.


## Usage

- Open `http://localhost:8080/` in multiple browser windows or tabs.
- Start typing in the text area. Your changes will be reflected in real-time across all connected clients.
- The application supports inserting and deleting text with basic conflict resolution.

## Details

- **Server**:
    - Built with Go.
    - Uses `net/http` for serving HTTP requests.
    - Uses `github.com/gorilla/websocket` for WebSocket connections.
    - Uses `github.com/go-redis/redis/v8` for Redis client.
    - Manages connected clients and broadcasts edits via WebSockets.

- **Client**:
    - Written in HTML and JavaScript. (not the best, nor the worst, just not doing much front-end)
    - Establishes a WebSocket connection to receive real-time updates.
    - Sends edits to the server via HTTP `POST` requests.
    - Uses a unique `clientId` to prevent applying its own edits multiple times.

- **Data Flow**:
    - **Edits**: When a user makes a change, the client calculates the difference and sends an edit to the server.
    - **Server Processing**: The server applies the edit to its text state and broadcasts the edit to all connected clients.
    - **Client Update**: Clients receive the edit via WebSocket and apply it to their local text area.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
