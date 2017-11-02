# mvcPractice

A backend server written on Go. It implements microservice architecture and can be used either in Docker container or directly in OS.

## Install
```bash
go get -u github.com/ferux/mvcPractice
go build -ldflags "-X main.outsideDocker=true" -o "app_name" .
```
If you decide to run outside the docker container you can skip -ldflags parameter.  

## Usage in container

If you are using this server inside container, you should specify enviroment varibales in your Dockerfile or docker-compose.yml (the example below is for docker-compose.yml):  
```docker-compose
environment: 
     - IP=backend 
     - PORT=8080 
     - DB_IP=db 
     - DB_PORT=5432 
     - DB_USER=user
     - DB_PWD=password 
     - DRIVER=postgres 
     - DB=simplemvc
     - EXT_IP=http://${EXTERNAL_IP}
```

## Usage standalone

If you are using the server outside of container, you should specify parameters in your console while running the application.
This server requires the following params:  

Parameter: ip (default: 127.0.0.1) Description: Hosting IP  
Parameter: port (default: 8080) Description: Hosting Port  
Parameter: dbAddr (default: 127.0.0.1) Description: Database IP  
Parameter: dbPort (default: 5432) Description: Database Port  
Parameter: user (default: user) Description: Database User  
Parameter: password (default: pass) Description: Database Password  
Parameter: db (default: default) Description: Database Name  
Parameter: driver (default: postgres) Description: Database driver  
Parameter: extIP (default: 127.0.0.1) Description: External IP for CORS

## License

MIT