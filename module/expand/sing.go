package expand

import (
	"bufio"
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"sync"

	"GoATuber-2.0/engine"
)

//唱歌 本项目目前不执行训练任务。唱歌需要自行进行训练与推理。
//以后要实现自动化的推理似乎也不是不行。
//那么音频文件需要格式化命名才行。
//统一存储在根目录下的songs目录里。
//TODO:整个语音传输模式都需要进行更改。

var (
	songList = make(map[string]string) //歌曲列表。歌曲以字典映射形式存储在本地的songs/songs.txt中。格式为歌曲名-编号。
	songNum  int                       //歌曲最大编号。像的索引那样只增不减。
	songMu   sync.Mutex                //song的互斥锁。保护关键信息的读写。
)

func GetSongList() (map[string]string, error) {
	//加载歌曲列表文件
	file, er := os.OpenFile("songs/songs.txt", os.O_RDONLY, 0666)
	if er != nil {
		return nil, er
	}
	//读取列表文本内容，加载字典映射。
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//按歌曲名-编号格式分割
		song := strings.Split(line, "-")
		songList[song[0]] = song[1]
		//更新最大编号
		num, er := strconv.Atoi(song[1])
		if er != nil {
			return nil, er
		}
		if num > songNum {
			songNum = num
		}
	}
	return songList, nil
}

// UploadSong 获取上传的歌曲。依次传入人声、伴奏、歌曲名
func UploadSong(voice, instrument []byte, songName string) error {
	//歌曲数增加，互斥锁进行保护
	songMu.Lock()
	songNum++
	songList[songName] = strconv.Itoa(songNum)
	//写入歌曲文件
	voiceFilePath := "songs/voice/" + strconv.Itoa(songNum) + ".wav"
	er := os.WriteFile(voiceFilePath, voice, 0644)
	if er != nil {
		return er
	}
	//写入伴奏文件
	if instrument != nil { //可以不加伴奏
		instrumentFilePath := "songs/instrument/" + strconv.Itoa(songNum) + ".wav"
		er = os.WriteFile(instrumentFilePath, instrument, 0644)
		if er != nil {
			return er
		}
	}
	//更新歌曲列表文件
	file, er := os.OpenFile("songs/songs.txt", os.O_WRONLY|os.O_APPEND, 0666)
	if er != nil {
		return er
	}
	_, er = file.WriteString(songName + "-" + strconv.Itoa(songNum) + "\n")

	//释放互斥锁
	songMu.Unlock()
	return er
}

func DeleteSong(songName string) error {
	//在map内删除歌曲
	delete(songList, songName)
	//在文件内删除歌曲
	file, er := os.OpenFile("songs/songs.txt", os.O_WRONLY|os.O_TRUNC, 0666)
	if er != nil {
		return er
	}

	//重新写入歌曲列表文件
	songMu.Lock()
	for k, v := range songList {
		_, er = file.WriteString(k + "-" + v + "\n")
		if er != nil {
			return er
		}
	}
	songMu.Unlock()

	return nil
}

// Sing 加载歌曲音频文件，传递给前端进行播放
func Sing(e *engine.Engine, singName string) error {
	//获取歌曲编号
	songNum := songList[singName]
	//读取歌曲文件
	voiceFilePath := "songs/voice/" + songNum + ".wav"
	voice, er := os.ReadFile(voiceFilePath)
	if er != nil {
		return er
	}
	//读取伴奏文件
	instrumentFilePath := "songs/instrument/" + songNum + ".wav"
	instrument, er := os.ReadFile(instrumentFilePath)
	if er != nil {
		return er
	}

	var message engine.PriorityMessage
	message.MessageType = engine.SongMessage
	message.Voice = base64.StdEncoding.EncodeToString(voice)
	message.Instrument = base64.StdEncoding.EncodeToString(instrument)
	message.UUID = "42"
	message.Username = "Administrator"
	message.Priority = engine.MaxPriority //设置优先级为最大值

	//将message加入优先队列
	e.Ch.ExpendToQueue <- message
	return nil
}
