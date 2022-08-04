package packet

//Message 消息结构体
type Message struct {
	id      uint32 //消息id
	dataLen uint32 //数据长度
	data    []byte //数据本身
}

//NewMessage  消息构造函数
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetData(bytes []byte) {
	m.data = bytes
}

func (m Message) SetMsgId(id uint32) {
	m.id = id
}

func (m Message) SetDataLen(len uint32) {
	m.dataLen = len
}
