# Channel-Service
post message, upload file/photo/video


## Environment
1. make sure the Mongodb start  
  database name :  channel-service  

2. start elasticearch  
  cd /opt/elasticsearch  
  ./bin/elasticsearch -d   



## Configuration
    "Elasticsearch": "http://127.0.0.1:9200",  
	"Server" :  "127.0.0.1:3000",   
	"Database" : "channel_service"  

## Build
make

## Run
sudo start channel-service

## Stop
sudo stop channel-service


## Test
curl -XPOST -H 'Content-Type: application/json' -d \
 '{"user-id": 101, "type": "text","active": true,  "text-message" : "Honey Roasted Peanuts" }' http://127.0.0.1:3000/v1/posts 




