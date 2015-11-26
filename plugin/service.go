package plugin

import (
//"fmt"
	"strings"
)

type Service struct{}

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
				return err
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

