package log

import (
	"balance/registry"
	"bytes"
	"fmt"
	stlog "log"
	"net/http"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{
		url: serviceURL,
	})
}

type clientLogger struct {
	url string
}

// 重写io writer方法，使log输出到系统的log服务中
func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message. Service responded with %v", res.StatusCode)
	}
	return len(data), nil
}
