import express from "express";
import http from "http";
import socketIo, { Server } from "socket.io";
const port = process.env.PORT || 4001;
import index from "./routes/index";
import { Socket } from "socket.io";
import status from "./routes/status";
import cors from "cors"
import dapr from "./routes/dapr";
import { publish } from "./infrastructure/dapr";

const router = express.Router();
const app = express();
app.use(cors({origin: "*" }))
app.use(express.json({ type: 'application/*+json' }))
app.use(express.json({ type: 'application/json' }))
let io: Socket | undefined = undefined;
app.use(index(publish)(router));
app.use(status(router));
app.use(dapr(io)(router));
const server = http.createServer(app);

io = (socketIo as any)(server, { cors: { orgin: "*" }}); // < Interesting!

if(io === undefined) {
  throw new Error("Can't start socket")
}

const getApiAndEmit = (socket: Socket) => {
    const response = new Date();
    // Emitting a new message. Will be consumed by the client
    socket.emit("FromAPI", response);
  };


let interval: NodeJS.Timeout;
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