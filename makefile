port=3000
wsPort=4000
wsServer=dev-ws-server/server.js

start-websocket-server:
	node $(wsServer)
watch-websocket-ping:
	watchexec --restart -w . -e css,js,tmpl -- "make websocket-ping"
websocket-ping:
	until nc -z localhost $(port); do sleep 0.5; done && echo '' | websocat ws://localhost:$(wsPort)


watch-prebuild:
	watchexec --restart --ignore "cmd/web2/main.go" -w . -e go -- "make prebuild"
prebuild:
	go run cmd/cli/prebuild.go

watch-run:
	watchexec --restart -w . -e go,tmpl -- "make run"
run:
	go run cmd/web/* -addr localhost:$(port)
