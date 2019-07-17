package tests

import (
	"scentrealm_bcc/datapkg"
	"scentrealm_bcc/datapkg/defines/CTLBCC"
	"scentrealm_bcc/datapkg/defines/CTLSP"
	"scentrealm_bcc/utils"
)
var logger = utils.Logger
func TestPKG() {
	pkg := &datapkg.DataPackageCTLSp{}
	pkg.Init()

	r := pkg.VerifyCMD(
		CTLSP.RespGetControllerPhysicalChannel,
		utils.Hex2bytes("f5a777011e55"), 6)
	logger.Println("RespGetControllerPhysicalChannel", r)
	utils.PrintByteArray(pkg.ToByteArray())
	pkg.Describe()
	r = pkg.VerifyCMD(
		CTLSP.RespGetControllerLogicalChannel,
		utils.Hex2bytesF("f5 b4 00 00 b4 55"), 6)
	logger.Println("RespGetControllerLogicalChannel", r)
	utils.PrintByteArray("[ToByteArray]", pkg.ToByteArray())
	pkg.Describe()

	utils.PrintByteArray("[AssembleGetControllerPC]", pkg.AssembleGetControllerPC())
	utils.PrintByteArray("[AssembleGetControllerC]", pkg.AssembleGetControllerC())

	pkgBcc := &datapkg.DataPackageCTLBcc{}
	pkgBcc.Init()
	InsRespQueryState := []byte{
		0xf5, 0xd1, 0x00, 0x12,
		0xf5,
		0x02, 0x07, 0xde,
		0x01, 0x00, 0x00,
		0xa1,
		0x00, 0x0f, 0xd8, 0x06, 0x34, 0x47, 0x31,
		0x03, 0x22, 0x55,
		0x05, 0x74, 0x55,
	}
	InsReqQueryState := []byte{
		0xf5, 0x51, 0x00, 0x0b,
		0xf5,
		0x01, 0x00, 0x01,
		0x02, 0x07, 0xde,
		0x21,
		0x01, 0x0a, 0x55,
		0x02, 0xbb, 0x55,
	}

	r = pkgBcc.VerifyCMD(
		CTLBCC.RespQueryDeviceState,
		InsRespQueryState, byte(len(InsRespQueryState)))
	logger.Println("ReqQueryDeviceState", r)
	utils.PrintByteArray("[ToByteArray]", pkgBcc.Payload)
	pkgBcc.Describe()
	r = pkgBcc.VerifyCMDReq(
		CTLBCC.ReqQueryDeviceState,
		InsReqQueryState, byte(len(InsReqQueryState)))
	logger.Println("ReqQueryDeviceState", r)
	pkgBcc.Describe()
	utils.PrintByteArray("[ToByteArray]", pkgBcc.Payload)
	ReqQueryArray := pkgBcc.AssembleQueryDeviceDevice(1000, 0)
	r = pkgBcc.VerifyCMDReq(
		CTLBCC.ReqQueryDeviceState,
		ReqQueryArray, byte(len(ReqQueryArray)))
	logger.Println("ReqQueryDeviceState", r)
	pkgBcc.Describe()

}



