package common

import "testing"

func TestCreateAddressDescriptor(t *testing.T) {
	args := []struct {
		addrType     string
		redeemScript string
		mSigns       int
		nKeys        int
		wants        string
	}{
		{
			addrType:     "P2PKH",
			redeemScript: "03dbf340db709173d676ac4d0c598cc5811e3d8a7e1aae1f33f4237349367c2925",
			mSigns:       1,
			nKeys:        1,
			wants:        "pkh(03dbf340db709173d676ac4d0c598cc5811e3d8a7e1aae1f33f4237349367c2925)",
		},
		{
			addrType:     "P2SH",
			redeemScript: "52210251c789f59bc870ae263db2fac71c1625bb16bff840eee169420bcba14443f20b210343f2d97938abfcf8201adf46ce50cf7119bb419684268c145902ef5cb7c3e76321027d8ddf369f5dbd880f6e8e04a74b728454445704d844c9042bc43bdc9cea3f6e53ae",
			mSigns:       2,
			nKeys:        3,
			wants:        "sh(multi(2,0251c789f59bc870ae263db2fac71c1625bb16bff840eee169420bcba14443f20b,0343f2d97938abfcf8201adf46ce50cf7119bb419684268c145902ef5cb7c3e763,027d8ddf369f5dbd880f6e8e04a74b728454445704d844c9042bc43bdc9cea3f6e))",
		},
		{
			addrType:     "P2WSH",
			redeemScript: "52210324f8cdeaf96781d99b95dc90af2184869c4dbae236bfb7955112bb8e34221ece21024110529824463881ce90e6d1bc3c093a11dc5ac326bfc1f4b763a211e37801052103b20e5987c8375d93986453757e8b6cbbe58666f847dfc2e8faa4fc11f49efd1253ae",
			mSigns:       2,
			nKeys:        3,
			wants:        "wsh(multi(2,0324f8cdeaf96781d99b95dc90af2184869c4dbae236bfb7955112bb8e34221ece,024110529824463881ce90e6d1bc3c093a11dc5ac326bfc1f4b763a211e3780105,03b20e5987c8375d93986453757e8b6cbbe58666f847dfc2e8faa4fc11f49efd12))",
		},
		{
			addrType:     "P2WSH",
			redeemScript: "5221026064e5b88c4fff7dba7dc0300db8dbfc1faff14f9ddbaacbcaa4f70124de0e93210331870350912385ca9a9d537e9cf9d80c6c9558e31d654f82f3164fdc5955e9642103c7b133a0f463a501d8c58c8eb8c7b6e9e4ddfb7d4a7bf6365a4732201569bc8353ae",
			mSigns:       2,
			nKeys:        3,
			wants:        "wsh(multi(2,026064e5b88c4fff7dba7dc0300db8dbfc1faff14f9ddbaacbcaa4f70124de0e93,0331870350912385ca9a9d537e9cf9d80c6c9558e31d654f82f3164fdc5955e964,03c7b133a0f463a501d8c58c8eb8c7b6e9e4ddfb7d4a7bf6365a4732201569bc83))",
		},
	}
	for _, tt := range args {
		res, err := CreateAddressDescriptor(tt.addrType, tt.redeemScript, tt.mSigns, tt.nKeys)
		if tt.wants != res || err != nil {
			t.Errorf("Want %s, Got %s", tt.wants, res)
		}
	}
}

func TestGuessAddressType(t *testing.T) {
	args := []struct {
		address string
		wants   string
	}{
		{
			address: "bc1quhruqrghgcca950rvhtrg7cpd7u8k6svpzgzmrjy8xyukacl5lkq0r8l2d",
			wants:   "P2WSH",
		},
		{
			address: "3CHeHsCpH9QmX2hmbzkZinqjtUtqseNWrV",
			wants:   "P2SH",
		},
	}
	for _, tt := range args {
		res := GuessAddressType(tt.address)
		if tt.wants != res {
			t.Errorf("Want %s, Got %s", tt.wants, res)
		}
	}
}
