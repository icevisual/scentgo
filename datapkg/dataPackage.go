package datapkg

import (
	"scentrealm_bcc/utils"
)

func AssembleGetControllerPC() []byte {
	// C[S]: f527 0027 55
	// C[R]: f5a7 9501 3c55
	return utils.Hex2bytes("f527002755")
}

func ValidateDataPkg(data []byte,length int) bool {

	if data[0] == 0xf5 && data[length - 1] == 0x55{


	}
	return false
}
//enum BCC_DeviceType
const
(
	Controller = 0x01
	SmellDevice = 0x02
	Syner = 0x05
	DeviceVehicleMounted = 0x08
)

//enum CtlSPInstructions
const (
	ReqTransparentTransmission         = 0x51
	ReqAssignControllerPhysicalChannel = 0x26
	ReqGetControllerPhysicalChannel    = 0x27
	RespGetControllerPhysicalChannel   = 0xa7
	RespAssignControllerPhysicalChannel= 0xa6
	RespTransparentTransmission        = 0xd1
)

//enum CtlBroadcastInstructions
const (
	ReqQueryDeviceState = 0x21
	RespQueryDeviceState = 0xa1
)

var CtlSpPkgPayloadConfig = map[int]int{
	ReqGetControllerPhysicalChannel:  0,
	RespGetControllerPhysicalChannel: 1,
	ReqAssignControllerPhysicalChannel:1,
	RespAssignControllerPhysicalChannel:1,
	RespTransparentTransmission:      -2,
	ReqTransparentTransmission:       -2,
}

var CtlBccPkgPayloadConfig = map[int]int{
	ReqQueryDeviceState:0,
	RespQueryDeviceState:7,
}


