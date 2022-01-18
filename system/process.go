package system

import (
	"bytes"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type PrecessState struct {
	Pid     int64   `json:"pid"`
	CpuPerc float64 `json:"cpuPerc"`
	MemPerc float64 `json:"memPerc"`
	VSZ     int64   `json:"vsz"`
	RSS     int64   `json:"rss"`
	State   string  `json:"state"`
	StartAt string  `json:"startAt"`
	Time    string  `json:"time"`
	Cmd     string  `json:"cmd"`
}

func GetProcessInfo(server string) (PrecessState, error) {
	var out PrecessState

	cmd := `ps aux | ` +
		`grep ` + server + ` | ` +
		`grep -v grep | ` +
		`awk {'printf("\"pid\":%s,\"cpuPerc\":%s,\"memPerc\":%s,\"vsz\":%s,\"rss\":%s,\"state\":\"%s\",\"startAt\":\"%s\",\"time\":\"%s\",\"cmd\":\"%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s\"\n",` +
		`$2, $3, $4, $5, $6, $8, $9, $10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)'}`
	data, err := BashCommand(cmd)
	if err != nil {
		return out, err
	}

	data = bytes.Trim(data, "\n\r ")
	data = bytes.TrimSpace(data)
	list := bytes.Split(data, []byte("\n"))

	if len(list) != 1 {
		return out, fmt.Errorf("server cmd num=%d is too many", len(list))
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonData []byte
	jsonData = append(jsonData, '{')
	jsonData = append(jsonData, list[0]...)
	jsonData = append(jsonData, '}')
	err = json.Unmarshal(jsonData, &out)
	if err != nil {
		return out, err
	}

	out.Cmd = strings.Trim(out.Cmd, " ")

	return out, nil
}
