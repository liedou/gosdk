package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Element interface {
	String() string
	toJson() string
}

func NewText(txt string) Text {
	return Text{
		Text: txt,
	}
}

type Text struct {
	Element
	Text string `json:"text"`
}

func (t Text) String() string {
	return t.Text
}

func (t Text) toJson() string {
	return fmt.Sprintf(
		`{"type": "text",
			 "data": {
				 "text": "%s"
			 }
			}`, ConvertEscapedChars(t.Text))
}

func NewFace(id string) Face {
	return Face{
		Id: id,
	}
}

type Face struct {
	Element
	Id string `json:"id"`
}

func (f Face) String() string {
	if f.Id != "" {
		return fmt.Sprintf(`[Face:%s]`, f.Id)
	}
	return "[Face]"
}

func (f Face) toJson() string {
	return fmt.Sprintf(
		`{"type": "face",
			 "data": {
				 "id": "%s"
			 }
			}`, f.Id)
}

func NewRecord(file string) Record {
	return Record{
		File: file,
	}
}

type Record struct {
	Element
	File string `json:"file"`
}

func (r Record) String() string {
	if r.File != "" {
		return fmt.Sprintf("[Record | %s]", r.File)
	}
	return "[Record]"
}

func (r Record) toJson() string {
	return fmt.Sprintf(
		`{"type": "record",
			 "data": {
				 "file": "%s"
			 }
			}`, r.File)
}

func NewVideo(file, url string) Video {
	return Video{
		File: file,
		Url:  url,
	}
}

type Video struct {
	Element
	File string `json:"file"`
	Url  string `json:"url"`
}

func (v Video) String() string {
	if v.Url != "" {
		return fmt.Sprintf("[Video | %s]", v.Url)
	} else if v.File != "" {
		return fmt.Sprintf("[Video | %s]", v.File)
	}
	return "[Video]"
}

func (v Video) toJson() string {
	return fmt.Sprintf(
		`{"type": "video",
			 "data": {
				 "file": "%s"
			 }
			}`, v.File)
}

func NewAt(id string) At {
	return At{
		QQ: id,
	}
}

type At struct {
	Element
	QQ string `json:"qq"`
}

func (a At) String() string {
	return fmt.Sprintf(`[At:%s]`, a.QQ)
}

func (a At) toJson() string {
	return fmt.Sprintf(
		`{"type": "at",
			 "data": {
				 "qq": "%s"
			 }
			}`, a.QQ)
}

func NewShare(url, title string) Share {
	return Share{
		Url:   url,
		Title: title,
	}
}

type Share struct {
	Element
	Url   string `json:"url"`
	Title string `json:"title"`
}

func (s Share) String() string {
	if s.Url != "" && s.Title != "" {
		return fmt.Sprintf(
			"[Share: %s | %s]",
			s.Title, s.Url,
		)
	} else if s.Url != "" && s.Title == "" {
		return fmt.Sprintf("[Share | %s]", s.Url)
	} else if s.Url == "" && s.Title != "" {
		return fmt.Sprintf("[Share: %s]", s.Title)
	} else {
		return "[Share]"
	}
}

func (s Share) toJson() string {
	return fmt.Sprintf(
		`{"type": "share",
			 "data": {
				 "url": "%s",
				 "title": "%s"
			 }
			}`, s.Url, s.Title)
}

func NewImage(file, url string) Image {
	return Image{
		File: file,
		Url:  url,
	}
}

type Image struct {
	Element
	File string `json:"file"`
	Url  string `json:"url"`
}

func (i Image) String() string {
	if i.Url != "" {
		return fmt.Sprintf("[Image | %s]", i.Url)
	} else if i.File != "" {
		return fmt.Sprintf("[Image | %s]", i.File)
	}
	return "[Image]"
}

func (i Image) toJson() string {
	return fmt.Sprintf(
		`{"type": "image",
			 "data": {
				 "file": "%s"
			 }
			}`, i.File)
}

