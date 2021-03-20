import { ChannelWrapper } from "amqp-connection-manager";
import { Router } from "express";
import { RabbitMqBus } from "../infrastructure/rabbitmq";

export class IndexController {
    #bus : RabbitMqBus
    #router : Router
    constructor(router: Router, bus: RabbitMqBus) {
        this.#bus = bus;
        this.#router = router;
    }

    routes() {
        this.#router.post("/",async (req, res) => {
            const { url } : { url: string } = req.body;

            await this.#bus.publish({ exchange: "crawl-media", message: { url } })

            res.send({ status: "ok" }).status(204);
        })
        return this.#router;
    }

    static from(router: Router, bus: RabbitMqBus)  {
        return new IndexController(router, bus)
    }
}
