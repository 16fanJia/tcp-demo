package main

import (
	"fmt"
	"io"
	"net"
	"tcp-demo/packet"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")

	conn, err := net.Dial("tcp", "localhost:2222")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	var count int

	for {
		//发封包message消息
		dp := packet.GetPackInstance()
		message := packet.NewMessage(1, []byte("我是你爹 测试数据。。。。。。"))

		//fmt.Println("======debug========", message.GetMsgId(), message.GetDataLen(), string(message.GetData()))

		msg, _ := dp.Packet(message)
		_, err := conn.Write(msg)

		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*packet.Message)
			data := make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			msg.SetData(data)

			fmt.Println("==> Recv Msg: ID=", msg.GetMsgId(), ", len=", msg.GetDataLen(), ", data=", string(msg.GetData()))
		}

		time.Sleep(1 * time.Second)
		count++
		if count == 10 {
			conn.Close()
		}
	}
}
