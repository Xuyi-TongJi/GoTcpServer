package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"server/iface"
	"server/utils"
)

type DataPack struct{}

func (d *DataPack) GetHeadLen() uint32 {
	// data len uint32 4bytes
	// data id uint32 4bytes
	return 8
}

// Pack 将封装好的Message以TLV格式输出为字节流
// 使用Buffer读
func (d *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	// write data len 4bytes
	err := binary.Write(buffer, binary.LittleEndian, msg.GetLen())
	if err != nil {
		return nil, err
	}
	// write id 4bytes
	err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	// write data
	err = binary.Write(buffer, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Unpack 读入并封装Message Head (len + id)
func (d *DataPack) Unpack(data []byte) (iface.IMessage, error) {
	msg := &Message{}
	// read data by reader
	reader := bytes.NewReader(data)
	/* 最后一个参数是数据类型， &msg.Len uint32 代表读4个字节 并且读入到msg.Len赋值 */
	if err := binary.Read(reader, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	// dataLen > MaxPackaging size
	if msg.Len > utils.GlobalObj.MaxPackingSize {
		err := errors.New(fmt.Sprintf("[UnPackaging Message Refused] Data is larger than max size, length = %d", msg.Len))
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	return msg, nil
}
