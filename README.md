# Log-Broker
A systemc consists of 3(4) parts to communicate with each other and send logs to the destinations with the rate of at least 10000 message per seconds and 50B to 8KB message size.

1. Sender: Produces messages with the given rate and size
2. Reciever: Gets the messages and send them to the broker
3. Broker: Gets messages and add them to message queue and send them to the destination. if destination gets offline it will retries and keep messages and log them.
4. Destination: Gets the messages and print their size

## Architecture
![architecture](https://github.com/alipourhabibi/log-broker/blob/master/log-broker.png?raw=true)
