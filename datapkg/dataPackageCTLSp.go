package datapkg

import (
	"scentrealm_bcc/datapkg/defines/CTLSP"
)

type DataPackageCTLSp struct {
	DataPackage
}

func (dataPkg *DataPackageCTLSp) Init() {
	dataPkg.Header = 0xf5
	dataPkg.End = 0x55
	dataPkg.PkgPayloadConfig = CTLSP.PkgPayloadConfig
	dataPkg.PkgCMDNames = CTLSP.PkgCMDNames
	dataPkg.MinLength = 5
	dataPkg.PayloadOffset = 2
	dataPkg.PayloadLengthOffset = 3
	dataPkg.FuncCodeIndex = 1
	dataPkg.AddressTypeOffset = 0
	dataPkg.AddressLength = 2

	dataPkg.dpi = dataPkg
}

func (dataPkg *DataPackageCTLSp) GetPayloadLength() int {
	payloadLength := dataPkg.DataPackage.GetPayloadLength()
	if payloadLength == -2 && dataPkg.PayloadLengthOffset > 0 && dataPkg.PayloadLengthOffset < dataPkg.DataLength {
		return int(dataPkg.Data[dataPkg.PayloadLengthOffset]) + 2
	}
	return payloadLength
}

func (dataPkg *DataPackageCTLSp) AssembleReqTransparentTransmission(payload []byte, payloadLength int, channel byte) []byte {
	dataPkg.FuncCode = CTLSP.ReqTransparentTransmission
	dataPkg.PayloadLength = 2 + byte(payloadLength)
	array := make([]byte, 2+payloadLength)
	array[0] = channel
	array[1] = byte(payloadLength)
	copy(array[2:], payload[:payloadLength])
	dataPkg.Payload = array
	return dataPkg.ToByteArray()
}

func (dataPkg *DataPackageCTLSp) AssembleGetControllerPC() []byte {
	dataPkg.FuncCode = CTLSP.ReqGetControllerPhysicalChannel
	dataPkg.PayloadLength = 0
	return dataPkg.ToByteArray()
}

func (dataPkg *DataPackageCTLSp) AssembleGetControllerC() []byte {
	dataPkg.FuncCode = CTLSP.ReqGetControllerLogicalChannel
	dataPkg.PayloadLength = 0
	return dataPkg.ToByteArray()
}
