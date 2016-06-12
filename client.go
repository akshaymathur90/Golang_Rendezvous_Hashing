package main

import (
	"strconv"
	"hash/crc64"
	"os"
	"strings"
	"log"
	"net/http"
	"io/ioutil"

)

type NodeRing struct {
	Nodes []string
}

func (r *NodeRing) AddNode(node string) {
	//node, _ := strconv.ParseUint(id,10,64)
	r.Nodes = append(r.Nodes, node)
	
}

func getHash(key string) uint64{
	crcTable := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum([]byte(key), crcTable)

}
func getWeight(node uint64,key string) uint64{
	a := uint64(1103515245)
    b := uint64(12345)
    keyHash := getHash(key)
    return (a * ((a * node + b) ^ keyHash) + b) % (2^31)
}
func (r *NodeRing) Get(key string) string {
	
	var weights []uint64
	wmap:= make(map[uint64]string)
	
	for _,n := range r.Nodes{
		nUint, _ := strconv.ParseUint(n,10,64)
		w:=getWeight(nUint,key)
		weights=append(weights,w)
		wmap[w]=n
		
	}
	max := weights[0] // assume first value is the smallest

         for _, value := range weights {
                 if value > max {
                         max = value // found another smaller value, replace previous value in max
                 }
         }
       //log.Println(wmap[max]);
       return wmap[max]
}
func main(){

	//Get nodes from command line
	portRange:= os.Args[1]
	portNums:=strings.Split(portRange,"-")
	//log.Println(portNums[0])
	//log.Println(portNums[1])
	portStart,_:=strconv.ParseInt(portNums[0],10,64)
	portEnd,_:=strconv.ParseInt(portNums[1],10,64)
	log.Println(portStart)
	log.Println(portEnd)
	//Get the key and values to be stored
	inputdata:= os.Args[2]
	keyValPairs:=strings.Split(inputdata,",")
	log.Println(keyValPairs)
	
	myRing:=&NodeRing{Nodes : []string{}}
	
	for i:=portStart;i<=portEnd;i++{
		myRing.AddNode(strconv.FormatInt(i,10))
	}
	
	for _,k:= range keyValPairs{
		keyVal:=strings.Split(k,"->")
		log.Println(keyVal[0]+"--"+keyVal[1])
		getNode:=myRing.Get(keyVal[0])
		log.Println(getNode)
		
		url:="http://localhost:"+getNode+"/"+keyVal[0]+"/"+keyVal[1]
		log.Println(url)
		req, err :=http.NewRequest("PUT", url,nil)
		req.Header.Set("X-Custom-Header", "myvalue")
    	req.Header.Set("Content-Type", "application/json")

    	client := &http.Client{}
	    resp, err := client.Do(req)
	    if err != nil {
	        panic(err)
	    }
	    defer resp.Body.Close()
	
	    log.Println("response Status:", resp.Status)
	    body, _ := ioutil.ReadAll(resp.Body)
	    log.Println("response Body:", string(body))
	}
}