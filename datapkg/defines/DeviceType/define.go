package DeviceType

import "fmt"

//enum BCC_DeviceType
const
(
	Controller byte = 0x01
	SmellDevice = 0x02
	Syner = 0x05
	DeviceVehicleMounted = 0x08
)


var PkgDeviceTypeNames = map[byte]string{
	Controller:  "Controller",
	SmellDevice: "SmellDevice",

	Syner:  "Syner",
	DeviceVehicleMounted : "DeviceVehicleMounted",
}


func GetDeviceType(typec byte) string {
	if v ,ok := PkgDeviceTypeNames[typec];ok {
		return v
	}
	return fmt.Sprintf("Unknown (%d)",typec)
}