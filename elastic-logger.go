package logs

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type logParams struct {
	url string
	host string
}

func New(url string, host string) logParams {
	l := logParams{url: url, host: host}
	return l
}

func toElastic(url string, host string, level string, message string) {
	message = strings.Replace(message, "\"", "'", -1)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	if host == "" {
		host = "Error"
	}
	postingUrl := url + "/" + host + "/logs"

	unixNano := time.Now().UnixNano()
	umillisec := unixNano / 1000000
	n := int64(umillisec)
	moment := strconv.FormatInt(n, 10)
	json := "{ \"Level\": \"" + level + "\", \"message\" : \"" + message + "\", \"timestamp\":" + moment + ", \"timestamp2\": \"" + time.Now().Format(time.RFC3339) +"\"}"

	_, err := client.Post(postingUrl, "application/json" ,
		bytes.NewBuffer([]byte(json)))
	if err != nil {
		fmt.Printf("Failed to log error message. Post failed %s \n", postingUrl, err)
	}
}

func Info(url string, host string, message string) {
	if url != "" {
		toElastic(url, host, "Info", message)
	} else {
		fmt.Printf("Info : %s \n",  message)
	}
}

func Error(url string, host string, message string) {
	if url != "" {
		toElastic(url, host, "Error", message)
	} else {
		fmt.Printf("Error : %s \n",  message)
	}
}

func (l logParams) Info(message string, args ...interface{}) {
	if l.url != "" {
		toElastic(l.url, l.host, "Info", fmt.Sprintf(message, args...))
	} else {
		fmt.Printf("Info : %s \n",  message)
	}
}

func (l logParams) Error(message string, args ...interface{}) {
	if l.url != "" {
		toElastic(l.url, l.host, "Error", fmt.Sprintf(message, args...))
	} else {
		fmt.Printf("Error : %s \n",  message)
	}
}