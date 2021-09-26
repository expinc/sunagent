if not exist gen mkdir gen
go mod tidy
go build -o gen/sunagentd.exe cmd/sunagentd/sunagentd.go

copy config.conf gen /Y
if not exist gen\grimoires mkdir gen\grimoires
xcopy etc\grimoires gen\grimoires /E /Y
