package tcp_net

import (
	"fmt"
	"testing"
)

func TestNewServer(t *testing.T) {
	server, _ := NewServer()
	fmt.Println(server)
}
