import { CrawledMedia } from "../domain/model";
import { Router } from "express";
import { Socket } from "socket.io";

export default function dapr(socket: Socket | undefined): (router: Router) => Router {
    return (router: Router) => {

        router.get("/dapr/subscribe", (req, res) => {
            const subs = [{ pubsubname: "pubsub", topic: "crawled-media", route: "crawled-media" }, { pubsubname: "pubsub", topic: "crawl-website", route: "test" }]
            res.send(subs).status(200)
        })
        router.post("/crawled-media", async (req, res) => {
            const model: CrawledMedia = req.body;
            console.debug("crawled", model)
            if (socket !== undefined) {
                socket.emit("new-crawled-media", model)
            }

            res.status(200).send({
                "status": "SUCCESS"
            })
        })

        router.post("/test", async (req, res) => {
            const model = req.body;
            console.debug("test", model)
            console.debug(req.headers["content-type"])
            res.status(200).send({
                "status": "SUCCESS"
            })
            res.end()
        })
        return router;
    }
}