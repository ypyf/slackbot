package plugin

import (
	"fmt"
	"strings"
)

var allows = []string{
	"man", "info", "history",
	"source", "export", "env", "uname", "hostname",
	"ls", "dir", "cd", "pwd", "tree",
	"mv", "cp", "scp", "touch", "mkdir",
	"ps", "pgrep", "crontab", "watch", "nohup", "jobs", "bg", "fg",
	"cat", "sort", "uniq", "seq", "wc", "echo", "split",
	"date", "cal",
	"diff", "od", "readelf", "file", "nm",
	"tar", "gzip", "gunzip", "bzip2", "bunzip2", "zip", "unzip",
	"df", "free", "du", "mount",
	"grep", "find", "ack", "ag", "locate",
	"head", "tail",
	"less", "more",
	"w", "who", "which", "where",
	"ping", "netstat", "lsof", "ifconfig", "dig", "tcpdump",
	"md5", "md5sum",
	"git", "curl", "wget",
	"docker", "consul",
}


type Shell struct{}

func (w *Shell) Matches(text string) bool {
	return strings.HasPrefix(text, ":")
}

func (w *Shell) Respond(msg *Message) error {
	cmd := strings.TrimPrefix(msg.Text, ":")
	if len(cmd) < 1 {
		// do nothing
		return nil
	}
	tokens := strings.Fields(cmd)
	bin := tokens[0]
	args := tokens[1:]
//	if bin == "top" {
//		args = append(args, "-n 1")
//	}
	if !contains(allows, bin) {
		return fmt.Errorf("这个命令太危险，我还是回火星吧！\n")
	}
	out, err := ExecShell(bin, args...)
	if err != nil {
		return err
	}
	msg.Send(out)
	msg.Done()
	return nil
}

func (w *Shell) Help() string {
	return ": - 执行 Shell 命令."
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
