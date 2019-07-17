package utils

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Logger = GetLogger()

func GetLogger()  *log.Logger {
	var logger = log.New(os.Stdout, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return logger
}

// Get Time Of Day like [09:15:11.1346699]
func GetTimeDT() string {
	segs := strings.Split(time.Now().String(), " ")
	return fmt.Sprintf("[%s]",segs[1])
}


func String2AnyType(strVal string, types string) (interface{}, error) {

	switch types {
	case "string":
	case "int":
		iv, err := strconv.ParseInt(strVal, 10, 32)
		// parsing \"12.3\" to int: invalid syntax
		return int(iv), err
	case "uint32":
		iv, err := strconv.ParseUint(strVal, 10, 32)
		return uint32(iv), err
	case "byte":
		iv, err := strconv.ParseUint(strVal, 10, 8)
		return byte(iv), err
	case "bool":
		iv, err := strconv.ParseBool(strVal)
		return iv, err
	default:
		panic("Unhandled type convert " + types)
	}
	return nil, nil
}

func AnyType2String(val interface{}) string {
	// {"smell":"123"} string
	// {"smell":123} float64
	// {"smell":123.22} float64
	strVal := ""
	switch vv := val.(type) {
	case string:
		strVal = vv
	case bool:
		strVal = strconv.FormatBool(vv)
	case int:
		strVal = strconv.Itoa(vv)
	case uint32:
		strVal = fmt.Sprintf("%d", vv)
	case byte:
		strVal = fmt.Sprintf("%d", vv)
	case float64:
		strVal = fmt.Sprintf("%f", vv)
		strVal = strings.Trim(strVal, "0")
		strVal = strings.Trim(strVal, ".")
		if strVal == "" {
			strVal = "0"
		}
	default:
		panic(fmt.Sprintf("AnyType2String Unhandled type %T", val))
	}
	return strVal
}
var 	ErrValidateSetting        = errors.New("Validate Setting Error ")

/**
err := SimpleValidate(pas, map[string]string{
	"smell":    "required|type:uint32|range:1,20",
	"duration": "required|type:uint32|range:500,200000",
	"channel":  "sometimes|default:0|type:byte",
})
if err != nil {
	return err
}
dt := sbo.pkgBcc.AssemblePlaySmell(pas["smell"].(uint32), pas["duration"].(uint32), pas["channel"].(byte))
*/
func SimpleValidate(pam map[string]interface{}, rules map[string]string) error {
	// required|type:int
	// sometimes|type:int|

	for k, v := range rules {
		if len(v) > 0 {
			vs := strings.ToLower(v)
			ruleSeg := strings.Split(vs, "|")
			for _, ruleSetting := range ruleSeg {
				ruleComp := strings.Split(ruleSetting, ":")
				var ruleParams []string
				if len(ruleComp) > 1 {
					ruleParams = strings.Split(ruleComp[1], ",")
				}
				switch ruleComp[0] {
				case "required":
					if _, ok := pam[k]; !ok {
						return errors.New(fmt.Sprintf("Field `%s` Is Required", k))
					}
				case "sometimes":
					if _, ok := pam[k]; !ok {
						break
					}
				case "type":
					if ruleParams == nil {
						return ErrValidateSetting
					}
					strVal := AnyType2String(pam[k])
					newVal, err := String2AnyType(strVal, ruleParams[0])
					if err != nil {
						return errors.New(fmt.Sprintf("Field `%s` Syntax Error", k))
					}
					pam[k] = newVal
				case "default":
					if ruleParams == nil {
						return ErrValidateSetting
					}
					if _, ok := pam[k]; !ok {
						pam[k] = ruleParams[0]
					}
				case "range":
					if ruleParams == nil || len(ruleParams) != 2 {
						return ErrValidateSetting
					}

					minVal, err :=	strconv.ParseInt(ruleParams[0],10,64)
					if err != nil{
						return err
					}
					maxVal, err :=	strconv.ParseInt(ruleParams[1],10,64)
					if err != nil{
						return err
					}

					var sVal int64 = 0
					switch pam[k].(type) {
					case int:
						sVal = int64(pam[k].(int))
					case int64:
						sVal = pam[k].(int64)
					case int32:
						sVal = int64(pam[k].(int32))
					case uint:
						sVal = int64(pam[k].(uint))
					case uint16:
						sVal = int64(pam[k].(uint16))
					case uint32:
						sVal = int64(pam[k].(uint32))
					case uint64:
						sVal = int64(pam[k].(uint64))
					default:
						return errors.New(fmt.Sprintf("Type `%s` Not Support Range ",reflect.TypeOf(pam[k])))
					}
					if sVal < minVal || sVal> maxVal {
						return errors.New(fmt.Sprintf("Field `%s` Out Of Range ",k))
					}
				default:
					panic("What ?" + ruleComp[0])
				}
			}
		}
	}
	return nil
}

func Hex2bytes(s string) []byte {
	a, _ := hex.DecodeString(s)
	return a
}

func Bytes2hex(s []byte) string {
	return hex.EncodeToString(s)
}
func Hex2bytesF(s string) []byte {
	s = strings.Replace(s, "0x", "", -1)
	s = strings.Replace(s, " ", "", -1)
	a, _ := hex.DecodeString(s)
	return a
}

func FullByteArray(d *[] byte, pas ...interface{}) {
	var ind = 0

	for _, v := range pas {
		switch i := v.(type) {
		case uint8:
			(*d)[ind] = byte(i)
			ind += 1
		case uint16:
			(*d)[ind] = byte(i >> 8)
			(*d)[ind+1] = byte(i)
			ind += 2
		case uint32:
			(*d)[ind] = byte(i >> 24)
			(*d)[ind+1] = byte(i >> 16)
			(*d)[ind+2] = byte(i >> 8)
			(*d)[ind+3] = byte(i)
			ind += 4
		case int32:
			(*d)[ind] = byte(i >> 24)
			(*d)[ind+1] = byte(i >> 16)
			(*d)[ind+2] = byte(i >> 8)
			(*d)[ind+3] = byte(i)
			ind += 4
		case int:
			(*d)[ind] = byte(i >> 24)
			(*d)[ind+1] = byte(i >> 16)
			(*d)[ind+2] = byte(i >> 8)
			(*d)[ind+3] = byte(i)
			ind += 4
		//case string:
		default:
			Logger.Println("Unresolved type ",i)
		}

	}
}

func Bytes2hexF(s []byte) string {
	str := ""
	for k, v := range s {
		str += fmt.Sprintf("%02x", v)
		if k > 0 && k % 2 == 1 {
			str += " "
		}
	}
	return str
}

func PrintByteArray(pas ...interface{}) {
	if len(pas) == 1 {
		if v, ok := pas[0].([]byte); ok {
			Logger.Println(Bytes2hexF(v))
		}
	}
	if len(pas) == 2 {
		v0, ok0 := pas[0].(string)
		v1, ok1 := pas[1].([]byte)
		if ok0 && ok1 {
			Logger.Println(v0, Bytes2hexF(v1))
		}
	}
}

func GetSerialPorts() []string {
	command := "/bin/bash"
	params := []string{"-c", "ls /dev/ttyUSB*"}

	r, _ := ExecCommand(command, params)
	// fmt.Println("After execCommand",r,err,len(r))
	return r
}

func ExecCommand(commandName string, params []string) ([]string, error) {
	var contentArray = make([]string, 0, 5)
	// contentArray = contentArray[0:0]
	cmd := exec.Command(commandName, params...)
	// Show CMD
	// fmt.Printf("Run CMD: %s\n", strings.Join(cmd.Args[1:], " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return contentArray, err
	}
	err = cmd.Start()
	if err != nil {
		return contentArray, err
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
		line = strings.Trim(line, "\n")
		// fmt.Print(line)
		index++
		contentArray = append(contentArray, line)
	}
	err = cmd.Wait()
	return contentArray, err
}
