test:
	env GO111MODULE=on go test -count=1 -v .

certs:
	env GO111MODULE=on PRINT_CERTS=true go test -v . -run TestGetCert

ctx:
	env GO111MODULE=on TEST_CTX=true go test -v . -count 1 -run TestCtx