# WebSocket Chat Application in Go

This project is a WebSocket-based chat application developed in Go, providing real-time communication between users. It's designed to showcase the use of WebSockets in a Go server environment, paired with a React frontend for a seamless chatting experience.

## Developers

- Maxim Turbulyak
- Kuanysh Kambarov

Group: SE-2216

## Features

- Real-time messaging between users.
- Support for multiple chat rooms.
- User authentication and chat history.
- Scalable architecture for handling concurrent connections.

## Technologies

- **Backend**: Go (Gorilla WebSocket package for WebSocket support)
- **Frontend**: React.js

## Getting Started

To get the project up and running on your local machine, follow these steps:

1. **Clone the repository**

```bash
git clone https://github.com/yourgithubusername/websocket-chat-app.git
cd websocket-chat-app
```

2. **Start the backend server**

Navigate to the backend directory and run:

```bash
go run main.go
```

3. **Start the frontend application**

Navigate to the frontend directory and install the dependencies:

```bash
npm install
```

Then, start the React application:

```bash
npm start
```

The application will be available at `http://localhost:3000`.

## Documentation

For more detailed information about the project structure and API endpoints, refer to the `docs` directory.

## Contributing

We welcome contributions! Please feel free to fork the repository and submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
