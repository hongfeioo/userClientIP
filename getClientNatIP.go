//使用方法：  sudo  ./getClientNatIPMAC

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aeden/traceroute" // 这个包用的时候注意自己的电脑使用一种网络，同时插着有线和无线会有问题。
)

const FILENAME = "debug.txt"

func main() {

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

	if err := trac("www.liepin.com", 6, 20); err != nil { //参数2是探针的个数，参数3是最大探测的跳数
		fmt.Println(err.Error())
	}

	if err := trac("www.zhaopin.com", 6, 20); err != nil { //参数2是探针的个数，参数3是最大探测的跳数
		fmt.Println(err.Error())
	}
	// c := make(chan string, 0)
	// <-c
	fmt.Println("------> end---------")
	debugFile(FILENAME, fmt.Sprintln("------> end---------"))

}

//tracefile 写入文件
func debugFile(fileName, strContent string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fdTime := time.Now().Format("2006-01-02 15:04:05")
	fdContent := strings.Join([]string{fdTime, "> ", strContent}, "")
	buf := []byte(fdContent)
	fd.Write(buf)
	fd.Close()
}

//20170511  改造这个函数将日志写入文件
func printHop(hop traceroute.TracerouteHop) {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	hostOrAddr := addr
	if hop.Host != "" {
		hostOrAddr = hop.Host
	}
	if hop.Success {
		fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime)
		debugFile(FILENAME, fmt.Sprintf("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime))

	} else {

		fmt.Printf("%-3d *\n", hop.TTL)
		debugFile(FILENAME, fmt.Sprintf("%-3d *\n", hop.TTL))

	}

}

func address(address [4]byte) string {
	return fmt.Sprintf("%v.%v.%v.%v", address[0], address[1], address[2], address[3])
}

//trac  适用于mac下的trace
func trac(www string, probeNum, hops int) error {

	options := traceroute.TracerouteOptions{}
	options.SetRetries(probeNum - 1) //Set the number of probes per "ttl" to nqueries (default is 1 probe).
	options.SetMaxHops(hops)         //Set the max time-to-live (max number of hops) used in outgoing probe packets (default is 64)

	ipAddr, err := net.ResolveIPAddr("ip", www)
	if err != nil {
		return err
	}

	fmt.Printf("traceroute to %v (%v), %v hops max, %v byte packets\n", www, ipAddr, options.MaxHops(), options.PacketSize())
	debugFile(FILENAME, fmt.Sprintf("traceroute to %v (%v), %v hops max, %v byte packets\n", www, ipAddr, options.MaxHops(), options.PacketSize()))

	c := make(chan traceroute.TracerouteHop, 0)
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				//fmt.Println("cddd")
				return
			}
			//fmt.Println(hop)
			printHop(hop)
		}
	}()

	_, err = traceroute.Traceroute(www, &options, c)
	if err != nil {
		return err
	}

	return nil
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
