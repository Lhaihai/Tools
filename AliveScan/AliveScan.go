package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

var (
	cips sync.Map
	filename = "alive.txt"
	thread int
	CIDR string
)

func main() {
	flag.IntVar(&thread,"t",512, "Threads")
	flag.StringVar(&CIDR,"host","","172")
	flag.Parse()
	t := strings.Split(CIDR,".")
	if CIDR == "192" {
		check192Alive(thread)
	}else if CIDR == "172" {
		check172Alive(thread)
	}else if CIDR == "10" {
		check10Alive(thread)
	}else if len(t)> 1 && strings.HasPrefix(CIDR,"10.") {
		checkBAlive(CIDR,thread)
	}else {
		usage()
	}

	//writeFile(time.Now().Format("2006-01-02 15:04:05")+"\n")
	//check192Alive(thread)
	//writeFile(time.Now().Format("2006-01-02 15:04:05")+"\n")
	//check172Alive(thread)
	//writeFile(time.Now().Format("2006-01-02 15:04:05"))
	//check10Alive(thread)
	//writeFile(time.Now().Format("2006-01-02 15:04:05"))
}

func usage() {
    msg := `
Usage:
    AliveScan host thread                 (默认使用512个并发协程) 
    Example: ./AliveScan -host 192 -t 50  用50个并发协程探测192.168 B段存活的C段
             ./AliveScan -host 172        探测172.16-32 B段存活的C段
             ./AliveScan -host 10         探测10 A段存活的C段
             ./AliveScan -host 10.172     探测10.172 B段存活的C段`

	fmt.Println(msg)
	os.Exit(0)
}

func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func sendICMPRequest(icmp ICMP, destAddr *net.IPAddr) error {
	conn, err := net.DialIP("ip4:icmp", nil, destAddr)
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	tStart := time.Now()

	conn.SetReadDeadline((time.Now().Add(time.Second * 2)))

	recv := make([]byte, 1024)
	receiveCnt, err := conn.Read(recv)

	if err != nil {
		return err
	}

	tEnd := time.Now()
	duration := tEnd.Sub(tStart).Nanoseconds() / 1e6

	fmt.Printf("%d bytes from %s: seq=%d time=%dms\n", receiveCnt, destAddr.String(), icmp.SequenceNum, duration)

	return err
}



func check192Alive(thread int){
	pools := New(thread)
	for i:=1;i<255;i++{
		host := "192.168."+strconv.Itoa(i)
		pools.Add(1)
		go func(host string){
			runICMP(host+".1")
			runICMP(host+".254")
			pools.Done()
		}(host)
	}
	pools.Wait()
}

func check172Alive(thread int){
	pools := New(thread)
	for j:=16 ; j < 33 ; j++{
		for i:=1;i<255;i++{
			host := "172."+strconv.Itoa(j)+"."+strconv.Itoa(i)
			pools.Add(1)
			go func(host string){
				runICMP(host+".1")
				runICMP(host+".254")
				pools.Done()
			}(host)
		}
	}
	pools.Wait()
}

func check10Alive(thread int){
	pools := New(thread)
	for j:=1 ; j < 255 ; j++{
		for i:=1;i<255;i++{
			host := "10."+strconv.Itoa(j)+"."+strconv.Itoa(i)
			pools.Add(1)
			go func(host string){
				runICMP(host+".1")
				runICMP(host+".254")
				pools.Done()
			}(host)
		}
	}
	pools.Wait()
}

func checkBAlive(t string,thread int){
	pools := New(thread)
	for i:=1;i<255;i++{
		host := t+"."+strconv.Itoa(i)
		pools.Add(1)
		go func(host string){
			runICMP(host+".1")
			runICMP(host+".254")
			pools.Done()
		}(host)
	}
	pools.Wait()
}

func runICMP(host string)  {
	raddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Printf("Fail to resolve %s, %s\n", host, err)
		return
	}

	//fmt.Printf("Ping %s (%s):\n\n", raddr.String(), host)

	for i:=1;i<2;i++{
		if err = sendICMPRequest(getICMP(uint16(i)), raddr); err != nil {
			fmt.Printf("Error: %s\n", err)
		}else {
			tmphost := strings.Split(host,".")
			cip := tmphost[0] + "." + tmphost[1] + "." + tmphost[2] + ".1/24"
			_ ,ok := cips.LoadOrStore(cip,1)
			if ok{

			} else {
				writeFile(cip+"\n")
			}
		}
	}

	//time.Sleep(1 * time.Second)
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
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

func writeFile(cip string){
	var f *os.File
	var err1 error
	if checkFileIsExist(filename){
		f,err1 = os.OpenFile(filename,os.O_APPEND,0666)
	}else {
		f,err1 = os.Create(filename)
	}
	defer f.Close()
	_, err1 = io.WriteString(f,cip)
	if err1 != nil {
		fmt.Println(err1)
	}
}