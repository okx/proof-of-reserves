package common

import "testing"

func TestGetFilAddressByPublicKey(t *testing.T) {
	args := []struct {
		publicKey string
		wants     string
	}{
		{
			publicKey: "0x04c7d2209a4b286046cdeaf457e499a40a9a1da5d7bc6e85c05e5ac9e6af9c7a35063c8a8efaa7cc4cd294c3b76dd4b0a3f5773cc421fef44e6a99914c8c85c971",
			wants:     "f12cs7ppvnhwhma3xzhkm4pavq2q47blmprcxvg6i",
		},
	}
	for _, tt := range args {
		res := GetFilAddressFromPublicKey(tt.publicKey)
		if tt.wants != res {
			t.Errorf("Want %s, Got %s", tt.wants, res)
		}
	}
}
