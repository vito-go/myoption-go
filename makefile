appName=myoption
check:
	go build -o /tmp/myoption-main ./cmd/ && rm /tmp/myoption-main
fmt:
	@gofmt -w -s ./
acp:fmt check
ifndef m
	@$(error error: 需要提交说明 请指定参数m, 例如 make acp m=fix)
else
	git add . && git commit -m '$(m)'  && git push origin --all
endif

build:
	go build -o ./bin/$(appName) ./cmd/


kill-9:
	- @lsof -sTCP:LISTEN -i :9130 | awk 'NR==2{print $$2}' | xargs kill -9
	echo "kill the process with signal 9 by port 9130"
kill-15:
	- @lsof -sTCP:LISTEN -i :9130 | awk 'NR==2{print $$2}' | xargs kill -15
	echo "kill the process with signal 15 by port 9130"


start: build kill-15
	./bin/$(appName) -out=false -env ./configs/myoption/test.yaml -daemon

run:
	go run ./cmd/ -env ./configs/myoption/test.yaml

lint:
	@golangci-lint run
build-initdb:
	- mkdir bin
	go build -o ./bin/initdb ./script/initdb.go
initdb: build-initdb
	./bin/initdb doc/postgresql.sql