import express from "express";
const router = express.Router();

router.get("/", (_, res) => {
    res.send({ response: "I am alive" }).status(200);
});

export default router;