package base

import (
	"os"

	"github.com/johntech-o/iphelper"
)

var ipStore *iphelper.IpStore

func init() {
	if ok, _ := fileExists("ip.dat"); ok {
		ipStore = iphelper.NewIpStore("ip.dat")
	}
}

// SimpleRegion 简单区域
func SimpleRegion(ip string) string {
	if ipStore == nil {
		return "<UNSUPPORT>"
	}
	if len(ip) == 0 {
		return "<EMPTY>"
	}
	geo, err := ipStore.GetGeoByIp(ip)
	if err != nil {
		return err.Error()
	}
	return geo["country"] + geo["city"] + geo["zone"]
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	// 为nil,说明文件或文件夹存在
	if err == nil {
		return true, nil
	}
	// 不存在错误，明确表示文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	// 其他错误类型，系统错误，不确定是否存在
	return false, err
}
