package datapkg

import (
	"fmt"
	"scentrealm_bcc/datapkg/defines/CTLBCC"
	"scentrealm_bcc/datapkg/defines/CTLSP"
	"scentrealm_bcc/datapkg/defines/DeviceType"
	"scentrealm_bcc/utils"
)


type DataPackageCTLBcc struct{
	DataPackage
	dpcs DataPackageCTLSp
}


func (dataPkg *DataPackageCTLBcc) Init(){
	dataPkg.Header = 0xf5
	dataPkg.End = 0x55
	dataPkg.PkgPayloadConfig = CTLBCC.PkgPayloadConfig
	dataPkg.PkgCMDNames = CTLBCC.PkgCMDNames
	dataPkg.MinLength = 11
	dataPkg.PayloadOffset = 8
	dataPkg.PayloadLengthOffset = 0
	dataPkg.FuncCodeIndex = 7
	dataPkg.AddressTypeOffset = 1
	dataPkg.AddressLength = 2
	dataPkg.FromType =  DeviceType.Controller
	dataPkg.FromAddress = 0x0001
	dataPkg.ToType = DeviceType.SmellDevice
	dataPkg.ToAddress = 0xffff

	dataPkg.dpi = dataPkg
}

func (dataPkg *DataPackageCTLBcc) VerifyCMDReq(CMD byte, Ins []byte, Length byte) bool {
	dataPkg.TargetFuncCode = CMD
	dataPkg.Data = Ins
	dataPkg.DataLength = Length
	return dataPkg.VerifyDataPackageReq()
}

func (dataPkg *DataPackageCTLBcc) VerifyDataPackageReq() bool {
	CSPDP := dataPkg.dpcs
	CSPDP.Init()
	CSPDP.Data = dataPkg.Data
	CSPDP.DataLength = dataPkg.DataLength
	CSPDP.TargetFuncCode = CTLSP.ReqTransparentTransmission
	TT := CSPDP.VerifyDataPackage()
	fmt.Println("TT ",TT)
	if TT {
		dataPkg.Data = CSPDP.GetPayloadAddr()[2:]
		dataPkg.DataLength = dataPkg.DataLength - CSPDP.MinLength - 2
		return dataPkg.DataPackage.VerifyDataPackage()
	}
	return TT
}

func (dataPkg *DataPackageCTLBcc) VerifyDataPackage() bool {
	CSPDP :=  dataPkg.dpcs
	CSPDP.Init()
	CSPDP.Data = dataPkg.Data
	CSPDP.DataLength = dataPkg.DataLength
	CSPDP.TargetFuncCode = CTLSP.RespTransparentTransmission
	TT := CSPDP.VerifyDataPackage()
	if TT {
		dataPkg.Data = CSPDP.GetPayloadAddr()[2:]
		dataPkg.DataLength = dataPkg.DataLength - CSPDP.MinLength - 2
		return dataPkg.DataPackage.VerifyDataPackage()
	}
	return TT
}

func (dataPkg *DataPackageCTLBcc) AssembleQueryDeviceDevice(device uint32,channel byte) []byte  {
	dataPkg.Init()
	dataPkg.FuncCode = CTLBCC.ReqQueryDeviceState
	dataPkg.PayloadLength = 0
	dataPkg.ToAddress = device
	payload := dataPkg.ToByteArray()
	dataPkg.dpcs.Init()
	return dataPkg.dpcs.AssembleReqTransparentTransmission(payload,len(payload),channel)
}

func (dataPkg *DataPackageCTLBcc) AssemblePlaySmell(smell uint32,duration uint32,channel byte) []byte  {
	// f551 0013 f501 0001 0200 ff01 0000 ee33 ffff 5566 04de 5507 6e55
	dataPkg.Init()
	dataPkg.FuncCode = CTLBCC.ReqPlaySmell
	dataPkg.PayloadLength = 8
	dataPkg.Payload = make([]byte,8)
	utils.FullByteArray(&dataPkg.Payload,smell,duration)
	payload := dataPkg.ToByteArray()
	dataPkg.dpcs.Init()
	return dataPkg.dpcs.AssembleReqTransparentTransmission(payload,len(payload),channel)
}

func (dataPkg *DataPackageCTLBcc) AssembleStopPlay(channel byte) []byte  {
	// f551 0013 f501 0001 0200 ff01 0000 ee33 ffff 5566 04de 5507 6e55
	dataPkg.Init()
	dataPkg.FuncCode = CTLBCC.ReqStopPlay
	dataPkg.PayloadLength = 0
	payload := dataPkg.ToByteArray()
	dataPkg.dpcs.Init()
	return dataPkg.dpcs.AssembleReqTransparentTransmission(payload,len(payload),channel)
}
