package iface

//数据封包和拆包，添加消息头 为了解决tcp 粘包问题

type IDataPack interface {
	//GetHeadLen 获取数据包的头长度的方法
	GetHeadLen() uint32
	//Packet 封包方法
	Packet(msg IMessage) ([]byte, error)
	//UnPack 拆包方法
	UnPack([]byte) (IMessage, error)
}
