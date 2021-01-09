//@ts-check
const rabbit = require("amqp-connection-manager")

// export function startPublisher() {
//     rabbit.connect
// }


/**
 * @param {import("amqp-connection-manager").ChannelWrapper} channel
 */
async function publishToRabbitMq(channel, { exchange, message, topic = "#"}) {
    channel.publish(exchange, topic, Buffer.from(JSON.stringify(message)))
}

module.exports = { publishToRabbitMq }