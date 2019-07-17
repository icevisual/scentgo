package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"runtime"
	"runtime/pprof"
	"scentrealm_bcc/datapkg"
	"scentrealm_bcc/datapkg/defines/CTLBCC"
	"scentrealm_bcc/datapkg/defines/CTLSP"
	"scentrealm_bcc/scentser"
	"scentrealm_bcc/utils"
	"sync"
	"time"
)

var logger = utils.Logger

var (
	ErrControllerFound        = errors.New("Controller Not Found ")
	ErrControllerNotConnected = errors.New("Controller Not Connected ")
	ErrMissingParameters      = errors.New("Missing Parameters ")
)

type ScentBccOperator struct {
	sp     *SerialPort
	pkgBcc datapkg.DataPackageCTLBcc
}

func (sbo *ScentBccOperator) IsConnected() bool {

	if sbo.sp == nil {
		return false
	}

	return sbo.sp.isConnected
}


func (sbo *ScentBccOperator) HandleDisconnect(pas map[string]interface{}) error {
	if !sbo.IsConnected() {
		return nil
	}

	if sbo.sp != nil {
		err := sbo.sp.s.Close()
		sbo.sp.isConnected = false
		return err
	}

	return nil
}


func (sbo *ScentBccOperator) HandleConnect(pas map[string]interface{}) error {
	if sbo.IsConnected() {
		return nil
	}
	s := FindConnectedCTL()
	if s == nil {
		// Notify Controller Not Found
		logger.Println("Controller Not Found ")
		return ErrControllerFound
	}
	sbo.sp = s
	return nil
}

func (sbo *ScentBccOperator) HandlePlaySmell(pas map[string]interface{}) error {
Retry:
	if !sbo.IsConnected() {
		err := sbo.HandleConnect(nil)
		if err != nil {
			return ErrControllerNotConnected
		}
	}

	if pas == nil {
		return ErrMissingParameters
	}
	err := utils.SimpleValidate(pas, map[string]string{
		"smell":    "required|type:uint32|range:1,20",
		"duration": "required|type:uint32|range:500,200000",
		"channel":  "sometimes|default:0|type:byte",
	})
	if err != nil {
		return err
	}

	dt := sbo.pkgBcc.AssemblePlaySmell(pas["smell"].(uint32), pas["duration"].(uint32), pas["channel"].(byte))
	err = sbo.sp.Send(dt)
	if err != nil {
		goto Retry
	}
	return err
}

func (sbo *ScentBccOperator) HandleStopPlay(pas map[string]interface{}) error {
Retry:
	if !sbo.IsConnected() {
		err := sbo.HandleConnect(nil)
		if err != nil {
			return ErrControllerNotConnected
		}
	}
	var channel byte = 0
	if pas != nil {
		err := utils.SimpleValidate(pas, map[string]string{
			"channel": "sometimes|default:0|type:byte",
		})
		if err != nil {
			return err
		}
		channel = pas["channel"].(byte)
	}

	dt := sbo.pkgBcc.AssembleStopPlay(channel)
	err := sbo.sp.Send(dt)

	if err != nil {
		goto Retry
	}

	return err
}

func (sbo *ScentBccOperator) HandleWakeUp(pas map[string]interface{}) error {
Retry:
	if !sbo.IsConnected() {
		err := sbo.HandleConnect(nil)
		if err != nil {
			return ErrControllerNotConnected
		}
	}
	var isBlocking = false
	if pas != nil {
		err := utils.SimpleValidate(pas, map[string]string{
			"blocking": "sometimes|default:false|type:bool",
		})
		if err != nil {
			return err
		}
		isBlocking = pas["blocking"].(bool)
	}
	if !isBlocking {
		go func(){
			err := sbo.sp.WakeUp()
			if err != nil {

			}
		}()
	} else {
		err := sbo.sp.WakeUp()
		if err != nil {
			goto Retry
		}
		return err
	}


	return nil
}
var addr = flag.String("addr", "localhost:8383", "http service address")

func MainRunner() {
	flag.Parse()
	var server = &scentser.ScentWsServer{}
	var handler = &ScentBccOperator{}
	server.Init()
	server.ServerAddr = *addr

	server.AddHandler(scentser.CmdConnect, handler.HandleConnect)
	server.AddHandler(scentser.CmdDisconnect, handler.HandleDisconnect)
	server.AddHandler(scentser.CmdPlaySmell, handler.HandlePlaySmell)
	server.AddHandler(scentser.CmdStopPlay, handler.HandleStopPlay)
	server.AddHandler(scentser.CmdWakeup, handler.HandleWakeUp)

	fmt.Println(server)
	server.RunServer()
}

func main() {
	//MainRunner()
	p := FindConnectedCTL()
	if p ==nil{
		logger.Println("Nt fD")
	}
	fmt.Println(pprof.Profiles())
}

type SerialPort struct {
	s           *serial.Port
	c           *serial.Config
	wl          sync.Mutex
	isConnected bool
}

