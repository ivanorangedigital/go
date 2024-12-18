run:
	watchexec --restart -w . -e go,tmpl -- go run cmd/web/*
