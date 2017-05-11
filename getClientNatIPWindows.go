package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gooops/mtr"
)

const DEBUGFILENAME = "debug.txt"

func main() {

	//fmt.Println(getCurrentDirectory(), FILENAME)

	FILENAME := getCurrentDirectory() + "\\" + DEBUGFILENAME

	//------get local ip-------------
	fmt.Println("-------> local IP ----------")
	debugFile(FILENAME, fmt.Sprintln("-------> local IP ----------"))

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, addr := range addrs {
		fmt.Printf("%d:%s\n", i+1, addr.String())
		debugFile(FILENAME, fmt.Sprintf("%d:%s\n", i+1, addr.String()))

	}

	//----get external ip from 3 web station----------------------
	fmt.Println("------> get external ip ---------")
	debugFile(FILENAME, fmt.Sprintln("------> get external ip ---------"))
	ipcn, err := getipcn()
	if err != nil {
		fmt.Println(err.Error())
	}

	external, err := getexternalip()
	if err != nil {
		fmt.Println(err.Error())
	}

	ipipnet, err := getipipnet()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("\n1. ip.cn:%s\n\n2. myexternalip.com:%s\n\n3. ipip.net:%s\n\n", ipcn, external, ipipnet)
	debugFile(FILENAME, fmt.Sprintf("\n1. ip.cn:%s\r\n\n2. myexternalip.com:%s\r\n\n3. ipip.net:%s\r\n\n", ipcn, external, ipipnet))

	//---begin to traceroute , 第二个参数为最大跳数-------

	fmt.Println("------> traceroute---------")
	debugFile(FILENAME, fmt.Sprintln("------> traceroute---------"))

	traceliepin, err := gomtrWindows("www.liepin.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	debugFile(FILENAME, strings.Replace(traceliepin, "\n", "\r\n", -1))

	tracezhaopin, err := gomtrWindows("www.zhaopin.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	debugFile(FILENAME, strings.Replace(tracezhaopin, "\n", "\r\n", -1))

	//fmt.Println("successful")
	debugFile(FILENAME, fmt.Sprintln("------> successful---------"))
	fmt.Println("检测完毕，报告存放处：", FILENAME, "\n 按Crtl+c退出")

	c := make(chan string, 0)
	<-c

}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
	//return strings.Replace(dir, "\\", "/", -1)
}

func debugFile(fileName, strContent string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fdTime := time.Now().Format("2006-01-02 15:04:05")
	fdContent := strings.Join([]string{fdTime, "> ", strContent, "\r\n"}, "")
	buf := []byte(fdContent)
	fd.Write(buf)
	fd.Close()
}

func gomtrWindows(www string) (string, error) {

	resultString := ""

	mtr.LocalAddr()
	destAddrs, err := mtr.DestAddr(www)
	if err != nil {
		return resultString, err
	}

	fmt.Println(www, " -> ip:", convert(destAddrs[:]))

	// c := make(chan mtr.TracerouteHop, 0)
	// go func() {
	// 	for {
	// 		hop, ok := <-c
	// 		if !ok {
	// 			//fmt.Println()
	// 			return
	// 		}
	// 		fmt.Println("---", hop.TTL, hop.Address, hop.AvgTime, hop.BestTime, hop.Loss)
	// 		//debugFile(FILENAME, fmt.Sprintln(hop.TTL, hop.Address, hop.AvgTime, hop.BestTime, hop.Loss))
	//
	// 	}
	// }()
	//
	// options := mtr.TracerouteOptions{}
	// _, err1 := mtr.Mtr(destAddrs, &options, c)
	// if err1 != nil {
	// 	return resultString, err1
	// }

	mm, err := mtr.T(www, true, 20, 52, 5, 5) //参数分别为，主机名，最大跳数，发送字节数，snt ，retries
	if err != nil {
		return resultString, err
	}

	return mm, nil

}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ".")
}

//getexternalip  从国外获取本机的出口ip
func getexternalip() (string, error) {

	ipString := ""

	resp, err := http.Get("http://myexternalip.com/raw") //ip.cn
	if err != nil {
		return ipString, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ipString, err
	}
	return strings.TrimSpace(string(body)), nil

}

//getipcn  从ip.cn获取本机的出口ip     	<div id="result"><div class="well"><p>您现在的 IP：<code>124.65.168.26</code></p><p>所在地理位置：<code>北京市 联通</code></p><p>GeoIP: Beijing, China</p></div></div>
func getipcn() (string, error) {

	ipString := ""

	resp, err := http.Get("http://ip.cn") //ip.cn
	if err != nil {
		return ipString, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ipString, err
	}

	re, _ := regexp.Compile("<code>[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}</code>")

	IP := re.Find([]byte(body))

	re, _ = regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	IP = re.Find(IP)

	return string(IP), nil

}

//getipipnet    <li>124.65.168.26</li>
func getipipnet() (string, error) {

	ipString := ""

	resp, err := http.Get("http://www.ipip.net/share.html") //ip.cn
	if err != nil {
		return ipString, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ipString, err
	}

	re, _ := regexp.Compile("<li>[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}</li>")

	IP := re.Find([]byte(body))

	re, _ = regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	IP = re.Find(IP)

	return string(IP), nil

}

// //gettaobaoip  从淘宝获取ip
// func gettaobaoip() (string, error) {
//
// 	ipString := ""
//
// 	resp, err := http.Get("http://ip.taobao.com") //ip.cn
// 	if err != nil {
// 		return ipString, err
// 	}
//
// 	defer resp.Body.Close()
//
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return ipString, err
// 	}
// 	return string(body), nil
//
// }
//
// //
