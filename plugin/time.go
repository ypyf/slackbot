package plugin

import (
	"fmt"
	"time"
	"strings"
)

var weekname = []string{"星期日", "星期一", "星期二", "星期四", "星期四", "星期五", "星期六"}
type Time struct{}

func (w *Time) Matches(text string) bool {
	return strings.HasPrefix(text, "time") || strings.Contains(text, "时间")
}

func (w *Time) Respond(msg *Message) error {
	timestamp := time.Now().Format("2006年01月02日 15:04:05")
	datetime := strings.Fields(timestamp)
	msg.Send(fmt.Sprintf("亲，今天是%s %s，当前时间是%s",
		datetime[0], weekname[time.Now().Weekday()], datetime[1]))
	msg.Done()
	return nil
}

func (w *Time) Help() string {
	return "time - 返回服务器当前时间."
}
