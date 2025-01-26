go test -v ./tests/... \
  -coverpkg=./pkg/... \
  -coverprofile=coverage.out \
  -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
echo "Report generated in coverage.html"