import http from "http";
export interface IMessage<T> {
    readonly message: T
}

const daprPort = process.env.DAPR_HTTP_PORT || 3504

export function publish<T>(msg: IMessage<T>) {

    return new Promise((resolve, reject) => {
        const message = {
            "data": {
                "message": msg.message,
            }
        };
        console.debug("Message = ", message)
        const data = JSON.stringify(message)
        const options: http.RequestOptions = {
            hostname: "localhost",
            port: daprPort,
            method: "POST",
            path: "/v1.0/publish/pubsub/crawl-website",
            headers: { 'Content-Type': 'application/json', 'Content-Length': data.length },
        } 
        const req = http.request(options, res => {
            console.log(`statusCode: ${res.statusCode}`)

            res.on('data', d => {
                console.debug("data", d);
            })

            resolve(res)
        })

        req.on('error', error => {
            console.error(error)
            reject(error)
        })
        
        req.write(data);
        req.end();
    });
}

