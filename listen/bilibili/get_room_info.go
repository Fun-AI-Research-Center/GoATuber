package bilibili

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andybalholm/brotli"
	"github.com/tidwall/gjson"
)

var (
	reqHeader = map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	}
)

// BiliClient 包含一个 HTTP 客户端，用于发送请求到 Bilibili API。
type BiliClient struct {
	Client *http.Client
}

// 向指定的 URL 发送一个 HTTP GET 请求
func (client *BiliClient) get(url string) (*gjson.Result, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//设置请求头
	for name, value := range reqHeader {
		req.Header.Add(name, value) //向请求头中添加一个新的字段
	}
	return handleResp(client.Client.Do(req))
}

// handleResp 函数用于处理 HTTP 响应。它接受一个指向 http.Response 结构体的指针和一个错误作为输入，并返回一个指向 gjson.Result 结构体的指针和一个错误。
func handleResp(resp *http.Response, err error) (*gjson.Result, error) {
	// 如果请求发生了错误，那么直接返回错误。
	if err != nil {
		return nil, err
	}
	// 在函数执行完毕后关闭响应体。
	defer func() {
		_ = resp.Body.Close()
	}()
	// 请求失败
	if resp.StatusCode != 200 {
		// 如果 HTTP 状态码不为 200，那么返回一个错误。
		return nil, errors.New(resp.Status)
	}

	// 解压缩响应体
	var reader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(resp.Body)
	case "br":
		reader = brotli.NewReader(resp.Body)
	default:
		// 如果响应头中没有指定压缩方式，则不进行解压缩。
		fmt.Printf(contentEncoding)
	}
	buf := &bytes.Buffer{}
	// 将响应体中的数据读入一个缓冲区中。
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	// 解析 JSON 响应体。
	r := gjson.ParseBytes(buf.Bytes())
	code := r.Get("code").Int()
	// 如果返回的 JSON 中包含错误码，则返回一个错误。
	if code != 0 {
		return nil, errors.New(r.Get("message").String())
	}
	// 否则，返回解析后的 JSON 数据。
	return &r, nil
}

// GetRoomInfo 获取直播间信息
func (client *BiliClient) GetRoomInfo(roomID int) (Room, error) {
	url := "https://api.live.bilibili.com/room/v1/Room/get_info?room_id=" + strconv.Itoa(roomID)
	resp, err := client.get(url)
	if err != nil {
		return Room{}, err
	}
	r := Room{
		RoomID: roomID,
	}
	//发起查询请求
	data := resp.Get("data")
	//获得真roomID
	r.RealRoomID = int(data.Get("room_id").Int())
	if err != nil {
		return Room{}, err
	}
	return r, nil
}
