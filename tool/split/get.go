package split

import (
	"os"
	"strconv"
	"strings"
)

//获取指定的文件分片

func GetFilePiece(path, filename, pieceRange string) (piece []byte, err error) {
	//打开文件
	file, err := os.Open(path + "/" + filename)
	if err != nil {
		return
	}
	defer file.Close()

	//解析range
	pieceStart, pieceSize, err := getPieceStartAndSize(pieceRange)
	if err != nil {
		return
	}

	//获取文件的指定片段
	filePiece := make([]byte, pieceSize)
	_, err = file.ReadAt(filePiece, int64(pieceStart))
	if err != nil {
		return
	}

	return filePiece, nil
}

// 解析range
func getPieceStartAndSize(pieceRange string) (pieceStart, pieceSize int, err error) {
	//range格式：bytes=0-1023
	//利用正则表达式获取起始位置和大小
	const prefix = "bytes="
	if !strings.HasPrefix(pieceRange, prefix) {
		return
	}
	rangeValue := strings.TrimPrefix(pieceRange, prefix)

	rangeParts := strings.Split(rangeValue, "-")

	//解析range头的值为整数
	pieceStart, err = strconv.Atoi(rangeParts[0])
	if err != nil {
		return
	}

	pieceEnd, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return
	}

	pieceSize = pieceEnd - pieceStart + 1
	return
}
