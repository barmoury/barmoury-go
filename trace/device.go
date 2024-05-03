package trace

import (
	"github.com/mileusna/useragent"
)

type Device struct {
	OsName         string `json:"os_name"`
	OsVersion      string `json:"os_version"`
	EngineName     string `json:"engine_name"`
	DeviceName     string `json:"device_name"`
	DeviceType     string `json:"device_type"`
	DeviceClass    string `json:"device_class"`
	BrowserName    string `json:"browser_name"`
	EngineVersion  string `json:"engine_version"`
	BrowserVersion string `json:"browser_version"`
}

func Build(userAgent string) Device {
	ua := useragent.Parse(userAgent)
	return Device{
		OsName:         ua.OS,
		BrowserName:    ua.Name,
		DeviceName:     ua.Device,
		BrowserVersion: ua.Version,
		OsVersion:      ua.OSVersion,
	}
}
