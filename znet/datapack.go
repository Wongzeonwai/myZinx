package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go-zinx/utils"
	"go-zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (p *DataPack) GetHeadLen() uint32 {
	// Datalen uint32（4字节） + ID uint32（4字节）
	return 8
}

func (p *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgID写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (p *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到datalen和msgID
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, errors.New("data len err:" + err.Error())
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, errors.New("msg id err:" + err.Error())
	}
	// 判断datalen是否超出maxpackagesize
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("datalen too large！")
	}
	return msg, nil
}
