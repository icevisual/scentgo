package datapkg

import (
	"fmt"
	"scentrealm_bcc/datapkg/defines/CTLSP"
	"scentrealm_bcc/datapkg/defines/DeviceType"
	"scentrealm_bcc/utils"
)

type DataPackage struct {
	///////////////////////数据包请求属性////////////////
	/* 传入的数据  */
	Data       []byte
	DataLength byte
	/* 目标功能码  */
	TargetFuncCode byte
	/* 超时时间（毫秒）  */
	TimeoutIntervalInMilliSecond uint32
	///////////////////////数据包类属性//////////////////
	/* 数据包包头 0xf5 */
	Header byte
	/* 数据包包尾0x55 */
	End byte
	/* 配置的各个类型的数据包的包长 */
	PkgPayloadConfig map[int]int
	PkgCMDNames      map[int]string
	/* 数据包最小长度 */
	MinLength byte
	/* 荷载数据偏移 */
	PayloadOffset byte
	/* Payload Length Offset */
	PayloadLengthOffset byte
	/* 功能码索引 */
	FuncCodeIndex byte
	/* 地址区域偏移 */
	AddressTypeOffset byte
	/* 地址字节数 */
	AddressLength byte
	///////////////////////数据包实例属性////////////////
	FuncCode byte
	Payload  []byte
	/* Payload Length, 0xff not set,0xfe from offset  */
	PayloadLength byte
	/* 源设备类别 */
	FromType    byte
	FromAddress uint32
	/* 目标设备类别 */
	ToType    byte
	ToAddress uint32
	/* 校验和 */
	CheckBit [2]byte
	/////////////////////////////////////////////////////
	/* 处理方法重载 */
	dpi DataPackageInter
}

type DataPackageInter interface {
	/* 验证数据包完整性、比对 TargetFuncCode  */
	VerifyDataPackage() bool
	/* TargetFuncCode、Data、DataLength 赋值 调用 VerifyDataPackage */
	VerifyCMD(CMD byte, Ins []byte, Length byte) bool
	/* 分析接收到的数据 */
	DataAnalysis(Ins []byte, Length byte) bool
	// WaitResponse(CMD byte, IntervalInMilliSecond uint32) bool
	GetPayloadLength() int
	// 计算校验和
	CalculateCheckSum(Ins []byte, Start byte, InsLength byte) uint16
	// 获取 Payload 切片
	GetPayloadAddr() []byte
	// 实例类转 byte 数组
	ToByteArray() []byte
	// 从 byte 数组载入实例属性
	LoadFromData()
	// 描述实例数组
	Describe()
	// 获取功能码名字
	GetFuncCodeName(funcCode byte) string
}

func (dataPkg *DataPackage) GetFuncCodeName(funcCode byte) string {
	if v, ok := dataPkg.PkgCMDNames[int(funcCode)]; ok {
		return v
	}
	return fmt.Sprintf("Unknown (0x%02x)", funcCode)
}

func (dataPkg *DataPackage) Describe() {
	str := fmt.Sprintf("Describe DataPackage %p\n",dataPkg)
	str += fmt.Sprintf("FuncCode : 0x%02x (%s)\n", dataPkg.FuncCode, dataPkg.GetFuncCodeName(dataPkg.FuncCode))
	str += fmt.Sprintf("From 	 : %5d \tType : %s\n", dataPkg.FromAddress, DeviceType.GetDeviceType(dataPkg.FromType))
	str += fmt.Sprintf("To   	 : %5d \tType : %s\n", dataPkg.ToAddress, DeviceType.GetDeviceType(dataPkg.ToType))
	str += fmt.Sprintf("Payload  : %s \n", utils.Bytes2hexF(dataPkg.Payload))
	utils.Logger.Println(str)
}

func (dataPkg *DataPackage) ToByteArray() []byte {
	InsLength := dataPkg.MinLength + dataPkg.PayloadLength
	bArray := make([]byte, InsLength)
	Ins := bArray
	Ins[0] = dataPkg.Header
	if dataPkg.AddressTypeOffset > 0 {
		Ins[dataPkg.AddressTypeOffset] = dataPkg.FromType
		FromAddressT := dataPkg.FromAddress
		ToAddressT := dataPkg.ToAddress
		for i := int(dataPkg.AddressLength - 1); i >= 0; i-- {
			Ins[int(dataPkg.AddressTypeOffset)+1+i] = byte(FromAddressT & 0xff)
			FromAddressT >>= 8
		}
		Ins[dataPkg.AddressTypeOffset+dataPkg.AddressLength+1] = dataPkg.ToType
		for i := int(dataPkg.AddressLength - 1); i >= 0; i-- {
			Ins[int(dataPkg.AddressTypeOffset)+int(dataPkg.AddressLength)+1+1+i] = byte(ToAddressT & 0xff)
			ToAddressT >>= 8
		}
	}
	Ins[dataPkg.FuncCodeIndex] = dataPkg.FuncCode
	if dataPkg.PayloadLength > 0 && dataPkg.Payload != nil {
		for i := 0; i < int(dataPkg.PayloadLength); i++ {
			Ins[int(dataPkg.PayloadOffset)+i] = dataPkg.Payload[i]
		}
	}
	checkSum := dataPkg.dpi.CalculateCheckSum(Ins, 1, InsLength)
	Ins[InsLength-3] = byte((checkSum >> 8) & 0xff)
	Ins[InsLength-2] = byte(checkSum & 0xff)
	Ins[InsLength-1] = dataPkg.End
	return bArray
}

