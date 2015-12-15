#Insert a new text passage post  
curl -XPOST -H 'Content-Type: application/json' -d \
 '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:3000/v1/posts 

curl -XPOST -H 'Content-Type: application/json' -d \
 '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:8082/test 


curl -XPOST -H 'Content-Type: application/json' -d \
 '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts", "created-at": "Nov 25 16:00:51 PST 2015",  "updated-at": "Nov 25 16:00:51 PST 2015" }' http://127.0.0.1:3000/v1/posts

 
curl -XGET -H 'Content-Type: application/json' http://127.0.0.1:9200/postindex/?pretty=true

#create a new image post, and store image into Mongodb
curl -XPOST -H 'Content-Type: application/json' -d \
'{"user-id": 201, "type": "image","active": true,  "title" : "mylogo", "link" : "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg", name: "logo", "comment" : "This is an image file" }' http://127.0.0.1:3000/v1/posts 

curl \
  -F "user-id=201" \
  -F "type=image" \
  -F "comment=This is an image file" \
  -F "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg" \
  http://127.0.0.1:3000/v1/posts

# create a new file post, store a file into MongoDB, and create index in Elasticsearch
 curl -i -X POST -H "Content-Type: multipart/form-data" \
-F "filename=@/Users/huazhang/git/channel-service/test/test.pdf" -v http://127.0.0.1:3000/v1/uploadfile
 

# Query the total acount of the posts
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/count

# Query the posts with a filter string
curl -XGET 'http://127.0.0.1:3000/v1/posts/?q=user-id:101'

curl -XGET 'http://127.0.0.1:3000/v1/posts/??limit=10&q=user-id:101'

curl -XGET 'http://127.0.0.1:3000/v1/posts/??limit=10&offset=xx&q=user-id:101'


# Delete post  DELETE /v1/posts/post-id
curl -XPUT 'http://127.0.0.1:3000/v1/posts/564fa99fd1261920bfa52557'




# Create comments 
curl -XPOST -H 'Content-Type: application/json' \
     -d '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' \
          'http://127.0.0.1:3000/v1/posts/564cf977d12619192199b1b3/comments'

# List Comments
curl -XGET -H "Content-Type: application/json"  http://127.0.0.1:3000/v1/posts/564cf977d12619192199b1b3/comments