package iface

//将消息封装到message中

type IMessage interface {
	//GetMsgId 获取消息的id号
	GetMsgId() uint32
	//GetDataLen 获取消息数据段的长度
	GetDataLen() uint32
	//GetData 获取消息数据的内容
	GetData() []byte

	//SetData 设置消息内容
	SetData([]byte)
	//SetMsgId 设置消息的ID
	SetMsgId(uint32)
	//SetDataLen 设置数据的长度
	SetDataLen(uint322 uint32)
}
