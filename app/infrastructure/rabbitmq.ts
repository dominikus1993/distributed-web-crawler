import { ChannelWrapper } from "amqp-connection-manager"

export interface IMessage<T> {
    readonly exchange: string
    readonly message: T
    readonly topic?: string
}

export async function publishToRabbitMq<T>(channel: ChannelWrapper, { exchange, message, topic = "#"}: IMessage<T>) {
    channel.publish(exchange, topic, Buffer.from(JSON.stringify(message)))
}
