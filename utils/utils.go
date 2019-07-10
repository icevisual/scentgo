package utils

import (
	"bufio"
	"encoding/hex"
	"io"
	"os/exec"
	"strings"
)

func Hex2bytes(s string) []byte{
	a,_ := hex.DecodeString(s)
	return a
}

func Bytes2hex(s []byte) string{
	return hex.EncodeToString(s)
}

func GetSerialPorts() []string {
	command := "/bin/bash"
	params := []string{"-c", "ls /dev/ttyUSB*"}

	r,_ := ExecCommand(command, params)
	// fmt.Println("After execCommand",r,err,len(r))
	return r
}

func ExecCommand(commandName string, params []string) ([]string,error) {
	var contentArray = make([]string, 0, 5)
	// contentArray = contentArray[0:0]
	cmd := exec.Command(commandName, params...)
	// Show CMD
	// fmt.Printf("Run CMD: %s\n", strings.Join(cmd.Args[1:], " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return contentArray,err
	}
	err = cmd.Start()
	if err != nil {
		return contentArray,err
	}
	// Start开始执行c包含的命令，但并不会等待该命令完成即返回。
	// Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。

	reader := bufio.NewReader(stdout)
	var index int
	// read msg from stdout
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		line = strings.Trim(line,"\n")
		// fmt.Print(line)
		index++
		contentArray = append(contentArray, line)
	}
	err = cmd.Wait()
	return contentArray,err
}

