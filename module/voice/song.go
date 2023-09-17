package voice

import "GoATuber-2.0/engine"

//处理唱歌消息

func HandelSongVoice(e *engine.Engine, m engine.PriorityMessage) {
	e.Voice.VType = Song

	var (
		voice      engine.MessageSlice
		instrument engine.MessageSlice
	)

	voice.Index = 0
	voice.Song = m.Voice
	voice.Content = m.SongName //第一条消息的内容是歌曲名字

	instrument.Index = 1
	instrument.Song = m.Instrument

	e.Message.MessageSlice = append(e.Message.MessageSlice, voice)
	e.Message.MessageSlice = append(e.Message.MessageSlice, instrument)

	//将消息发送到前端
	e.Ch.ToFront <- struct{}{}
}
