package tcpserver

import proto "github.com/golang/protobuf/proto"

import (
	
	"net"
	"fmt"
	"time"
	"bytes"
	"encoding/binary"
	"sync/atomic"
)

type client struct {

	nc net.Conn

	receivedBytes    []byte
	
	sendBytesChan chan []byte

	serverInstance *server

	id uint64
}

func newclient(conn net.Conn, s *server) (*client){
		
	c := &client{nc:conn}
	c.id = atomic.AddUint64(&s.gcid, 1)

	c.sendBytesChan = make(chan []byte, 100)

	go c.readloop()
	go c.writeloop()

	c.serverInstance = s

	return c
}

func (c *client) readloop() {
	
	bytes := make([]byte, 32768)
	for {

		c.nc.SetReadDeadline(time.Now().Add(5 * time.Minute))
		i, err := c.nc.Read(bytes)
		if err != nil {
			
			// handle read error
			// close connection
			return
		}

		if i == 0 {
			continue
		}

		data := bytes[:i]
		if c.receivedBytes != nil {
			c.receivedBytes = append(c.receivedBytes, data...)
		} else {
			c.receivedBytes = data
		}

		go c.serverInstance.dispatch(c, c.receivedBytes);


		c.receivedBytes = c.receivedBytes[len(c.receivedBytes):]
	}
}

func (c *client) send(pakcetid uint16, pakcetdata interface{}) {
	
	dataSending, err := proto.Marshal(pakcetdata.(proto.Message))

	if err != nil {

		// error handle

		return
	}

	var sendbuffer bytes.Buffer

	handleId := uint16(pakcetid)
	size := uint16(len(dataSending))
	fmt.Println("size : ", (size))
	fmt.Println("handleId : ", (handleId))

	ext := uint32(1)
	
	binary.Write(&sendbuffer, binary.BigEndian, size);
	
	binary.Write(&sendbuffer, binary.BigEndian, handleId);
	
	binary.Write(&sendbuffer, binary.BigEndian, ext);

	binary.Write(&sendbuffer, binary.BigEndian, dataSending);

	data := sendbuffer.Bytes()

	c.sendBytesChan <- data
}

func (c *client) writeloop() {
	
	for {

		select {

		case data := <-c.sendBytesChan:

			c.nc.Write(data)
		}
	}
}