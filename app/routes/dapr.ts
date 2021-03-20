import { CrawledMedia } from "../domain/model";
import { Router } from "express";
import { Socket } from "socket.io";

export class DaprSubcriptionController {
    #socket : Socket | undefined
    #router : Router
    constructor(router: Router, io: Socket | undefined) {
        this.#socket = io;
        this.#router = router;
    }

    routes() {
        this.#router.get("/dapr/subscribe", (req, res) => {
            const subs = [{pubsubname: "pubsub", topic: "crawled-media", route: "crawled-media"}]
            res.send(subs).status(200)
        })
        this.#router.post("/crawled-media",async (req, res) => {
            const model : CrawledMedia = req.body;
            console.debug("crawled", model)
            if(this.#socket !== undefined) {
                this.#socket.emit("new-crawled-media", model)
            }
        
            res.status(200);
        })
        return this.#router;
    }

    static from(router: Router, io: Socket | undefined)  {
        return new DaprSubcriptionController(router, io)
    }
}