func NewReply(id string) Reply {
	return Reply{
		Id: id,
	}
}

type Reply struct {
	Element
	Id string `json:"id"`
}

func (r Reply) String() string {
	if r.Id != "" {
		return fmt.Sprintf("[Reply: %s]", r.Id)
	}
	return "[Reply]"
}

func (r Reply) toJson() string {
	return fmt.Sprintf(
		`{"type": "reply",
			 "data": {
				 "id": "%s"
			 }
			}`, r.Id)
}

func NewPoke(id int64) Poke {
	return Poke{
		QQ: id,
	}
}

type Poke struct {
	Element
	QQ int64 `json:"qq"`
}

func (p Poke) String() string {
	return fmt.Sprintf(
		`[戳:%d]`,
		p.QQ)
}

func (p Poke) toJson() string {
	return fmt.Sprintf(
		`{"type": "poke",
			 "data": {
				 "qq": "%d"
			 }
			}`, p.QQ)
}

func NewTTS(text string) TTS {
	return TTS{
		Text: text,
	}
}

type TTS struct {
	Element
	Text string `json:"text"`
}

func (t TTS) String() string {
	return fmt.Sprintf(
		`[TTS:%s]`,
		t.Text)
}

func (t TTS) toJson() string {
	return fmt.Sprintf(
		`{"type": "tts",
			 "data": {
				 "text": "%s"
			 }
			}`, t.Text)
}

type MessageData struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (md MessageData) toJson() string {
	var buffer bytes.Buffer
	for k, v := range md.Data {
		buffer.WriteString(fmt.Sprintf(`"%s": "%s",`, k, v))
	}
	s := buffer.String()
	if s != "" {
		s = s[:len(s)-1]
	}
	return fmt.Sprintf(`{%s}`, s)
}

type MessageChain []MessageData

func (mc MessageChain) Plain() string {
	var buffer bytes.Buffer
	for _, data := range mc {
		if data.Type == "text" {
			buffer.WriteString(data.Data["text"])
		}
	}
	return buffer.String()
}

func (mc MessageChain) fromJson() []Element {
	chain := make([]Element, 0, len(mc))
	for _, ele := range mc {
		switch ele.Type {
		case "text":
			var e Text
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "at":
			var e At
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "image":
			var e Image
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "face":
			var e Face
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "record":
			var e Record
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "video":
			var e Video
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "share":
			var e Share
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "reply":
			var e Reply
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		case "tts":
			var e TTS
			if err := json.Unmarshal([]byte(ele.toJson()), &e); err == nil {
				chain = append(chain, e)
			}
		}
	}
	return chain
}

func (mc MessageChain) String() string {
	var buffer bytes.Buffer
	for _, v := range mc.fromJson() {
		buffer.WriteString(fmt.Sprint(v))
	}
	return buffer.String()
}

type Message struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
}

type Sender struct {
	Age      uint8  `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserId   int64  `json:"user_id"`
}

type PrivateMessage struct {
	Message
	MessageId int32        `json:"message_id"`
	SubType   string       `json:"sub_type"`
	TargetId  int64        `json:"target_id"`
	MsgChain  MessageChain `json:"message"`
	RawMsg    string       `json:"raw_msg"`
	Sender    Sender       `json:"sender"`
}

type GroupMessage struct {
	Message
	MessageId int32        `json:"message_id"`
	User      int64        `json:"user_id"`
	Group     int64        `json:"group_id"`
	RawMsg    string       `json:"raw_message"`
	MsgChain  MessageChain `json:"message"`
	Sender    Sender       `json:"sender"`
}

type BaseResp struct {
	Status  string `json:"status"`
	RetCode int    `json:"retcode"`
	Msg     string `json:"msg"`
}

type OkayResp struct {
	BaseResp
	Data map[string]int64 `json:"data"`
}
