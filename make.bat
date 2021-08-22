if not exist gen mkdir gen
go mod tidy
go build -o gen/sunagentd.exe cmd/sunagentd/sunagentd.go
