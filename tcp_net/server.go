package tcp_net

import (
	"fmt"
	"log"
	"net"
	"sync"
	"tcp-demo/config"
	"tcp-demo/conn_management"
	"tcp-demo/coroutinepool"
	"tcp-demo/iface"
	"tcp-demo/msg_manage"
)

//Server 服务器实例
type Server struct {
	Name    string
	IpVer   string
	Ip      string
	Port    int
	MaxConn int

	MsgManage  iface.IMessageManage
	Pool       *coroutinepool.Pool
	Wg         sync.WaitGroup
	ConnManage iface.IConnectionManage //server的连接管理器

	//hook 函数
	ConnStartingFunc func(connection iface.IConnection)
	ConnStoppingFunc func(connection iface.IConnection)
}

func NewServer() (iface.IServer, error) {
	//初始化协程池
	pool, err := coroutinepool.NewPool(config.ServerCon.MaxPoolCapacitySize, config.ServerCon.LimitTask)
	if err != nil {
		log.Fatalln(err)
	}
	return &Server{
		Name:       config.ServerCon.Name,
		IpVer:      config.ServerCon.IpVer,
		Ip:         config.ServerCon.Ip,
		Port:       config.ServerCon.Port,
		MaxConn:    config.ServerCon.MaxConn,
		MsgManage:  msg_manage.MsgManageInstance,
		Pool:       pool,
		ConnManage: conn_management.ConnectionManageClient,
	}, nil
}

func (s *Server) Start() {
	fmt.Printf("[START TCP SERVER] Server Name [%s] Server Listener At Ip:[%s:%d]\n", s.Name, s.Ip, s.Port)
	var (
		err         error
		addr        *net.TCPAddr
		tcpListener *net.TCPListener
		tcpConn     *net.TCPConn
	)

	if addr, err = net.ResolveTCPAddr(s.IpVer, fmt.Sprintf("%s:%d", s.Ip, s.Port)); err != nil {
		log.Fatalln(err)
	}

	//listen TCP
	if tcpListener, err = net.ListenTCP(s.IpVer, addr); err != nil {
		log.Fatalln(err)
	}
	var connId uint32 = 0

	for {
		//accept
		if tcpConn, err = tcpListener.AcceptTCP(); err != nil {
			fmt.Printf("Accept err:%s", err.Error())
			continue
		}
		//判断连接是否超过最大连接数量
		if s.ConnManage.Count() > s.MaxConn {
			//取消这个连接  进行下一次的accept
			tcpConn.Close()
			continue
		}

		//实例化一个连接
		dealConn := NewConnection(s, tcpConn, connId, s.MsgManage)
		connId++

		go dealConn.StartConnection()
	}
}

func (s *Server) Stop() {
	fmt.Printf("[STOP TCP SERVER] [NAME:%s]", s.Name)
	//等待所有业务协程运行结束
	s.Wg.Wait()
	//关闭协程池
	s.Pool.CloseAndRelease()
	//清除server 中保存的所有连接
	s.ConnManage.EliminateConn()
}

func (s *Server) Server() {
	go s.Start()
	//预留之后 可以在服务器开启后做一些其他的业务操作
	select {}
}

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.MsgManage.AddRouter(msgId, router)
}

//SetConnStarting 设置该Server的连接开始前的Hook函数
func (s *Server) SetConnStarting(function func(iface.IConnection)) {
	s.ConnStartingFunc = function
}

//SetConnStopping 设置该Server的连接断开时的Hook函数
func (s *Server) SetConnStopping(function func(iface.IConnection)) {
	s.ConnStoppingFunc = function
}
