package tcp_net

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"tcp-demo/iface"
	"tcp-demo/packet"
	"tcp-demo/property"
)

var (
	ConnectionIsClosed    = errors.New("连接已经关闭")
	PacketMsgError        = errors.New("封包失败")
	DataSendToClientError = errors.New("数据发送到客户端失败")
)

type Connection struct {
	TcpServer *Server //当前Connection 属于哪一个server
	TcpConn   *net.TCPConn
	ConnId    uint32
	isClosed  bool
	//退出通知
	ctx    context.Context
	cancel context.CancelFunc
	//该连接处理的方法
	MsgManage iface.IMessageManage
	//读写分离的chan
	msgChan chan []byte

	//连接的属性 eg:[name:zs] [form:127.0.0.1:9888]
	attribute *property.Property
}

//NewConnection connection 构造函数
func NewConnection(server *Server, tcpConn *net.TCPConn, connId uint32, router iface.IMessageManage) *Connection {
	connection := &Connection{
		TcpServer: server,
		TcpConn:   tcpConn,
		ConnId:    connId,
		isClosed:  false,
		MsgManage: router,
		msgChan:   make(chan []byte),
		attribute: property.PropertyInstance,
	}
	//添加连接
	server.ConnManage.AddConn(connection)
	return connection
}

func (c *Connection) StartConnection() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//启动读数据的功能
	go c.startReader()

	go c.startWrite()

	if c.TcpServer.ConnStartingFunc != nil {
		c.TcpServer.ConnStartingFunc(c)
	}

	//TODO 启动从当前连接写数据的功能

	//检测连接是否关闭
	for {
		select {
		case <-c.ctx.Done():
			c.StopConnection()
			return
		}
	}
}

func (c *Connection) StopConnection() {
	if !c.isClosed {
		//关闭连接之前
		if c.TcpServer.ConnStoppingFunc != nil {
			c.TcpServer.ConnStoppingFunc(c)
		}

		//关闭连接 回收资源
		c.TcpConn.Close()
		//从连接管理中 移除本连接
		c.TcpServer.ConnManage.RemoveConn(c)
		c.isClosed = true
	}
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.TcpConn
}

func (c *Connection) GetConnectionId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.TcpConn.RemoteAddr()
}

// SendMsg 直接将msg 数据发送到tcp 客户端
func (c *Connection) SendMsg(msgId uint32, bytes []byte) error {
	var (
		err  error
		pack []byte
	)
	if c.isClosed {
		return ConnectionIsClosed
	}

	//将数据封包 发送给客户端
	msg := packet.NewMessage(msgId, bytes)

	if pack, err = packet.GetPackInstance().Packet(msg); err != nil {
		return PacketMsgError
	}

	//写回客户端
	c.msgChan <- pack
	return nil
}

func (c *Connection) startReader() {
	fmt.Println("[READER GOROUTINE IS RUNNING...]")

	var (
		err  error
		msg  iface.IMessage
		data []byte
	)
	defer c.cancel()

	for {
		dp := packet.GetPackInstance()

		//读取客户端的msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Printf("readFull to headdata from conn err:%s\n", err.Error())
			return
		}

		//拆包
		if msg, err = dp.UnPack(headData); err != nil {
			fmt.Printf("unpack err:%s\n", err.Error())
			return
		}

		//根据dataLen 获取data 放入msg 中
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Printf("readFull to data from conn err:%s\n", err.Error())
				return
			}
		}

		msg.SetData(data)

		req := &Request{
			connection: c,
			msg:        msg,
		}

		c.TcpServer.Pool.Submit(func() {
			c.TcpServer.Wg.Add(1)
			c.MsgManage.DoMsgHandler(req)
			c.TcpServer.Wg.Done()
		})
	}
}

//写消息goroutine ，将数据发送给客户端
func (c *Connection) startWrite() {
	fmt.Println("[WRITER GOROUTINE IS RUNNING...]")
	var (
		err error
	)
	for {
		select {
		case data := <-c.msgChan:
			if _, err = c.TcpConn.Write(data); err != nil {
				fmt.Errorf("%s err:%s", DataSendToClientError, err.Error())
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//暴露给外部使用的接口

func (c *Connection) SetProperty(key string, value interface{}) {
	c.attribute.SetProperty(key, value)
}

func (c *Connection) GetProperty(key string) interface{} {
	return c.attribute.GetProperty(key)
}

func (c *Connection) RemoveProperty(key string) {
	c.attribute.RemoveProperty(key)
}
