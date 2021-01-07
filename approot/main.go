package main

import (
	"encoding/hex"
	"fmt"

	_ "github.com/tivizi/forarun/approot/config"
	_ "github.com/tivizi/forarun/pkg/api/admin"
	_ "github.com/tivizi/forarun/pkg/api/site"
	_ "github.com/tivizi/forarun/pkg/base"
	_ "github.com/tivizi/forarun/pkg/daemon"
	_ "github.com/tivizi/forarun/pkg/extra/api"
	_ "github.com/tivizi/forarun/pkg/resources"
	_ "github.com/tivizi/forarun/pkg/resources/admin"
	_ "github.com/tivizi/forarun/pkg/resources/site"

	_ "github.com/tivizi/forarun/approot/docs"
)

// @title FORARUN 自助建站系统开放接口文档
// @version 1.0
// @description 轻量级自助建站系统，随时随地维护站点
// @host fora.run

// @contact.name Tivizi
// @contact.url https://fora.run
// @contact.email tivizi@163.com
func main() {
	// site.DefaultEngine().Run()
	fmt.Println(hex.EncodeToString([]byte("滑稽")))
}
