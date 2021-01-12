import express from "express";
import http from "http";
import socketIo, { Server } from "socket.io";
import rabbit from "amqp-connection-manager";

const port = process.env.PORT || 4001;
import { IndexController } from "./routes/index";
import { Socket } from "socket.io";
import { RabbitMqBus } from "./infrastructure/rabbitmq";
import { CrawledMedia } from "./domain/model";

const router = express.Router();
const app = express();
app.use(express.json())

const bus = RabbitMqBus.from(process.env.RABBITMQ_CONNECTION ?? "amqp://guest:guest@localhost:5672/")


app.use(IndexController.from(router, bus).routes());

const server = http.createServer(app);

const io: Socket = (socketIo as any)(server); // < Interesting!

bus.consume({ exchange: "crawled-media", queue: "crawled-media-app"}, (model: CrawledMedia) => {
  console.log(model.url);
  io.emit("new-crawled-media", model)
});

const getApiAndEmit = (socket: Socket) => {
    const response = new Date();
    // Emitting a new message. Will be consumed by the client
    socket.emit("FromAPI", response);
  };


let interval: NodeJS.Timeout;
io.emit("")
io.on("connection", (socket: Socket) => {
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