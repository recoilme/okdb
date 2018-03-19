package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/recoilme/okdb/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	var err error
	var f = "db/1.db"
	var f2 = "db/2.db"

	conn, err = grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewOkdbClient(conn)

	// SayOk
	response, err := c.SayOk(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Message)

	// Set
	_, err = c.Set(context.Background(), &api.CmdSet{File: f, Key: []byte("1"), Val: []byte("1")})
	if err != nil {
		log.Println(err)
	}

	// Get
	b, err := c.Get(context.Background(), &api.CmdGet{File: f, Key: []byte("1")})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("b:", b.Bytes, "str:", string(b.Bytes))

	//Sets
	var a [][]byte
	for i := 0; i < 10; i++ {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(i))
		a = append(a, bs)
		a = append(a, bs)
	}
	_, err = c.Sets(context.Background(), &api.CmdSets{File: f2, Keys: a})
	if err != nil {
		log.Println(err)
	}

	// Keys
	cmdKeys := &api.CmdKeys{File: f2, From: nil, Limit: 2, Offset: 0, Asc: false}
	keys, err := c.Keys(context.Background(), cmdKeys)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("keys:%+v\n", keys.Keys)
	for i, key := range keys.Keys {
		fmt.Println(i, binary.BigEndian.Uint32(key))
	}

	// Gets
	resPairs, err := c.Gets(context.Background(), &api.CmdGets{File: f2, Keys: keys.Keys})
	if err != nil {
		log.Println(err)
	}
	for k, v := range resPairs.Pairs {
		if k%2 == 0 {
			//key
			fmt.Println("Key:", v)
		} else {
			//val
			fmt.Println("Val:", v)
		}
	}
}
