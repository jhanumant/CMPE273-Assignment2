package main

import ("fmt"
		"net/http"
		"encoding/json"
		"io/ioutil"
		"gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        "log"
        "github.com/julienschmidt/httprouter"
        "strconv"
        "strings"
        "time"
)

const(
	MongoDBHosts = "ds043324.mongolab.com:43324" 
 	AuthDatabase = "users" 
 	AuthUserName = "jhanumant" 
 	AuthPassword = "hanumant" 
)

type Response struct{
	Id int `json:"id" bson:"id"` 
	Name string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
	City string `json:"city" bson:"city"`
	State string `json:"state" bson:"state"`
	Zip string `json:"zip" bson:"zip"`
	Coor Coordinate `json:"coordinate" bson:"coordinate"`
}

type Coordinate struct{
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type Result struct{
	Id int  `bson:"id"`
}

func main() {
	mux := httprouter.New()
    mux.POST("/locations",PostLocations)
    mux.GET("/locations/:locationid",GetLocations)
    mux.PUT("/locations/:locationid",PutLocations)
    mux.DELETE("/locations/:locationid",DeleteLocations)
     server := http.Server{
            Addr:        "0.0.0.0:8082",
            Handler: mux,
    }
    server.ListenAndServe()

}

func DeleteLocations(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	params,_:= strconv.Atoi(p.ByName("locationid"))
	dBDialInfo := &mgo.DialInfo{ 
	Addrs:    []string{MongoDBHosts}, 
 	Timeout:  60 * time.Second, 
 	Database: AuthDatabase, 
 	Username: AuthUserName, 
 	Password: AuthPassword, 
	} 
	session, err := mgo.DialWithInfo(dBDialInfo)
	if(err!=nil){
		defer session.Close()
	}
	c := session.DB("users").C("locations")	
	err = c.Remove(bson.M{"id":params})
	if(err!=nil){
		log.Printf("RunQuery : ERROR : %s\n", err) 
		fmt.Fprintln(rw,err)
				return
	}else{
		fmt.Fprintln(rw,"Deleted a record")
	}
}

func PutLocations(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	rq,err:= ioutil.ReadAll(request.Body)
	params,_:= strconv.Atoi(p.ByName("locationid"))
	dBDialInfo := &mgo.DialInfo{ 
	Addrs:    []string{MongoDBHosts}, 
 	Timeout:  60 * time.Second, 
 	Database: AuthDatabase, 
 	Username: AuthUserName, 
 	Password: AuthPassword, 
	} 
	session, err := mgo.DialWithInfo(dBDialInfo)
	if(err!=nil){
		defer session.Close()
	}
	var data Response
	c := session.DB("users").C("locations")	
	json.Unmarshal(rq,&data)
	add:= strings.Replace(data.Address," ","+",-1)
	add= add + ",+"+ strings.Replace(data.City," ","+",-1)
	add= add + ",+" + data.State
	var jsonLocation interface{}
	response,err:= http.Get("http://maps.google.com/maps/api/geocode/json?address="+add+"&sensor=false")
	if err!=nil{
		fmt.Println("Error:",err)
	}else{
		defer response.Body.Close()
		Locationcontents,_:= ioutil.ReadAll(response.Body)
		json.Unmarshal(Locationcontents,&jsonLocation)		
		latitude:= (jsonLocation.(map[string]interface{})["results"]).([]interface{})[0].(map[string]interface{})["geometry"].
		           (map[string]interface{})["location"].(map[string]interface{})["lat"]
		longitude:= (jsonLocation.(map[string]interface{})["results"]).([]interface{})[0].(map[string]interface{})["geometry"].
		           (map[string]interface{})["location"].(map[string]interface{})["lng"]		
		data.Coor.Lat=latitude.(float64)
		data.Coor.Lng=longitude.(float64)
		err = c.Update(bson.M{"id":params},bson.M{"$set":bson.M{"address":data.Address,"city":data.City,"state":data.State,"zip":data.Zip,"coordinate.lat":data.Coor.Lat,"coordinate.lng":data.Coor.Lng}})
		if(err!=nil){
			log.Printf("RunQuery : ERROR : %s\n", err) 
			fmt.Fprintln(rw,err)
					return
		}else{
			err = c.Find(bson.M{"id":params}).Select(bson.M{"_id":0}).One(&data)
			if(err!=nil){
				log.Printf("RunQuery : ERROR : %s\n", err) 
				fmt.Fprintln(rw,err)
					return
			}
			result,_:=json.Marshal(data)
			fmt.Fprintln(rw,string(result))	
		}
	}
}

func GetLocations(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	params,_:= strconv.Atoi(p.ByName("locationid"))
	dBDialInfo := &mgo.DialInfo{ 
	Addrs:    []string{MongoDBHosts}, 
 	Timeout:  60 * time.Second, 
 	Database: AuthDatabase, 
 	Username: AuthUserName, 
 	Password: AuthPassword, 
	} 
	session, err := mgo.DialWithInfo(dBDialInfo)
	if(err!=nil){
		defer session.Close()
	}
	var data Response
	c := session.DB("users").C("locations")
	err = c.Find(bson.M{"id":params}).Select(bson.M{"_id":0}).One(&data)
	if(err!=nil){
		log.Printf("RunQuery : ERROR : %s\n", err) 
		fmt.Fprintln(rw,err)
				return
	}else{
		result,_:=json.Marshal(data)
		fmt.Fprintln(rw,string(result))	
	}
}

func PostLocations(rw http.ResponseWriter, request *http.Request,p httprouter.Params) {
	
	var idResult Result
	dBDialInfo := &mgo.DialInfo{ 
	Addrs:    []string{MongoDBHosts}, 
 	Timeout:  60 * time.Second, 
 	Database: AuthDatabase, 
 	Username: AuthUserName, 
 	Password: AuthPassword, 
	} 
	session, err := mgo.DialWithInfo(dBDialInfo)
	if(err!=nil){
		defer session.Close()
	}
	var jsonLocation interface{}
	var data Response
	rq,err:= ioutil.ReadAll(request.Body)
	json.Unmarshal(rq,&data)
	add:= strings.Replace(data.Address," ","+",-1)
	add= add + ",+"+ strings.Replace(data.City," ","+",-1)
	add= add + ",+" + data.State
	response,err:= http.Get("http://maps.google.com/maps/api/geocode/json?address="+add+"&sensor=false")
	if err!=nil{
		fmt.Println("Error:",err)
	}else{
		defer response.Body.Close()
		Locationcontents,_:= ioutil.ReadAll(response.Body)
		json.Unmarshal(Locationcontents,&jsonLocation)
		latitude:= (jsonLocation.(map[string]interface{})["results"]).([]interface{})[0].(map[string]interface{})["geometry"].
		           (map[string]interface{})["location"].(map[string]interface{})["lat"]
		longitude:= (jsonLocation.(map[string]interface{})["results"]).([]interface{})[0].(map[string]interface{})["geometry"].
		           (map[string]interface{})["location"].(map[string]interface{})["lng"]
		data.Coor.Lat=latitude.(float64)
		data.Coor.Lng=longitude.(float64)
		c := session.DB("users").C("locations")
		idResult.Id=0
		count,_:=c.Count()
		if(count > 0){
			err := c.Find(nil).Select(bson.M{"id":1}).Sort("-id").One(&idResult)
			if(err!=nil){
				log.Printf("RunQuery : ERROR : %s\n", err) 
				fmt.Fprintln(rw,err)
				return 
			}
			data.Id = idResult.Id + 1
	        err = c.Insert(data)
	        if err != nil {
	                log.Fatal(err)
	        }
	        result,_:=json.Marshal(data)
			fmt.Fprintln(rw,string(result))
		}else{
			data.Id = idResult.Id + 1
	        err = c.Insert(data)
	        if err != nil {
	                log.Fatal(err)
	        }
	        result,_:=json.Marshal(data)
			fmt.Fprintln(rw,string(result))	
		}
	}
}

