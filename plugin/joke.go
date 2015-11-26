package plugin

import (
	"time"
	"fmt"
	"math/rand"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const MAX_COUNT = 999

type Joke struct{}

func (w *Joke) Matches(text string) bool {
	return strings.HasPrefix(text, "joke") || strings.Contains(text, "笑话")
}

func (w *Joke) Respond(msg *Message) error {
	resp, err := http.Get(fmt.Sprintf("http://m2.qiushibaike.com/article/list/suggest?count=%d&page=1", MAX_COUNT))
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	count, _ := data["count"].(float64)
	items, _ := data["items"].([]interface{})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(int(count))
	item, _ := items[index].(map[string]interface{})
	joke, _ := item["content"].(string)

	msg.Send(joke)
	msg.Done()
	return nil
}

func (w *Joke) Help() string {
	return "joke - 讲笑话."
}
