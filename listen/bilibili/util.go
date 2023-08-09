package bilibili

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/andybalholm/brotli"
)

// 数据解包，分离出头部的信息和数据体，lens为数据包长度，ver为数据类型，op为操作码
func unpackPacket(data []byte) (op int, body []byte) {
	head := data[:16]
	op = int(binary.BigEndian.Uint32(head[8:12]))
	body = data[16:]
	return
}

// 数据装包
func pack(ver, op int, data []byte) []byte {
	totalLen := 16 + len(data)
	head := make([]byte, 0, totalLen)
	head = putInt(head, int32(totalLen), 4)     //整个数据包长度，4个字节
	head = append(head, 0x00, 0x10)             //数据包头部长度，固定为16，2个字节
	head = putInt(head, int32(ver), 2)          //数据类型，2个字节
	head = putInt(head, int32(op), 4)           //操作码，4个字节
	head = append(head, 0x00, 0x00, 0x00, 0x01) //作用未知，固定值，4个字节
	return append(head, data...)
}

// 将 val 的 size 个字节写入到 data中，写入方式为大端序
// 例如：putInt(data, 0x0000f0f0, 2) 表示，将 0xf0f0 写入 data
func putInt(data []byte, val int32, size int) []byte {
	for ; size > 0; size-- {
		data = append(data, byte((val>>((size-1)*8))&0xff))
	}
	return data
}

// 拆包，将原始的数据内容解压、拆分
func unpack(data []byte) [][]byte {
	head, body := data[:16], data[16:]
	ver := int(binary.BigEndian.Uint16(head[6:8]))
	packets := make([][]byte, 0)
	switch ver {
	case 0, 1:
		packets = append(packets, data)
	case 2:
		buf := &bytes.Buffer{}
		reader, _ := zlib.NewReader(bytes.NewReader(body))
		_, _ = io.Copy(buf, reader)
		packets = append(packets, splitPackets(buf.Bytes())...)
	case 3:
		buf := &bytes.Buffer{}
		reader := brotli.NewReader(bytes.NewReader(body))
		_, _ = io.Copy(buf, reader)
		packets = append(packets, splitPackets(buf.Bytes())...)
	}
	return packets
}

// 数据拆包，被压缩的数据内容会有一个或多个数据包，将其拆分出来
func splitPackets(data []byte) [][]byte {
	packets := make([][]byte, 0)
	for size := uint32(0); len(data) != 0; {
		size = binary.BigEndian.Uint32(data[:4])
		packets = append(packets, data[:size])
		data = data[size:]
	}
	return packets
}
