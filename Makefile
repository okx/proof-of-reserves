.PHONY: build-local

all:checkbalance merklevalidator verifyaddress

checkbalance:
	 go build -o build/CheckBalance cmd/checkbalance/main.go

merklevalidator:
	go build -o build/MerkleValidator cmd/merklevalidator/main.go

verifyaddress:
	go build -o build/VerifyAddress cmd/verifyaddress/main.go

