package CTLSP

//enum CtlSPInstructions
const (
	ReqTransparentTransmission         = 0x51
	ReqAssignControllerPhysicalChannel = 0x26
	ReqGetControllerPhysicalChannel    = 0x27
	ReqGetControllerLogicalChannel     = 0x34
	RespGetControllerLogicalChannel    = 0xb4
	RespGetControllerPhysicalChannel   = 0xa7
	RespAssignControllerPhysicalChannel= 0xa6
	RespTransparentTransmission        = 0xd1
)

var PkgPayloadConfig = map[int]int{
	ReqGetControllerPhysicalChannel:  0,
	RespGetControllerPhysicalChannel: 1,

	ReqGetControllerLogicalChannel:  0,
	RespGetControllerLogicalChannel : 1,


	ReqAssignControllerPhysicalChannel:1,
	RespAssignControllerPhysicalChannel:1,
	RespTransparentTransmission:      -2,
	ReqTransparentTransmission:       -2,
}

var PkgCMDNames = map[int]string{
	ReqGetControllerPhysicalChannel:  "ReqGetControllerPhysicalChannel",
	RespGetControllerPhysicalChannel: "RespGetControllerPhysicalChannel",

	ReqGetControllerLogicalChannel:  "ReqGetControllerLogicalChannel",
	RespGetControllerLogicalChannel : "RespGetControllerLogicalChannel",


	ReqAssignControllerPhysicalChannel:"ReqAssignControllerPhysicalChannel",
	RespAssignControllerPhysicalChannel:"RespAssignControllerPhysicalChannel",
	RespTransparentTransmission:      "RespTransparentTransmission",
	ReqTransparentTransmission:       "ReqTransparentTransmission",
}