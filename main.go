package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"reflect"
	"scentrealm_bcc/datapkg"
	"scentrealm_bcc/utils"
	"time"
)

func main() {
	sps := utils.GetSerialPorts()
	fmt.Println("SerialPorts ",sps)

	ch := make(chan *serial.Port)
	for _, v := range sps  {
		go AutoConnectCTL(v,ch)
	}
	s := <- ch
	fmt.Printf("%+v %s\n",s,reflect.TypeOf(s))
}

func AutoConnectCTL(portname string,ch chan *serial.Port){

	c := &serial.Config{
		Name: portname,
		Baud: 115200,
		ReadTimeout:time.Millisecond * 600,
	}
	// fmt.Printf("%+v %s\n",c,reflect.TypeOf(c))
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("OpenPort err",err)
		return
	}

	defer func () {
		fmt.Println("defer",portname)
		if err != nil{
			err = s.Close()
			log.Fatal(err)
		}
	}()

	n, err := s.Write(datapkg.AssembleGetControllerPC())
	fmt.Println(portname,"Write = ",n,err)
	if err != nil {
		fmt.Println("Write err",err)
		return
	}
	buf := make([]byte, 256)
	n, err = s.Read(buf)
	fmt.Println(portname,"Read = ",n,err)
	if err != nil {

		// io.EOF == err2
		fmt.Println("Read err",err)
		return
	}
	fmt.Println(utils.Bytes2hex(buf[:n]))
	ch <- s
}