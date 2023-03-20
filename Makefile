.PHONY: build-local

all:checkbalance merklevalidator verifyaddress

checkbalance:
	 go build -o build/CheckBalance cmd/checkbalance/main.go

merklevalidator:
	go build -o build/MerkleValidator cmd/merklevalidator/main.go

merklevalidator-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/MerkleValidator-linux cmd/merklevalidator/main.go

merklevalidator-windows:
	GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o build/MerkleValidator.exe cmd/merklevalidator/main.go

verifyaddress:
	go build -o build/VerifyAddress cmd/verifyaddress/main.go

