import { Router } from "express";


export default function status(router: Router): Router {

    router.get("/ping", async (_, res) => {
        res.status(200).send({ message: "pong" });
        res.end()
    })
    return router;

}
