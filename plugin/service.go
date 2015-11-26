package plugin

import (
	"fmt"
	"strings"
)

type Service struct {
	Name string
}

func (w *Service) Matches(text string) bool {
	return strings.HasPrefix(text, "service")
}

func (w *Service) Respond(msg *Message) error {
	tokens := strings.Fields(msg.Text)
	if len(tokens) > 1 {
		switch tokens[1] {
		case "list":
			out, err := ExecShell("consul", "members")
			if err != nil {
				return fmt.Errorf("%s遇到了点麻烦，正在紧张处理中...", w.Name)
			}
			msg.Send(out)
		}
	}
	msg.Done()
	return nil
}

func (w *Service) Help() string {
	return "service - 管理 Consul 服务."
}

