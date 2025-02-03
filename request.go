package gosdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Bot 创建一个新的 Bot 实例
type Bot struct {
	// go-cqhttp 监听的 地址
	IP string
	// go-cqhttp 监听的 端口号
	Port uint16
	// go-cqhttp 统一资源定位符
	baseURL string

	// 在 config.yml 中的反向 HTTP 地址
	PostIP string
	// 在 config.yml 中的反向 HTTP 端口号
	PostPort uint16
	// 在 config.yml 中的反向 HTTP 统一资源定位符
	postURL string

	GroupMessageFuncPool []func(*Bot, *GroupMessage)
	// todo
	// PrivateMessageFuncPool []func(*PrivateMessage)
	// TempMessageFuncPool []func(*TempMessage)
}

// Run 创建新的 goroutine 并开始运行实例
func (b *Bot) Run() {
	b.baseURL = fmt.Sprintf("http://%s:%d", b.IP, b.Port)
	b.postURL = fmt.Sprintf("%s:%d", b.PostIP, b.PostPort)
	http.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		msgByte, err := io.ReadAll(r.Body)
		if err != nil {
			color.Red("读取响应流失败: ", err.Error())
		}
		msg, err := b.parseMsg(msgByte)
		if err != nil {
			color.Red("解析消息失败: %v", err.Error())
		}
		switch v := msg.(type) {
		case *GroupMessage:
			b.HandleGroupMessage(v)
		}
	})
	log.Fatal(http.ListenAndServe(b.postURL, nil))
}

// parseMsg 解析消息内容
func (b *Bot) parseMsg(message []byte) (interface{}, error) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, err
	}
	if msg.PostType == "message" || msg.PostType == "message_sent" {
		switch msg.MessageType {
		case "group":
			var groupMsg GroupMessage
			if err := json.Unmarshal(message, &groupMsg); err != nil {
				return nil, err
			}
			return &groupMsg, nil
		}
	}
	return nil, errors.New("unsupported event: " + msg.MessageType)
}

// parseResp 解析响应数据
func (b *Bot) parseResp(message string) (map[string]int64, error) {
	var msg BaseResp
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}
	if msg.Status == "ok" {
		var groupMsg OkayResp
		if err := json.Unmarshal([]byte(message), &groupMsg); err != nil {
			return nil, err
		}
		return groupMsg.Data, nil
	}
	return nil, errors.New(msg.Msg)
}

// post 向 go-cqhttp 发送 POST 请求
func (b *Bot) post(node string, data string) (*http.Response, error) {
	url := b.baseURL + node
	if data != "" {
		resp, err := http.Post(url, "application/json", strings.NewReader(data))
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		resp, err := http.Post(url, "", nil)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// SendGroupMessage 发送一条群消息, 无返回值
func (b *Bot) SendGroupMessage(group interface{}, msg []Element) {
	if gId, ok := group.(int64); ok && len(msg) != 0 {
		var buffer bytes.Buffer
		for _, ele := range msg {
			buffer.WriteString(ele.toJson())
			buffer.WriteRune(',')
		}
		s := buffer.String()[:buffer.Len()-1]
		message := fmt.Sprintf(
			`{"group_id": %d, "message": [%s]}`,
			gId, s)
		fmt.Println(message)
		resp, err := b.post("/send_group_msg", message)
		if err != nil {
			color.Red(err.Error())
			return
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Blue(err.Error())
			return
		}
		cont, err := b.parseResp(string(content))
		if err != nil {
			fmt.Println(string(content))
			return
		}
		if msgId := cont["message_id"]; msgId != 0 {
			output := func() string {
				var buffer bytes.Buffer
				for _, v := range msg {
					buffer.WriteString(v.String())
				}
				return buffer.String()
			}()
			color.New(color.FgGreen, color.Bold).Printf("[%s] [INFO]: #%d Group(%d) <- %s\n",
				time.Now().Format(time.DateTime),
				msgId, gId, output,
			)
		}
		//else {
		//	color.(`[!*] "%s"`, message)
		//}
	} else {
		color.Red("Group number is invalid.")
	}
}

// HandleGroupMessageFunc 处理消息的函数
func (b *Bot) HandleGroupMessageFunc(f func(*Bot, *GroupMessage)) {
	b.GroupMessageFuncPool = append(b.GroupMessageFuncPool, f)
}

// HandleGroupMessage 处理消息的函数，可自定义
func (b *Bot) HandleGroupMessage(message *GroupMessage) {
	fmt.Printf("[%s] [INFO]: #%d Group(%d) - 「%s」(%d) -> %s\n",
		time.Unix(message.Time, 0).Format(time.DateTime),
		message.MessageId,
		message.Group,
		message.Sender.Nickname,
		message.User,
		message.MsgChain)
	for _, f := range b.GroupMessageFuncPool {
		go f(b, message)
	}
}
