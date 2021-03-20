import rabbit, { ChannelWrapper, AmqpConnectionManager } from "amqp-connection-manager"
import { Channel, ConsumeMessage } from "amqplib";

export interface IMessage<T> {
    readonly exchange: string
    readonly message: T
    readonly topic?: string
}

export interface ISubscription {
    readonly exchange: string
    readonly queue: string
    readonly topic?: string
}

function onMessage<T>(action: (obj: T) => void) {
    return (ch: Channel) => {
        return (data: ConsumeMessage | null) => {
            if (data) {
                const msg: T | null | undefined = JSON.parse(data.content.toString())
                if (msg) {
                    action(msg)
                    ch.ack(data)
                }
            }
        }
    }
}

export class RabbitMqBus {
    #connection: AmqpConnectionManager;
    #channel: ChannelWrapper
    #subscriptions: ChannelWrapper[] = []

    private constructor(connection: AmqpConnectionManager, channel: ChannelWrapper) {
        this.#connection = connection;
        this.#channel = channel;
    }

    async publish<T>({ exchange, message, topic = "#" }: IMessage<T>) {
        this.#channel.publish(exchange, topic, Buffer.from(JSON.stringify(message)))
    }

    async consume<T>({ exchange, queue, topic = "#" }: ISubscription, action: (obj: T) => void) {
        const ch = this.#connection.createChannel({
            setup: (channel: Channel) => {
                return Promise.all([
                    channel.assertQueue(queue, { exclusive: false, autoDelete: false, durable: true }),
                    channel.assertExchange(exchange, 'topic'),
                    channel.prefetch(1),
                    channel.bindQueue(queue, exchange, topic),
                    channel.consume(queue, onMessage(action)(channel))
                ])
            }
        })
        await ch.waitForConnect();
        console.log("Listening for messages")
        this.#subscriptions.push(ch)
    }

    static from(url: string): RabbitMqBus {
        const connection = rabbit.connect([url ?? "amqp://guest:guest@localhost:5672/"]);
        const channel = connection.createChannel();
        return new RabbitMqBus(connection, channel);
    }
}