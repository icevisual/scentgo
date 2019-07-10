package datapkg

import "scentrealm_bcc/utils"

func AssembleGetControllerPC() []byte {
	// C[S]: f527 0027 55
	// C[R]: f5a7 9501 3c55
	return utils.Hex2bytes("f527002755")
}