import { Router } from "express";
import {IMessage} from "../infrastructure/dapr"

export default function index(publish: (msg: IMessage<any>) => Promise<any>) : (router: Router) => Router {
    return (router: Router) => {
        router.post("/",async (req, res) => {
            const { url } : { url: string } = req.body;

            await publish({ message: { url } })

            res.send({ status: "ok" }).status(204);
            res.end()
        })

        return router;
    }
}
