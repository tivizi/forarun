package base

import (
	"os"

	"github.com/bububa/ip2region-go"
)

var ipStore *ip2region.Ip2Region

func init() {
	if ok, _ := fileExists("ip2region.db"); ok {
		ipStore, _ = ip2region.New("ip2region.db")
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
	geo, err := ipStore.BinarySearch(ip)
	if err != nil {
		return err.Error()
	}
	return geo.String()
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
