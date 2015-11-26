package plugin

import (
	"fmt"
	"strings"
)

type Hello struct{}

func (w *Hello) Matches(text string) bool {
	return strings.HasPrefix(text, "hello") || strings.Contains(text, "你好")
}

func (w *Hello) Respond(msg *Message) error {
	msg.Send(fmt.Sprintf("你好，我是机器人 %s，有什么需要我帮助的吗？", msg.BotName))
	msg.Done()
	return nil
}

func (w *Hello) Help() string {
	return "hello - 打招呼."
}
