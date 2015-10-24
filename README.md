# CMPE273-Assignment2
# Building a location and trip planner service in Go
Part I - CRUD Location Service

## Usage
Build REST endpoints to store and retrieve locations.
This application uses MongoDB to store and retrieve data
### Install

If the packages ("github.com/julienschmidt/httprouter","gopkg.in/mgo.v2","gopkg.in/mgo.v2/bson") are not installed, go to the code file's folder and do
```
go get 
```
```
go build

```
Start the  server:
```
go run RESTfulService.go
```
### Make POST calls from Client to the URL using cURL

```
curl -X POST -d "{\"name\" : \"John Smith\",\"address\":\"123 Main St\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip\":\"94113\"}" http://localhost:8082/locations
```
```
Response:
{"id":1,"name":"John Smith","address":"123 Main St","city":"San Francisco","state":"CA","zip":"94113","coordinate":{"lat":37.7917618,"lng":-122.3943405}}
```
### Make GET calls From Client to the URL using cURL
```
curl http://localhost:8082/locations/1
```
```
Response:
{"id":1,"name":"John Smith","address":"123 Main St","city":"San Francisco","state":"CA","zip":"94113","coordinate":{"lat":37.7917618,"lng":-122.3943405}}
```
### Make PUT calls from Client to the URL using cURL
```
curl -X PUT -d "{\"address\":\"1600 Amphitheatre Parkway\",\"city\":\"Mountain View\",\"state\":\"CA\",\"zip\":\"94043\"}" http://localhost:8082/locations/1
```
```
Response:
{"id":1,"name":"John Smith","address":"1600 Amphitheatre Parkway","city":"Mountain View","state":"CA","zip":"94043","coordinate":{"lat":37.4220352,"lng":-122.0841244}}
```
### Make DELETE calls from Client to the URL using cURL
```
curl -X DELETE http://localhost:8082/locations/1
```
```
Response:
Deleted a record

```
Comments: Run the RESTfulService.go file first using "go run RESTfulService.go" (without double quotes) and in separate command prompt run the cURL command as specified above
The same can be achieved using postman(application) where you can post raw values to URL and retrieve data
