package dict

import (
	"errors"
	"strconv"
	"strings"

	"GoATuber-2.0/engine"
	"GoATuber-2.0/err"
)

type mood struct {
	Happy    int //开心的
	Health   int //健全的（皮套人直播术语，，，
	Mad      int //生气
	Sad      int //伤心的
	Fear     int //害怕
	Disgust  int //厌恶的
	Surprise int //惊讶
}

func GetMood(e *engine.Engine, msg *engine.MessageSlice) {
	search := New()
	er := search.LoadWordDict("config/cfg/dict/"+e.Config.Application.Dict.MoodDictName, 0)
	if er != nil {
		err.Error(errors.New("加载情感词汇本体失败："+er.Error()), err.Normal)
		return
	}

	moodWords := search.FindAll(msg.Content)
	var mood mood
	for _, moodWord := range moodWords {
		moodWord = search.SearchLine(moodWord)
		words := strings.Split(moodWord, ",")
		for i, word := range words {
			if i == 1 || i == 3 {
				num, er := strconv.Atoi(words[i+1]) //获得情感值
				if er != nil {
					err.Error(errors.New("mood,转化数字失败："+er.Error()), err.Normal)
				}
				if i == 3 {
					num = int(float32(num) * 0.5)
				}
				if word == "PA" || word == "PE" {
					mood.Happy += num
				} else if word == "PD" || word == "PH" || word == "PG" || word == "PB" || word == "PK" {
					mood.Health += num
				} else if word == "NA" {
					mood.Mad += num
				} else if word == "NB" || word == "NJ" || word == "NH" || word == "PF" {
					mood.Sad += num
				} else if word == "NI" || word == "NC" || word == "NG" {
					mood.Fear += num
				} else if word == "NE" || word == "ND" || word == "NN" || word == "NK" || word == "NL" {
					mood.Disgust += num
				} else if word == "PC" {
					mood.Surprise += num
				}
			}
		}
	}
	var moodNum [7]float64
	moodNum[0] = float64(mood.Happy)
	moodNum[1] = float64(mood.Mad)
	moodNum[2] = float64(mood.Sad)
	moodNum[3] = float64(mood.Disgust)
	moodNum[4] = float64(mood.Surprise)
	moodNum[5] = float64(mood.Fear)
	moodNum[6] = float64(mood.Health)
	msg.Emotion.Emo = max(moodNum)
}

func max(moodNum [7]float64) string {
	var max float64
	var maxFlag int
	for i := 0; i < 7; i++ {
		if moodNum[i] > max {
			max = moodNum[i]
			maxFlag = i
		}
	}
	if maxFlag == 0 {
		return "happy"
	} else if maxFlag == 1 {
		return "mad"
	} else if maxFlag == 2 {
		return "sad"
	} else if maxFlag == 3 {
		return "disgust"
	} else if maxFlag == 4 {
		return "surprise"
	} else if maxFlag == 5 {
		return "fear"
	} else if maxFlag == 6 {
		return "health"
	}
	return ""
}
