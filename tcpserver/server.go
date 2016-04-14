package tcpserver

import proto "github.com/golang/protobuf/proto"

import (

	"net"
	"fmt"
	"puzzleclash/proto_struct"
)

type server struct {

	listener net.Listener
	clients  map[uint64]*client
	gcid     uint64
}

func Newserver() *server {
	
	s := &server{gcid:0}
	s.clients = make(map[uint64]*client, 100)

	return s
}

func (s *server) Start() {

	l, err := net.Listen("tcp", "127.0.0.1:8086");

	if err !=nil {

		//handle listen error
	}

	s.listener = l;

	for {

		connection,err := s.listener.Accept();

		if err != nil {
			
			// handle err
		}

		// handle connection
		c := newclient(connection, s);

		s.clients[c.id] = c
	}
}

func (s *server) dispatch(c *client, data []byte) {

	fmt.Println("Message Received: ", string(data))

	addressbook := &proto_struct.AddressBook{}

	p := &proto_struct.Person{

		Name : proto.String("john"),
		Id : proto.Int32(1),
		Email : proto.String("bysdtd@firewing.com"),
		Phone : []*proto_struct.Person_PhoneNumber{

				{Type: proto_struct.Person_HOME.Enum(), Number: proto.String("13811825285")},
				{Type: proto_struct.Person_MOBILE.Enum(), Number: proto.String("18101037263")},
			},
	}

	//addressbook.Person = []*proto_struct.Person{}

	//addressbook.Person = append(addressbook.Person, p)
	addressbook.Person = p
	
	// testproto := &proto_struct.TestProto{

	// 	Index : proto.Int32(147),
	// 	Name : proto.String("bYsdtd"),
	// 	Numbers : []string{
	// 		"147",
	// 		"258",
	// 		"369",
	// 	},
	// }

	c.send(1, addressbook)
}