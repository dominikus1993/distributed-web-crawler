import dapr from 'dapr-client';
import util from 'util';
var messages = dapr.dapr_pb; 
var services = dapr.dapr_grpc;
var grpc = require('grpc');

const PORT = process.env.DAPR_GRPC_PORT || 50001;

export interface IMessage<T> {
    readonly message: T
}

export class DaprClient {
    #dapr
    constructor() {
        this.#dapr = new services.DaprClient(`localhost:${PORT}`, grpc.credentials.createInsecure())
    }

    async publish<T>(msg: IMessage<T>) {
        const event = new messages.PublishEventRequest();
        event.setTopic("crawl-website")
        event.setPubsubName("pubsub")
        event.setDataContentType('text/plain');
        console.log(msg.message)
        event.setData(Buffer.from(JSON.stringify(msg.message)))
        this.#dapr.publishEvent(event, (err, res) => {
            if(err) {
                console.error(err)
                return;
            }
            console.debug(res)
        })
    }
}