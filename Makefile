.PHONY: windows
windows:
	GOOS=windows CGO_ENABLED=0 go build -o hive.exe .
