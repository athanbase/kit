package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	pb "github.com/athanbase/kit/protobuf/generate/tutorialpb"

	"google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func main() {
	// 北京时间
	fmt.Println(time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().In(time.FixedZone("CST", 8*60*60)))
	writeBookToFile()
	readBookFromFile()
}

func readBookFromFile() {
	in, err := ioutil.ReadFile("book")
	if err != nil {
		panic(err)
	}
	book := &pb.AddressBook{}

	if err := proto.Unmarshal(in, book); err != nil {
		panic(err)
	}
	fmt.Println(book.GetPeople())
}

func writeBookToFile() {
	p := pb.Person{
		Name:        "John Doe",
		Id:          1,
		Email:       "jdoe@foo.com",
		Phones:      []*pb.Person_PhoneNumber{{Number: "123232", Type: pb.Person_HOME}},
		LastUpdated: &timestamp.Timestamp{Seconds: 13333333},
	}
	book := pb.AddressBook{}
	book.People = append(book.People, &p)
	out, err := proto.Marshal(&book)
	if err != nil {
		panic("marshal book error: " + err.Error())
	}
	err = ioutil.WriteFile("book", out, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("write sucess")
}

var p = pb.Person{
	Name:        "John Doe",
	Id:          1,
	Email:       "jdoe@foo.com",
	Phones:      []*pb.Person_PhoneNumber{{Number: "123232", Type: pb.Person_HOME}},
	LastUpdated: &timestamp.Timestamp{Seconds: 13333333},
}
var book = pb.AddressBook{People: []*pb.Person{&p}}

func marshalJson() ([]byte, error) {
	return json.Marshal(&book)
}

func marshalPb() ([]byte, error) {
	return proto.Marshal(&book)
}
