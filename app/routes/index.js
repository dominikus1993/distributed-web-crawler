const { Router } = require("express");
const { publishToRabbitMq } = require("../infrastructure/rabbitmq")
class Routing {
    /** @type {import("amqp-connection-manager").ChannelWrapper} */
    #channel
    /** @type {Router} */
    #router
    constructor(router, channel) {
        this.#channel = channel;
        this.#router = router;
    }

    routes() {
        this.#router.post("/",async (req, res) => {
            const { url } = req.body;
            await publishToRabbitMq(this.#channel, { exchange: "crawl-media", message: JSON.stringify({ url }) })
            res.send({ status: "ok" }).status(204);
        })
        return this.#router;
    }

    /**
     * 
     * @param {Router} router 
     * @param {import("amqp-connection-manager").ChannelWrapper} channel 
     */
    static from(router, channel)  {
        return new Routing(router, channel)
    }
}


module.exports = Routing;