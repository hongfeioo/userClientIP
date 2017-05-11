package main

import (
	"fmt"
	"net/http"
)

func ip(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<table><tr><td>%s</td><td>%s</td></tr><tr><td>%s</td><td>%s</td></tr></table>", "本机IP", r.RemoteAddr, "浏览器", r.UserAgent())
	fmt.Fprintf(w, "%s<a href=\"getClientNatIPWindows.exe\">%s</a>,<a href=\"getClientNatIPLinux\">%s</a>,<a href=\"getClientNatIPMAC\">%s</a>", "获取自检工具请点击", "windows_X64", "linux_x64", "MacOS")
	fmt.Fprintf(w, "</p>%s", "注意1：windows平台需要关闭防火墙，并以管理员身份运行")
	fmt.Fprintf(w, "</p>%s", "注意2：MacOS平台需要以管理员身份运行")
}

func main() {

	fmt.Println("listen:8088")
	http.HandleFunc("/ip", ip)
	http.Handle("/", http.FileServer(http.Dir("./")))

	if err := http.ListenAndServe(":8088", nil); err != nil {
		fmt.Println(err)
	}
}
