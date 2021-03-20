import { Router } from "express";

export class StatusController {
    #router : Router
    constructor(router: Router) {
        this.#router = router;
    }

    routes() {
        this.#router.get("/ping",async (_, res) => {
            res.send({ message: "pong" }).status(200);
        })
        return this.#router;
    }

    static from(router: Router)  {
        return new StatusController(router)
    }
}
