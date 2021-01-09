//@ts-check

const express = require("express");
const http = require("http");
const socketIo = require("socket.io");
const rabbit = require("amqp-connection-manager")

const port = process.env.PORT || 4001;
const Index = require("./routes/index");
const router = express.Router();
const app = express();
app.use(express.json())
const connection = rabbit.connect([process.env.RABBITMQ_CONNECTION ?? "amqp://guest:guest@localhost:5672/"]);
const channel = connection.createChannel();
app.use(Index.from(router, channel).routes());

const server = http.createServer(app);

// @ts-ignore
const io = socketIo(server); // < Interesting!

/**
 * @param {{ emit: (arg0: string, arg1: Date) => void; }} socket
 */
const getApiAndEmit = socket => {
    const response = new Date();
    // Emitting a new message. Will be consumed by the client
    socket.emit("FromAPI", response);
  };

let interval;
/**
 * @param {{ on?: any; emit?: (arg0: string, arg1: Date) => void; }} socket
 */
io.on("connection", (socket) => {
    console.log("New client connected");
    if (interval) {
      clearInterval(interval);
    }
    interval = setInterval(() => getApiAndEmit(socket), 1000);
    socket.on("disconnect", () => {
      console.log("Client disconnected");
      clearInterval(interval);
    });
  });

server.listen(port, () => console.log(`Listening on port ${port}`));