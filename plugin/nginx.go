package plugin

import (
	"fmt"
	"strings"
)

type Nginx struct {
	Name string
}

func (w *Nginx) Matches(text string) bool {
	return strings.HasPrefix(text, "nginx")
}

func (w *Nginx) Respond(msg *Message) error {
	tokens := strings.Fields(msg.Text)
	if len(tokens) > 1 {
		switch tokens[1] {
		case "log":
			if len(tokens) > 2 {
				var filename string
				switch tokens[2] {
				case "access":
					filename = "access.log"
				case "error":
					filename = "error.log"
				}
				out, err := ExecShell("tail", "/var/log/nginx/" + filename)
				if err != nil {
					return fmt.Errorf("%s遇到了点麻烦，正在紧张处理中...", w.Name)
				}
				msg.Send(out)
			}
		}
	}
	msg.Done()

	return nil
}

func (w *Nginx) Help() string {
	return "nginx - 管理 Nginx.\n nginx log access\nnginx log error"
}

