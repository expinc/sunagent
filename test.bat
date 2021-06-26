if not exist gen mkdir gen
go test -v -coverpkg=./... ./... -coverprofile gen/coverprofile > gen/test.log
go tool cover -html ./gen/coverprofile -o gen/test.html
