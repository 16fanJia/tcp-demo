package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"tcp-demo/config"
	"tcp-demo/iface"
)

const (
	//DefaultHeaderLen 默认头部长度 ID uint32 (4 byte)+ DataLen uint32(4 byte)
	DefaultHeaderLen uint32 = 8
)

var (
	dataLenIsTooLarge = errors.New("数据包长度超出允许的最大值")
	dataPackInstance  iface.IDataPack
)

//DataPack 封包 拆包 结构体
type DataPack struct{}

func init() {
	dataPackInstance = &DataPack{}
}

//GetPackInstance 返回包操作 全局唯一实例 饿汉式
func GetPackInstance() iface.IDataPack {
	return dataPackInstance
}

func (d *DataPack) GetHeadLen() uint32 {
	return DefaultHeaderLen
}

//Packet 封包
func (d *DataPack) Packet(msg iface.IMessage) ([]byte, error) {
	var err error
	//创建一个存放bytes字节的缓冲区
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen 写入dataBuff中
	if err = binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//将msgId 写入dataBuff中
	if err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data 写入dataBuff中
	if err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//UnPack 拆包
func (d *DataPack) UnPack(binaryData []byte) (iface.IMessage, error) {
	var err error

	//创建一个从输入二进制数据的ioReader
	reader := bytes.NewReader(binaryData)

	msg := &Message{}

	//读dataLen
	if err = binary.Read(reader, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}

	//判断dataLen 是否超出了我们允许的最大长度
	if config.ServerCon.MaxPacketSize > 0 && msg.dataLen > uint32(config.ServerCon.MaxPacketSize) {
		return nil, dataLenIsTooLarge
	}
	//读msgID
	if err = binary.Read(reader, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	//只需要读出head 的数据包
	return msg, nil
}
