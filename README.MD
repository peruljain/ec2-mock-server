# Problem Statement
How ec2 server creation works from user perspective ?

Short Polling : Client will call server in every 5 seconds
Long Polling : Client will call server and server will call database/another-server in every 5 seconds


In this example, we will implement short-polling

## Steps to follow

1. docker-compose up -d
2. ./table_create.sh
3. go run main.go
4. Paste in browser : http://localhost:8080/create
![alt text](image.png)
5. paste in browser : http://localhost:8080/status?id=2
![alt text](image-1.png)


In second screenshot, you see created only after 20 seconds because we have logic in create server
such that after 20 status second , server status changes to created