package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var datas sync.Map
var iplist []string

func main() {
	h:=flag.String("h","", "host")
	p:=flag.String("p","80", "Enter the Ports!!")
	o:=flag.String("o","url.txt", "-o url.txt")
	f:=flag.String("f","", "-f domain or ip or cip list")

	flag.Parse()
	host := *h
	listPort := toPorts(*p)
	filename := *f


	if host != "" {
		iplist = toIPs(host)
	}else if filename != ""{
		iplist = readIPFile(filename)
	}else {
		flag.Usage()
		os.Exit(3)
	}

	run(iplist,listPort,*o)
	//writeData(*o)
}

func readIPFile(filename string) []string {

	var urllist []string
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return []string{}
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a),"/"){
			urllist = append(urllist,toIPs(string(a))...)
		}else {
			urllist = append(urllist,strings.TrimSpace(string(a)))
		}
	}
	return urllist
}

func run (ip1s []string,portlist []string,filename string){

	//pools := New(10)
	for _, ip := range ip1s[1:] {
		for _, port := range portlist{
			//fmt.Println(port)
			//pools.Add(1)

			//go func(ip string,port string) {
				var url string
				if port == "443" {
					url = "https://"+ip+":"+port
				}else {
					url = "http://"+ip+":"+port
				}
				WriteFile(filename,url)
				//_ ,ok := datas.LoadOrStore(url,url)
				//if ok{
				//
				//} else {
				//
				//}
				//datas [url] = ""
				//pools.Done()
			//}(ip,port)
		}
	}
	//pools.Wait()
}

func WriteFile(filename , str string){
	var f *os.File
	var err1 error
	f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666) //打开文件
	defer f.Close()

	_ , err1 = f.WriteString(str+"\n")
	if err1 != nil {
		log.Println(err1)
	}
}

func writeData(filename string){

	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		//fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		//fmt.Println("文件不存在")
	}
	defer f.Close()


	fs := func(_, str interface{}) bool {

		//这个函数的入参、出参的类型都已经固定，不能修改
		//可以在函数体内编写自己的代码，调用map中的k,v
		_, err1 = f.WriteString(str.(string)) //写入文件(字符串)
		if err1 != nil {
			log.Println(err1)
		}
		return true
	}

	datas.Range(fs)

	f.Sync()
}

//--------- Port scope ------------
/*
80
80,8080-9000
8000-9000
*/
func toPorts(ports string)[]string{
	var listPort []string
	if strings.Contains(ports,",") && strings.Contains(ports,"-"){
		tmplistport := strings.Split(ports,",")
		for _,port:= range tmplistport{
			if strings.Contains(port,"-"){
				start, _ := strconv.Atoi(strings.Split(port,"-")[0])
				end,_ := strconv.Atoi(strings.Split(port,"-")[1])

				for i:=start ; i<=end ; i++{
					listPort = append(listPort,strconv.Itoa(i))
				}
			}else {
				listPort = append(listPort,port)
			}
		}
	}else if strings.Contains(ports,"-"){
		start, _ := strconv.Atoi(strings.Split(ports,"-")[0])
		end,_ := strconv.Atoi(strings.Split(ports,"-")[1])

		for i:=start ; i<=end ; i++{
			listPort = append(listPort,strconv.Itoa(i))
		}
	}else if strings.Contains(ports,","){
		listPort = strings.Split(ports,",")
	}else {
		listPort = append(listPort,ports)
	}
	return listPort
}
//--------- Port scope ------------

func toIPs(host string) []string {

	var list []string

	if ( strings.Contains(host,"/") && len(strings.Split(host,".")) == 4 ) {
		list = Iplist(host)
	}else if strings.Contains(host,"-"){

		addrs := strings.Split(host,"-")
		addr := addrs[0]
		s := string([]byte(addr)[:strings.LastIndex(addr,".")])
		start, _ := strconv.Atoi(strings.Split(addr,".")[3])
		end,_ := strconv.Atoi(addrs[1])

		for i:=start ; i < end ; i++ {
			list = append(list, s+"."+strconv.Itoa(i))
		}

	}else {
		list = append(list,host)
	}
	return list
}

func incIP(ip net.IP) {

	for j := len(ip) - 1; j > 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Iplist(cidr string) []string {
	var list []string
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		ip = net.ParseIP(cidr)
		list = append(list, ip.String())
		return list
	}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {

		list = append(list, ip.String())
	}
	return list
}

type pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

func New(size int) *pool {
	if size <= 0 {
		size = 1
	}
	return &pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

func (p *pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

func (p *pool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *pool) Wait() {
	p.wg.Wait()
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}