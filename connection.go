package neo4jgo

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
)

type neo4jInfo struct {
	url  string
	auth string
}

type Executer interface {
	Execute(string, map[string]interface{}) string
}

func NewExecuter(host string, user string, password string) Executer {
	return &neo4jInfo{"http://" + host + "/db/data/transaction/commit", auth(user, password)}
}

func valueToString(value interface{}) string {
	switch t := value.(type) {
	case string:
		return "\"" + t + "\""
	case int:
		return strconv.FormatInt(int64(t), 10)
	case []int:
		res := ""
		for _, v := range t {
			if res != "" {
				res = res + ", "
			}
			res = res + valueToString(v)
		}
		return "[" + res + "]"
	case []string:
		res := ""
		for _, v := range t {
			if res != "" {
				res = res + ", "
			}
			res = res + valueToString(v)
		}
		return "[" + res + "]"
	}
	return "err"
}

func paramsToJson(params map[string]interface{}) string {
	res := ""
	for key, value := range params {
		if res != "" {
			res = res + ", "
		}
		res = res + "\"" + key + "\": " + valueToString(value)
	}
	return "{" + res + "}"
}

func createBody(stmt string, params map[string]interface{}) string {
	c := stmt
	p := paramsToJson(params)
	res := `{"statements": [{"statement": "` + c + `", "parameters": ` + p + `}]}`
	return res
}

func (con *neo4jInfo) Execute(stmt string, params map[string]interface{}) string {
	client := new(http.Client)
	body := createBody(stmt, params)
	req, err := http.NewRequest("POST", con.url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+con.auth)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}

func auth(user string, password string) string {
	e := base64.StdEncoding
	return e.EncodeToString([]byte(user + ":" + password))
}
