import { ChannelWrapper } from "amqp-connection-manager";
import { Router } from "express";
import { publishToRabbitMq } from "../infrastructure/rabbitmq";

export class IndexController {
    #channel : ChannelWrapper
    #router : Router
    constructor(router: Router, channel: ChannelWrapper) {
        this.#channel = channel;
        this.#router = router;
    }

    routes() {
        this.#router.post("/",async (req, res) => {
            const { url } : { url: string } = req.body;

            await publishToRabbitMq(this.#channel, { exchange: "crawl-media", message: { url } })

            res.send({ status: "ok" }).status(204);
        })
        return this.#router;
    }

    static from(router: Router, channel: ChannelWrapper)  {
        return new IndexController(router, channel)
    }
}
