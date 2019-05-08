package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net"
)

type msg struct {
	Mode    byte
	LAddr   string   //本地ip
	Raddr   []string //远程ip
	Content string
}

var conns *list.List

type Tcp struct {
	ServAddr string
}

func init() {
	conns = list.New()
}

func (t Tcp) sendMsg(m msg) {
	for e := conns.Front(); e != nil; e = e.Next() {
		conn := e.Value.(net.Conn)
		curIP := conn.RemoteAddr().String()

		for _, ip := range m.Raddr {
			if ip == curIP {
				bytes, _ := json.Marshal(m)

				conn.Write(bytes)
				fmt.Println("send msg to :", ip)
			}
		}

	}
}

func (t Tcp) handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			for e := conns.Front(); e != nil; e = e.Next() {
				if e.Value.(net.Conn) == conn {
					conns.Remove(e)
					fmt.Print("client %s is offline./n", conn.RemoteAddr().String())
					return
				}

			}

		}
		if n == 0 {
			continue
		}

		msg := msg{}
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			continue
		}
		fmt.Println(msg)
		go t.sendMsg(msg)
	}
}

func (t *Tcp) serve() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.ServAddr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("有一个客户端上线：", conn.RemoteAddr().String())

		conns.PushBack(conn)
		t.broadcast()
		go t.handleClient(conn)
	}
}

//广播所有上线的ip地址
func (t Tcp) broadcast() {
	msg := msg{
		Mode:  1,
		Raddr: make([]string, 1),
	}
	for e := conns.Front(); e != nil; e = e.Next() {
		c := e.Value.(net.Conn)
		ip := c.RemoteAddr().String()
		msg.Raddr = append(msg.Raddr, ip)

	}

	for e := conns.Front(); e != nil; e = e.Next() {
		c := e.Value.(net.Conn)
		bytes, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		c.Write(bytes)
	}
}
