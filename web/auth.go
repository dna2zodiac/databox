package web

import (
	"os"
	"time"
	"strings"
	"log"
	"net/http"
	"io/ioutil"
)

type ServerAuthBasic struct {
	FileName string
	Value string
	Mtime time.Time
	Watcher *time.Timer
}

var DefaultAuth = ServerAuthBasic{}

func watchAuthBasic(t *time.Timer, d time.Duration, a *ServerAuthBasic) {
	<- t.C
	if (a.checkModified()) {
		a.loadBasicAuth()
	}
	t.Reset(d)
	go watchAuthBasic(t, d, a)
}

func (a *ServerAuthBasic) CheckAuth(r *http.Request) bool {
	if a.FileName == "" {
		return true
	}
	if a.Value == "" {
		a.loadBasicAuth()
	}
	if a.Watcher == nil {
		d := time.Minute
		a.Watcher = time.NewTimer(d)
		go watchAuthBasic(a.Watcher, d, a)
	}

	value := strings.Trim(r.Header.Get("Authorization"), " \r\n\t")
	if value == a.Value {
		return true
	}
	return false
}

func (a *ServerAuthBasic) loadBasicAuth() {
	file, err := os.Open(a.FileName)
	if err != nil {
		log.Printf("failed to load basic auth: %v", err)
		return
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to load basic auth: %v", err)
		return
	}
	nextValue := strings.Trim(string(buf), " \r\n\t")
	if (a.Value == "" && nextValue == "" && a.Watcher == nil) || (a.Value != "" && nextValue == "") {
		log.Printf("set [empty] value to basic auth ...")
	}
	a.Value = nextValue

	stat, err := file.Stat()
	if err == nil {
		a.Mtime = stat.ModTime()
	}
}

func (a *ServerAuthBasic) checkModified() bool {
	if a.FileName == "" {
		return false
	}
	file, err := os.Open(a.FileName)
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	if a.Mtime.Equal(stat.ModTime()) {
		return false
	}
	return true
}
