package main

import (
	"container/list"
	"encoding/json"
	"fmt"

	"net"
	"time"
)

type Tcp struct {
	ServAddr string
}

type msg struct {
	Mode    byte     //0:broadcast ips 1:send msg
	LAddr   string   //本地ip
	Raddr   []string //远程ip
	Content string
}

var users *list.List

func init() {
	users = list.New()
}
func (t Tcp) write(conn net.Conn) {
	defer conn.Close()
	var m msg

	for {
		m = msg{
			Mode:    0,
			Content: "heelo",
			Raddr:   make([]string, 1),
		}
		for e := users.Front(); e != nil; e = e.Next() {
			ip := e.Value.(string)
			m.Raddr = append(m.Raddr, ip)
		}
		bytes, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(bytes)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

func (t Tcp) read(conn net.Conn) {
	defer conn.Close()
	inbuf := make([]byte, 1024)
	for {
		n, err := conn.Read(inbuf)
		if err != nil {
			fmt.Println("client 连接出错：")
			return
		}
		if n == 0 {
			continue
		}

		var msg msg

		err = json.Unmarshal(inbuf[:n], &msg)

		if err == nil {
			if msg.Mode == 1 {
				for _, v := range msg.Raddr {
					for e := users.Front(); e != nil; e = e.Next() {
						ip := e.Value.(string)
						if v == ip {
							users.Remove(e)
						}
					}
					users.PushBack(v)
				}
			} else if msg.Mode == 0 {
				fmt.Println("get msg from:", msg.LAddr, msg.Content)
			}
		} else {
			panic(err)
		}
	}
}

func (t Tcp) handleConn(conn net.Conn) {
	go t.read(conn)
	t.write(conn)
}

func (t *Tcp) serve() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.ServAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		panic(conn)
	}

	go t.handleConn(conn)
	time.Sleep(time.Hour)
}
