 Channel-service  Readme
 
# 0. SSH Tunnel to demo machine
To use channel-service  at http://127.0.0.1:3000
you can use ssh tunnel to see it,  create this section in ~/.ssh/config
Host demo
HostName 64.62.163.189
User ogrunner
Port 8822
 
then run: 
$ ssh demo -qTfnN -L 127.0.0.1:3000:192.1.199.122:3000

after that you use channel-service now.





# 1. Insert a new text passage post  
curl -XPOST -H 'Content-Type: application/json' -d \
 '{"user-id": 101,
  "type": "text",
  "active": true,
  "text-message" : "Honey Roasted Peanuts",
  "created-at": "Nov 25 16:00:51 PST 2015",  
  "updated-at": "Nov 25 16:00:51 PST 2015" 
  }'  \
  http://127.0.0.1:3000/v1/posts


#2. create a new image post, and store image into Mongodb
curl \
  -F "user-id=201" \
  -F "type=image" \
  -F "comment=This is an image file" \
  -F "image=@/home/ogrunner/git/channel-service/test/mylogo.jpg" \
  http://127.0.0.1:3000/v1/posts


#3. Query the total acount of the posts
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/count






# Optional part
# optional -  check Elasticsearch works  
curl -XGET -H 'Content-Type: application/json' http://127.0.0.1:9200/postindex/?pretty=true


# create a new file post, store a file into MongoDB, and create index in Elasticsearch
 curl -i -X POST -H "Content-Type: multipart/form-data" \
-F "filename=@/home/ogrunner/git/channel-service/test/test.pdf" -v http://127.0.0.1:3000/v1/uploadfile
 