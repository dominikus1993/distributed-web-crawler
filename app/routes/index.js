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
        this.#router.get("/",async (_, res) => {
            await publishToRabbitMq(this.#channel, { exchange: "crawl-media", message: JSON.stringify({ url: "https://jbzd.com.pl/" }) })
            res.send({ response: "I am alive" }).status(200);
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