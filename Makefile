REGISTRY ?= localhost:5000
BIG_VER ?= 1.0
m ?= ci

all: swag run

# 版本
version:
	$(eval VERSION = ${BIG_VER}-$(shell git log -1 --pretty=format:"%h"))
# 运行
run: build
	cd approot;../bin/forarun
# 编译
build:
	go build -o bin/forarun approot/main.go
# 清理编译
clean:
	rm bin/ -rf
# 构建容器镜像
docker: version
	docker build -t ${REGISTRY}/${USERNAME}/forarun:${VERSION} .
# 推送镜像
push: version
	docker push ${REGISTRY}/${USERNAME}/forarun:${VERSION}
# 安装
install: version docker push
	kubectl create deployment -n forarun --image ${REGISTRY}/${USERNAME}/forarun:${VERSION} forarun
	kubectl create service -n forarun clusterip --tcp 80:8080 forarun
uninstall:
	kubectl delete deployment -n forarun forarun
	kubectl delete service -n forarun forarun
# 持续集成
ci: commit version swag build docker push
	kubectl set image -n forarun deployment/forarun forarun=${REGISTRY}/${USERNAME}/forarun:${VERSION}
cd: version push
	kubectl set image -n forarun deployment/forarun forarun=${REGISTRY}/${USERNAME}/forarun:${VERSION}
# 二进制
release: clean build
	cp bin/forarun bin/forarun-${BIG_VER}
# 文档
swag: 
	cd approot;~/.go/bin/swag init --dir ./,../pkg/api/admin,../pkg/api/site
# 提交代码
commit:
	git add .
	git commit -m ${m}