func (sp *SerialPort) WakeUp() error {
	sp.wl.Lock()
	defer sp.wl.Unlock()
	Interval := 150
	TotalSecond := 2000
	Ins := []byte{0xf5, 0x51, 0xff, 0x03, 0xf5, 0x71, 0x55, 0x03, 0x0e, 0x55}
	for i := 0; i < TotalSecond; i += Interval {
		n, err := sp.s.Write(Ins)
		if err != nil {
			sp.isConnected = false
			logger.Println("Write Error = ", n, err)
			return err
		}
		time.Sleep(time.Millisecond * time.Duration(Interval))
	}
	time.Sleep(time.Millisecond * 300)
	return nil
}

func (sp *SerialPort) Reconnect() error {
	sp.wl.Lock()
	defer sp.wl.Unlock()

	return nil
}

func (sp *SerialPort) Send(data []byte) error {
	sp.wl.Lock()
	defer sp.wl.Unlock()
	n, err := sp.s.Write(data)
	if err != nil {
		sp.isConnected = false
		logger.Println("Write Error = ", n, err)
	}

	return err
}

func (sp *SerialPort) SendWait(data []byte, cmd byte, dpi datapkg.DataPackageInter) {
	sp.wl.Lock()
	defer sp.wl.Unlock()
	n, err := sp.s.Write(data)
	if err != nil {
		logger.Println("Write err", err)
		return
	}
	buf := make([]byte, 256)
	var res bytes.Buffer
	for {
		n, err = sp.s.Read(buf)
		if n == 0 {
			break
		}
		res.Write(buf[:n])
		if err != nil {
			sp.isConnected = false
			// io.EOF == err2
			return
		}
	}
	if len(res.Bytes()) > 0 {
		array := res.Bytes()
		r := dpi.VerifyCMD(cmd, array, byte(len(array)))
		if r {
			// Matches
		} else {
			// Not match
		}
	} else {
		// time out
	}
}

func GetSerialPortList() []string {
	switch runtime.GOOS {
	case "darwin":
		return utils.GetSerialPorts()
	case "windows":
		return []string{"COM3", "COM13"}
	case "linux":
		return utils.GetSerialPorts()
	default:
		panic("Unknown OS ")
	}
	return nil
}

func FindConnectedCTL() *SerialPort {
	sps := GetSerialPortList()
	logger.Println("SerialPorts ", sps)
	ch := make(chan *SerialPort,1)

	for _, v := range sps {
		go TestSerialPortAsController(v, ch)
	}
	logger.Println("After Go ")
	select {
	case s := <-ch:
		logger.Printf("MATCH %+v\n", s.c)
		return s
	case <-time.After(time.Second * 22):
		logger.Println("FindConnectedCTL Timeout")
	}
	return nil
}

func TestSerialPortAsController(portname string, ch chan *SerialPort) {
	logger.Println("TestSerialPortAsController")
	c := &serial.Config{
		Name:        portname,
		Baud:        115200,
		ReadTimeout: time.Millisecond * 400,
	}
	// fmt.Printf("%+v %s\n",c,reflect.TypeOf(c))
	s, err := serial.OpenPort(c)
	if err != nil {
		logger.Println("OpenPort err", portname, err)
		return
	}
	logger.Println("OpenPort")
	defer func() {

		if err != nil {
			logger.Println("defer close", portname)
			err = s.Close()
		}
	}()

	n, err := s.Write(utils.Hex2bytes("f527002755"))
	logger.Println("WriteHex2bytes")
	if err != nil {
		logger.Println("Write err", err)
		return
	}
	buf := make([]byte, 256)
	var res bytes.Buffer
	for {
		n, err = s.Read(buf)
		if n == 0 {
			break
		}

		res.Write(buf[:n])
		if err != nil {
			// io.EOF == err2
			return
		}
		if n == 6 {
			break
		}
	}
	if len(res.Bytes()) > 0 {
		array := res.Bytes()
		pkg := &datapkg.DataPackageCTLSp{}
		pkg.Init()
		r := pkg.VerifyCMD(CTLSP.RespGetControllerPhysicalChannel, array, byte(len(array)))
		if r {
			ch <- &SerialPort{s, c, sync.Mutex{}, true}

		} else {
			err = io.EOF
		}
	} else {
		err = io.EOF
	}
}

func TestConnectCTLAndPlaySmell() {

	pkgBcc := &datapkg.DataPackageCTLBcc{}
	pkgBcc.Init()

	s := FindConnectedCTL()
	if s == nil {
		// Notify Controller Not Found
		logger.Println("Controller Not Found ")
		return
	}

	dt := pkgBcc.AssemblePlaySmell(1, 20000, 0)
	utils.PrintByteArray("AssemblePlaySmell", dt)
	_ = s.Send(dt)
	time.Sleep(time.Second)

	go func() {
		s.SendWait(pkgBcc.AssembleQueryDeviceDevice(2014, 0), CTLBCC.RespQueryDeviceState, pkgBcc)
		pkgBcc.Describe()
	}()

	go func() {
		time.Sleep(time.Second * 2)
		dt = pkgBcc.AssembleStopPlay(0)
		utils.PrintByteArray("AssembleStopPlay", dt)
		s.Send(dt)
	}()

	go func() {
		time.Sleep(time.Second * 4)
		s.SendWait(pkgBcc.AssembleQueryDeviceDevice(2014, 0), CTLBCC.RespQueryDeviceState, pkgBcc)
		pkgBcc.Describe()
	}()

	<-time.After(time.Second * 8)
}
