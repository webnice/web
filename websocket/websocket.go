package websocket

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"

	"golang.org/x/net/websocket"
)

// New creates a new object and return interface
func New() Interface {
	var wst = new(impl)
	return wst
}

func echoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
