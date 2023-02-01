# GoTrader

<img src="https://raw.githubusercontent.com/egonelbre/gophers/10cc13c5e29555ec23f689dc985c157a8d4692ab/vector/friends/crash-dummy.svg" alt="gopher" width="30%"/>


A cryptocurrency trading bot supporting kraken futures written in Golang.

---

## Current Features

* Support for sending any order on kraken futures (mkt, lmt, etc...)
* Support trading on kraken futures using stop loss & take profit indicator
* REST API support for kraken futures
* Websocket API support for kraken futures
* JWT Token auth support with deleting token on logout from device
* Telegram bot 
* Swagger documentation

---

## Planned Features

* Support multiple kraken api tokens

---

## Exchange support table

| Exchange            | REST API | Streaming API | 
|---------------------|----------|---------------|
| Kraken futures demo | Yes      |  Yes          |
| Kraken futures      | Yes      |  Yes          |

---

## Tech stack

* [Go](https://github.com/golang/go)
* [Postgres](https://www.postgresql.org)
* [Redis](https://redis.io)

---


## Local installation of server
### Linux/OSX

*
    ```shell
    git clone {this repo}
    cd {this repo}/course_project/trade-bot
    ```

* #### Assume you have ```config.yml``` or ```config.yaml``` file in configs folder of type:
    ```yaml
    server:
      port: (int) 
      websocket:
        readBufferSize: (int) 1024 by derfault
        writeBufferSize: (int) 1024 by default
        checkOrigin: (true | false) true by default
    
    client:
      # url of server
      url: (string) example - http://localhost:8000
      
    telegram:
      apiToken: (string) your telegram api token from bot father
      webhookUrl: (string) example service for webhooks - ngrok
    
    postgreDatabase:
      host: (string) example - localhost
      port: (string) eample - 8000
      username: (string)
      dbname:  (string)
      sslmode: (string)
    
    redisDatabase:
      host: (string) example - localhost
      port: (string)
    
    kraken:
      apiurl: (string)
    
    krakenWS:
      requests:
        writeWaitInSeconds: (int) 10 by default
        pongWaitInSeconds: (int) 60 by default
        pingPeriodInSeconds: (int) 10 by default
        maxMessageSize: (int) 512 by default
      kraken:
        wsapiurl: (string)
    ```

* #### Assume you have ```.env``` file at the root of project with following:
    ```.dotenv
    DB_PASSWORD = (your postgres db password)
    
    JWT_ACCESS_SIGNING_KEY = (key for signing jwt tokens)
    
    PUBLIC_API_KEY = (public key from kraken futures)
    PRIVATE_API_KEY = (private key from kraken futures)
    ```

* #### Run postgres with settings from your config file
    ```shell
    # Example using docker
    
    docker pull postgres
    docker run --name postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d postgres
    ```

* #### Run redis with settings from your config file
    ```shell
    # Example using docker
    
    docker pull redis
    docker run --name redis -p 6379:6379 -d redis
    ```

* #### Run migrate files for postgres using ```migrate```
* __migrate installation__
* 
    ```
    curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
    apt-get update
    apt-get install -y migrate  
    ```
* __run migrate__
* 
    ```shell
    migrate -path ./schema -database 'postgres://{postgres_username}:{postgres_password}@{host}:{port}/postgres?sslmode={sslmode}' up
    ```

* #### Then run server

    ```shell
    go run cmd/api/main.go
    ```

---

## Installation of server using Docker
### Linux/OSX

* #### You'll need Docker Compose
* #### Make sure you have all config files like in local installation with some fixes
* __Fixes__
    ```yaml
      postgreDatabase:
        host: db (like db service name in docker-compose file)
      
      ###
      
      redisDatabase:
        host: redis (like redis service name in docker-compose file)
    ```

* #### Run docker-compose
    ```shell
    docker-compose up --build server
    ```
* #### Now simply run migration files like in local installation

---

## Swagger

__When server started:__ ```url: http://{host}:{port}/swagger/index.html```

---
