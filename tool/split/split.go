package split

import (
	"bytes"
	"encoding/gob"
	"os"
	"strconv"
)

//文件分割服务

type Pieces []piece

type piece struct {
	Index int    `json:"index"`
	Range string `json:"range"`
}

// GetPieceData 将文件分片，将分片数据存储备用。
// filePath:文件路径。只到文件目录，不包括文件名。目录后不加/。
// fileName:文件名。包括后缀。
func GetPieceData(filePath, fileName string, pieceSize int) error {
	//打开文件
	file, er := os.Open(filePath + "/" + fileName)
	defer file.Close()

	if er != nil {
		return er
	}

	//获取文件分片信息
	info, er := file.Stat()
	if er != nil {
		return er
	}
	pieceNum := chunkFileNum(int(info.Size()), pieceSize)

	//开始获取分片数据
	var pieces Pieces
	for i := 0; i < pieceNum; i++ {
		var (
			piece         piece
			pieceStart    int //分片起始位置
			tempPieceSize int //分片大小
		)
		piece.Index = i

		//计算分片起始位置
		pieceStart = i * pieceSize
		//计算分片大小
		if i == pieceNum-1 {
			tempPieceSize = int(info.Size()) - pieceStart
		} else {
			tempPieceSize = pieceSize
		}

		//获得分片范围
		piece.Range = "bytes=" + strconv.Itoa(pieceStart) + "-" + strconv.Itoa(pieceStart+tempPieceSize-1)
		pieces = append(pieces, piece)
	}

	//将分片数据利用gob编码后存储
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if er := enc.Encode(pieces); er != nil {
		return er
	}
	file, er = os.Create(filePath + "/" + fileName + ".data")
	if er != nil {
		return er
	}
	defer file.Close()

	_, er = file.Write(buf.Bytes())
	if er != nil {
		return er
	}

	return nil
}

// 确定文件分片的数目
func chunkFileNum(fileSize int, pieceSize int) int {
	chunks := fileSize / pieceSize
	if fileSize%pieceSize != 0 {
		chunks++
	}
	return chunks
}

// GetPieceInfo 获取分片数据
func GetPieceInfo(filePath, fileName string) (Pieces, error) {
	//打开文件
	file, er := os.Open(filePath + "/" + fileName + ".data")
	defer file.Close()
	if er != nil {
		return nil, er
	}

	//解码
	var pieces Pieces
	dec := gob.NewDecoder(file)
	if er := dec.Decode(&pieces); er != nil {
		return nil, er
	}

	return pieces, nil
}
