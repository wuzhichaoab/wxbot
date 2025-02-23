﻿package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Wxid               string `json:"wxid"`
	Content            string `json:"content"`
	ToUser             string `json:"toUser"`
	Msgid              uint64 `json:"msgid"`
	OriginMsg          string `json:"originMsg"`
	ChatRoomSourceWxid string `json:"chatRoomSourceWxid"`
	MsgSource          string `json:"msgSource"`
	Type               uint32 `json:"type"`
	DisplayMsg         string `json:"displayMsg"`
}

func wsClient() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)

		var msg Message
		json.Unmarshal(message, &msg)
		log.Printf("msg: %+v\n", msg)
	}
}

func httpServer() {
	g := gin.Default()
	g.POST("/callback", func(c *gin.Context) {
		var msg Message

		err := c.BindJSON(&msg)
		if err != nil {
			log.Printf("bind json faild: %s\n", err)
			return
		}

		log.Printf("msg: %+v\n", msg)
	})

	g.Run(*addr)
}

func escapeQuotes(s string) string {
	var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	return quoteEscaper.Replace(s)
}

func sendFormImg() {
	var client http.Client

	// 要上传的文件
	file, _ := os.Open(*img_path)
	defer file.Close()

	// 设置body数据并写入缓冲区
	bodyBuff := bytes.NewBufferString("") //bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)

	// math/rand
	// _ = bodyWriter.SetBoundary(fmt.Sprintf("-----------------------------%d", rand.Int()))
	// 加入图片二进制
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("image"), escapeQuotes(filepath.Base(file.Name()))))
	h.Set("Content-Type", "image/jpg")
	part, err := bodyWriter.CreatePart(h)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}
	// 其他字段
	err = bodyWriter.WriteField("wxid", *wxid)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	// 填充boundary结尾
	bodyWriter.Close()

	// 组合创建数据包
	req, err := http.NewRequest("POST", *addr+"/sendimgmsg", bodyBuff)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	req.ContentLength = int64(bodyBuff.Len())
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	fmt.Printf("form-img send msg response: %s\n", data)
}

func sendJsonImg() {
	var client http.Client

	type ImgInfo struct {
		Wxid  string `json:"wxid"`
		Image []byte `json:"image"`
	}

	data, err := ioutil.ReadFile(*img_path)
	if err != nil {
		log.Fatalf("read file faild: %s\n", err)
	}

	ii := ImgInfo{
		Wxid:  *wxid,
		Image: data,
	}

	j_data, err := json.Marshal(ii)
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}

	req, err := http.NewRequest("POST", *addr+"/sendimgmsg", bytes.NewReader(j_data))
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}
	// req.ContentLength = int64(bodyBuff.Len())
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	fmt.Printf("json-img send msg response: %s\n", data)
}

func sendFormFile() {
	var client http.Client

	// 要上传的文件
	file, _ := os.Open(*file_path)
	defer file.Close()

	// 设置body数据并写入缓冲区
	bodyBuff := bytes.NewBufferString("") //bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)

	// math/rand
	// _ = bodyWriter.SetBoundary(fmt.Sprintf("-----------------------------%d", rand.Int()))
	// 加入图片二进制
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes(filepath.Base(file.Name()))))
	h.Set("Content-Type", "text/plain")
	part, err := bodyWriter.CreatePart(h)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}
	// 其他字段
	err = bodyWriter.WriteField("wxid", *wxid)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	// 填充boundary结尾
	bodyWriter.Close()

	// 组合创建数据包
	req, err := http.NewRequest("POST", *addr+"/sendfilemsg", bodyBuff)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	req.ContentLength = int64(bodyBuff.Len())
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	fmt.Printf("form-file send msg response: %s\n", data)
}

func sendJsonFile() {
	var client http.Client

	type FileMsg struct {
		Wxid     string `json:"wxid"`
		File     []byte `json:"file"`
		FileName string `json:"filename"`
	}

	data, err := ioutil.ReadFile(*file_path)
	if err != nil {
		log.Fatalf("read file faild: %s\n", err)
	}

	fm := FileMsg{
		Wxid:     *wxid,
		File:     data,
		FileName: path.Base(*file_path),
	}

	j_data, err := json.Marshal(fm)
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}

	req, err := http.NewRequest("POST", *addr+"/sendfilemsg", bytes.NewReader(j_data))
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}
	// req.ContentLength = int64(bodyBuff.Len())
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("json Marshal faild: %s\n", err)
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("CreatePart faild: %s\n", err)
	}

	fmt.Printf("json-file send msg response: %s\n", data)
}

var addr = flag.String("addr", "localhost:8080", "Http service address")
var mode = flag.String("mode", "json-file", "Select the startup mode. The optional values are ws, http, form-img, json-img, form-file and json-file")
var img_path = flag.String("img", "../1.jpg", "Specify image path when sending image messages")
var file_path = flag.String("file", "../1.txt", "Send file message specifying file path")
var wxid = flag.String("wxid", "47331170911@chatroom", "Send message recipient's wxid")

func main() {
	flag.Parse()
	log.SetFlags(0)

	switch *mode {
	case "ws":
		wsClient()
	case "http":
		httpServer()
	case "form-img":
		sendFormImg()
	case "json-img":
		sendJsonImg()
	case "form-file":
		sendFormFile()
	case "json-file":
		sendJsonFile()
	}
}
