#source-listener-service-go

## To build
```bash
go build -o=app
chmod a+x app
```

## Prepare database:
```mysql
create table store (
  client_id varchar(100) not null, 
  object_type varchar(30) not null, 
  object_id varchar(100) not null
);

create unique index store_unique_index
    on store (client_id, object_type, object_id);
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
* AMQP_EXCAHNGE_NAME
* AMQP_QUEUE

```bash
./app
```