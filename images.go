package main

import (
	//"mvdan.cc/xurls"
	"fmt"
	"net/http"
	"strings"
)

func isImage(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Err:", err)
	}
	defer resp.Body.Close()
	switch strings.Split(resp.Header.Get("Content-Type"), "/")[0] {
	case "image":
		return true
	default:
		return false
	}
}
