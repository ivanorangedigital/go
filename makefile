processManagerPath=cmd/dev/process-manager/

start:
	gcc $(processManagerPath)/process-manager.c -o $(processManagerPath)/exec && $(processManagerPath)/exec
