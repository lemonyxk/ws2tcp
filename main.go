/**
* @program: go
*
* @create: 2024-07-10 12:07
**/

package main

import (
	websocket2 "github.com/fasthttp/websocket"
	"io"
	"log"
	"net"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		Exit("args should like: ws2tcp wsAddr > tcpAddr or ws2tcp wsAddr < tcpAddr")
	}

	var args = os.Args[1]
	var arr = strings.Split(args, " ")
	if len(arr) != 3 {
		Exit("args should like: ws2tcp wsAddr > tcpAddr or ws2tcp wsAddr < tcpAddr")
	}

	if arr[1] != ">" && arr[1] != "<" {
		Exit("second args must be > or <")
	}

	if !strings.HasPrefix(arr[0], "ws://") &&
		!strings.HasPrefix(arr[0], "wss://") &&
		!strings.HasPrefix(arr[2], "ws://") &&
		!strings.HasPrefix(arr[2], "wss://") {
		Exit("can not find ws or wss")
	}

	check(arr)
}

func check(arr []string) {
	if strings.HasPrefix(arr[0], "ws://") || strings.HasPrefix(arr[0], "wss://") {
		if arr[1] == ">" {
			// ws to tcp
			ws2tcp(arr[0], arr[2])
		} else {
			// tcp to ws
			tcp2ws(arr[2], arr[0])
		}

	} else {
		if arr[1] == ">" {
			// tcp to ws
			tcp2ws(arr[0], arr[2])
		} else {
			// ws to tcp
			ws2tcp(arr[2], arr[0])
		}
	}
}

func ws2tcp(wsAddr, tcpAddr string) {

	var url, err = url2.Parse(wsAddr)
	if err != nil {
		Exit(err)
	}

	var subs = GetArgs("-sub", "--sub-protocols")
	var subProtocols = strings.Split(subs, ", ")

	// create a websocket server
	var httpServer = &http.Server{
		Addr: url.Host,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var upgrade = websocket2.Upgrader{
				HandshakeTimeout: time.Second * 6,
				ReadBufferSize:   4096,
				WriteBufferSize:  4096,
				Subprotocols:     subProtocols,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}

			conn, err := upgrade.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			defer func() { _ = conn.Close() }()

			log.Printf("ws %s to tcp %s\n", url.Host, tcpAddr)

			remoteConn, err := net.Dial("tcp", tcpAddr)
			if err != nil {
				log.Println(err)
				return
			}

			var closeCh = make(chan struct{}, 1)

			go func() {
				_, _ = io.Copy(remoteConn, conn.NetConn())
				_ = remoteConn.Close()
				_ = conn.Close()
			}()

			go func() {
				_, _ = io.Copy(conn.NetConn(), remoteConn)
				_ = conn.Close()
				_ = remoteConn.Close()
				closeCh <- struct{}{}
			}()

			<-closeCh

			log.Printf("ws %s to tcp %s closed\n", url.Host, tcpAddr)
		}),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		Exit(err)
	}
}

func tcp2ws(tcpAddr, wsAddr string) {

	listen, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		Exit(err)
	}

	defer func() { _ = listen.Close() }()

	var subs = GetArgs("-sub", "--sub-protocols")
	var subProtocols = strings.Split(subs, ", ")

	for {

		localConn, err := listen.Accept()
		if err != nil {
			break
		}

		var client = &websocket2.Dialer{
			HandshakeTimeout: time.Second * 6,
			ReadBufferSize:   4096,
			WriteBufferSize:  4096,
			Subprotocols:     subProtocols,
		}

		conn, _, err := client.Dial(wsAddr, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("tcp %s to ws %s\n", tcpAddr, wsAddr)

		go func() {

			var closeCh = make(chan struct{}, 1)

			go func() {
				_, _ = io.Copy(conn.NetConn(), localConn)
				_ = localConn.Close()
				_ = conn.Close()
			}()

			go func() {
				_, _ = io.Copy(localConn, conn.NetConn())
				_ = conn.Close()
				_ = localConn.Close()
				closeCh <- struct{}{}
			}()

			<-closeCh

			log.Printf("tcp %s to ws %s closed\n", tcpAddr, wsAddr)
		}()
	}
}
