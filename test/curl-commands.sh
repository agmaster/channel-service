#Insert a new text passage post  
curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:3000/v1/posts 
 
#upload a new image post  
curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 201, "type": "image","active": true,  "title" : "mylogo",  "comment" : "This is an image file" , "link" : "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg"}' http://127.0.0.1:3000/v1/posts 

# Query the total acount of the posts
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/count

# Query the posts with a filter string
curl -XGET 'http://127.0.0.1:3000/v1/posts/?q=user-id:101'

curl -XGET 'http://127.0.0.1:3000/v1/posts/??limit=10&q=user-id:101'

curl -XGET 'http://127.0.0.1:3000/v1/posts/??limit=10&offset=xx&q=user-id:101'

curl \
  -F "user-id=201" \
  -F "type=image" \
  -F "comment=This is an image file" \
  -F "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg" \
  http://127.0.0.1:3000/v1/posts
  


# upload a file into MongoDB
 curl -i -X POST -H "Content-Type: multipart/form-data" \
-F "filename=@/Users/huazhang/test.txt" -v http://127.0.0.1:3000/v1/uploadfile
 
 #create a new image post into MongoDB
 curl -XPOST -H 'Content-Type: application/json' \
      -d '{"user-id": 301, "type": "image","active": true, "title" : "mylogo",  "comment" : "This is an image file" , "link" : "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg"}'\
            http://127.0.0.1:3000/v1/posts 
 
 
 curl -F  "file=@/home/devop/elk/test/lograge_production.log"  --cert "/home/devop/elk/logstash-forwarder/server.crt" -H  -v http://192.168.30.35:6782/
 
 
 
 curl -F  "file=@/home/devop/elk/test/lograge_production.log"  --cacert server.crt -H  -v "http://192.168.199.35:6782/"
 
 
  curl -F  "file=@/home/devop/elk/test/lograge_production.log"  -H -v "http://192.168.199.35:9200/"
  
  
  logger  -t your_tag -p "local1.info" --file /home/devop/elk/test/lograge_production.log --server 192.168.199.35 --tcp --port 9200 
  
  curl -XPOST 'http://192.168.199.35:9200/posts' -d  @/home/devop/elk/test/lograge_production.log
  
  
  curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' 'http://192.168.199.35:9200/comments'