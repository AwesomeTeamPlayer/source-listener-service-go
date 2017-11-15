# source-listener-service-go

This service is responsible for:
1) storing information about which client listens on changes on which object
2) receiving all events from the system and providing them to specific clients using appropriate queues 

To store pairs client - object we use 3 remote methods: RegisterClient, UnregisterClient, RemoveClient.
#### "RegisterClient"
Input data: 
```json
{
   "QueueName": "queue-123",
   "ClientId": "abc123defGHI",
   "ObjectType": "user",
   "ObjectId": "123"
}
``` 

#### "UnregisterClient"
Input data: 
```json
{
   "ClientId": "abc123defGHI",
   "ObjectType": "user",
   "ObjectId": "123"
}
``` 

#### "RemoveClient"
This method remove specific client assigned to all objects
Input data: 
```json
{
   "ClientId": "abc123defGHI"
}
``` 

The worker connect _AMQP_QUEUE_ queue and consume messages from it:
1. Worker takes one message. The message has to have format:
    ```json
    {
      "ObjectType": "user",
      "ObjectId": "123",
      "data": {
        // event details 
      }
    } 
    ```
2. Worker takes all clients listen for this ObjectType and ObjectId. For each client:

2.1. Worker connects to a queue assigned to the specific client.

2.2. If the queue does not exist worker finish this loop and take next client.


2.3. Otherwise the worker adds `"ClientId":"abc123defGHI"` to received message:
```json
{
  "ClientId": "abc123defGHI",
  "ObjectType": "user",
  "ObjectId": "123",
  "data": {
    // event details 
  }
} 
```
2.4. ... and pushes this message to client's queue.  


## To build
```bash
go build -o=app
chmod a+x app
```

## Prepare database:
```mysql
create table store (
  queue_name varchar(100) not null, 
  client_id varchar(100) not null, 
  object_type varchar(30) not null, 
  object_id varchar(100) not null
);

create unique index store_unique_index
    on store (queue_name, client_id, object_type, object_id);
```

## To run built file

Required ENV variables:

* APP_PORT

* MYSQL_HOST
* MYSQL_PORT
* MYSQL_USER
* MYSQL_PASSWORD
* MYSQL_DATABASE

* AMQP_HOST
* AMQP_PORT
* AMQP_USER
* AMQP_PASSWORD
* AMQP_QUEUE

```bash
./app
```