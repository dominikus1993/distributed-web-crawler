//@ts-check

import express from "express";
import http from "http";
import ioserver, { Socket } from 'socket.io';

const port = process.env.PORT || 4001;
import index from "./routes/index";

const app = express();
app.use(index);

const server = http.createServer(app);

const io = (ioserver as any)(server); // < Interesting!

/**
 * @param {{ emit: (arg0: string, arg1: Date) => void; }} socket
 */
const getApiAndEmit = (socket: Socket) => {
    const response = new Date();
    // Emitting a new message. Will be consumed by the client
    socket.emit("FromAPI", response);
  };

let interval: NodeJS.Timeout;
/**
 * @param {{ on?: any; emit?: (arg0: string, arg1: Date) => void; }} socket
 */
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