func (dataPkg *DataPackage) VerifyCMD(CMD byte, Ins []byte, Length byte) bool {
	dataPkg.TargetFuncCode = CMD
	dataPkg.Data = Ins
	dataPkg.DataLength = Length
	return dataPkg.dpi.VerifyDataPackage()
}

func (dataPkg *DataPackage) DataAnalysis(Ins []byte, Length byte) bool {
	dataPkg.Data = Ins
	dataPkg.DataLength = Length
	verify := dataPkg.dpi.VerifyDataPackage()
	return verify
}

func (dataPkg *DataPackage) VerifyDataPackage() bool {
	Ins := dataPkg.Data
	dataPkg.PayloadLength = byte(dataPkg.dpi.GetPayloadLength())

	if dataPkg.PayloadLength == 0xff {
		fmt.Printf("TargetFuncCode 0x%x NotFound ", dataPkg.TargetFuncCode)
		return false
	}
	if dataPkg.DataLength < dataPkg.MinLength {
		fmt.Printf("DataLength (%d) less than MinLength(%d)", dataPkg.DataLength, dataPkg.MinLength)
		return false
	}
	InsLength := dataPkg.MinLength + dataPkg.PayloadLength
	dataPkg.Payload = Ins[dataPkg.PayloadOffset : dataPkg.PayloadOffset+dataPkg.PayloadLength]
	if dataPkg.DataLength >= InsLength {
		if Ins[0] == dataPkg.Header && Ins[dataPkg.FuncCodeIndex] == dataPkg.TargetFuncCode && Ins[InsLength-1] == dataPkg.End {
			calChecksum := dataPkg.dpi.CalculateCheckSum(Ins, 1, InsLength)
			dataPkg.CheckBit[0] = byte(0xff & (calChecksum >> 8))
			dataPkg.CheckBit[1] = byte(0xff & calChecksum)
			if Ins[InsLength-2] == dataPkg.CheckBit[1] && Ins[InsLength-3] == dataPkg.CheckBit[0] {
				dataPkg.dpi.LoadFromData()
				return true
			}
		}
		fmt.Printf("H,T,E Error %t,%t,%t ", Ins[0] == dataPkg.Header , Ins[dataPkg.FuncCodeIndex] == dataPkg.TargetFuncCode , Ins[InsLength-1] == dataPkg.End)
		return false
	}
	fmt.Printf("Length Error %d,%d ", dataPkg.DataLength , InsLength)
	return false
}

func (dataPkg *DataPackage) Init() {
	dataPkg.Header = 0xf5
	dataPkg.End = 0x55

	dataPkg.PkgPayloadConfig = CTLSP.PkgPayloadConfig
	dataPkg.PkgCMDNames = CTLSP.PkgCMDNames
	dataPkg.MinLength = 5
	dataPkg.PayloadOffset = 2
	dataPkg.PayloadLengthOffset = 3
	dataPkg.FuncCodeIndex = 1
	dataPkg.AddressTypeOffset = 0
	dataPkg.AddressLength = 0

	dataPkg.dpi = dataPkg
}

func (dataPkg *DataPackage) GetPayloadAddr() []byte {
	return dataPkg.Data[dataPkg.PayloadOffset:]
}

func (dataPkg *DataPackage) CalculateCheckSum(Ins []byte, Start byte, InsLength byte) uint16 {
	var calChecksum uint16 = 0
	for i := Start; i < InsLength-3; i++ {
		calChecksum += uint16(Ins[i])
	}
	return calChecksum
}

func (dataPkg *DataPackage) GetPayloadLength() int {
	key := int(dataPkg.Data[dataPkg.FuncCodeIndex])
	if length, ok := dataPkg.PkgPayloadConfig[key]; ok {
		return length
	}
	return 0xff
}

func (dataPkg *DataPackage) LoadFromData() {
	dataPkg.FuncCode = dataPkg.Data[dataPkg.FuncCodeIndex]
	dataPkg.Payload = dataPkg.Data[dataPkg.PayloadOffset : dataPkg.PayloadOffset+dataPkg.PayloadLength]
	if dataPkg.AddressTypeOffset > 0 {
		dataPkg.FromType = dataPkg.Data[dataPkg.AddressTypeOffset]
		dataPkg.FromAddress = 0
		for i := 0; i < int(dataPkg.AddressLength); i++ {
			if i > 0 {
				dataPkg.FromAddress <<= 8
			}
			dataPkg.FromAddress |= uint32(dataPkg.Data[int(dataPkg.AddressTypeOffset)+i+1]);
		}
		dataPkg.ToType = dataPkg.Data[dataPkg.AddressTypeOffset+dataPkg.AddressLength+1]
		dataPkg.ToAddress = 0
		for i := 0; i < int(dataPkg.AddressLength); i++ {
			if i > 0 {
				dataPkg.ToAddress <<= 8
			}
			dataPkg.ToAddress |= uint32(dataPkg.Data[int(dataPkg.AddressTypeOffset)+int(dataPkg.AddressLength)+1+i+1])
		}
	}
}
