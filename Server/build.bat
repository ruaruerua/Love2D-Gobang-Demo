@echo off
set GOPROXY=https://goproxy.cn

set servers=Server\Center;Server\Login;Server\TestClient

for %%I in (%servers%) do (
	echo build %%I
	@echo on
	go build -o bin -v %%I
	@echo off
)

echo Build Finish!