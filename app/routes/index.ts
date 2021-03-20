import { Router } from "express";
import dapr from 'dapr-client';
import {DaprClient} from "../infrastructure/dapr"
const services = dapr.dapr_grpc;

export class IndexController {
    #bus : DaprClient
    #router : Router
    constructor(router: Router, client: DaprClient) {
        this.#bus = client;
        this.#router = router;
    }

    routes() {
        this.#router.post("/",async (req, res) => {
            const { url } : { url: string } = req.body;

            await this.#bus.publish({ message: { url } })

            res.send({ status: "ok" }).status(204);
        })
        return this.#router;
    }

    static from(router: Router, bus: DaprClient)  {
        return new IndexController(router, bus)
    }
}
