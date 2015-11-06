curl -XPOST -H 'Content-Type: application/json' -d '{"name": "Bob Smith", "gender": "male", "age": 50}' http://127.0.0.1:3000/user
{"id":"5497246c380a967ff1000003","name":"Bob Smith","gender":"male","age":50}



curl -H "Content-Type: application/json" -X POST -d '{"userId": 1192293, "orderNumber": "J2326", "orderTotal": 34.45, "orderDate": "2015-09-30T07:00:00.000Z", "orderState": "paid" }' -v http://127.0.0.1:8098/v1/orders
echo "Test case 4 Done "



  #Insert a new text passage post  
 curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:3000/v1/posts 
 
   #upload a new image post  
 curl -XPOST -H 'Content-Type: application/json' -d '{"user-id": 201, "type": "image","active": true,  "title" : "mylogo",  "comment" : "This is an image file" , "link" : "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg"}' http://127.0.0.1:3000/v1/posts 


curl \
  -F "user-id=201" \
  -F "type=image" \
  -F "comment=This is an image file" \
  -F "image=@/Users/huazhang/git/post-message-service/test/mylogo.jpg" \
  http://127.0.0.1:3000/v1/posts
  


 #Query an existing post
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/posts/563cfa9dd12619398ed7c71c



 # Insert a new product 
 curl -XPOST -H 'Content-Type: application/json' -d '{"Name": "peanuts", "Description": "Honey Roasted peanuts", "MetaDescription": "Good taste food", "Permalink": "www.peanuts.com"}' http://127.0.0.1:3000/v1/products 

 curl -XPOST -H 'Content-Type: application/json' -d '{"Name": "boil egg", "Description": "Safeway signure egg", "MetaDescription": " best by 10 days", "Permalink": "www.goodegg.net"}' http://127.0.0.1:3000/v1/products 



 # Query an existing product
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/products/1

 # Query an existing product with q (search string)
curl -H "Content-Type: application/json" -X GET -v http://127.0.0.1:3000/v1/products/1 -d '{"q" : "boil egg"}'