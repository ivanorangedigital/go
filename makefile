host=localhost
port=3000

dev:
	watchexec --restart -w . -e go,tmpl -- go run cmd/web/* -addr $(host):$(port)
ws:
	watchexec --restart -w . -e css,js,tmpl -- "until nc -z $(host) $(port); do sleep 0.5; done && echo '' | websocat ws://localhost:4000"
