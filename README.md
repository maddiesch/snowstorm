## Snowstorm

Generate a [KSUID](https://github.com/segmentio/ksuid) over a Redis protocol.

(I'm not saying this is a good idea, but it's possible.)

`docker run maddiesch/snowstorm:latest`

Options:

- `host` The Redis host to connect to

- `port` The Redis port to connect to

- `db` The Redis DB to use.

- `queue` The queue to publish new requests into.

- `prefix` The delivery queue prefix.

- `count` The number of listeners to use between 1 & 20

- `wait-ms` The amount of time (in milliseconds) a listener should wait until it checks for new messages.

### Request

`$ redis-cli rpush snowstorm-generate '{"ClientID":"my-client","RequestID":"1b"}'`

### Response

`$ redis-cli lpop snowstorm-delivery.my-client.1b`
