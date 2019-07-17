package CTLBCC

//enum CtlBroadcastInstructions
const (
	ReqQueryDeviceState  = 0x21
	ReqPlaySmell         = 0x01
	ReqStopPlay          = 0x02
	ReqPlayScript        = 0x03 //  播放脚本
	ReqStopScript        = 0x04 //  停止脚本
	RespQueryDeviceState = 0xa1
)

var PkgPayloadConfig = map[int]int{
	ReqQueryDeviceState:  0,
	RespQueryDeviceState: 7,
	ReqPlaySmell:         8,
	ReqStopPlay:          0,
	ReqPlayScript:        6,
	ReqStopScript:        0,
}

var PkgCMDNames = map[int]string{
	ReqQueryDeviceState:  "ReqQueryDeviceState",
	RespQueryDeviceState: "RespQueryDeviceState",
	ReqPlaySmell:         "ReqPlaySmell",
	ReqStopPlay:          "ReqStopPlay",
	ReqPlayScript:        "ReqPlayScript",
	ReqStopScript:        "ReqStopScript",
}
