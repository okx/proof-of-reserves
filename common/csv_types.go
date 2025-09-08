package common

// CSV line data structure
type CSVLine struct {
	LineNumber     int
	DigitalAsset   string
	Network        string
	Address        string
	SignedMessage  string
	SignedMessage2 string
	Message        string
	PublicKey      string
	Owner1         string
	Owner2         string
	RawLine        string
}

// Verification result structure
type VerifyResult struct {
	Line    CSVLine
	Success bool
	Coin    string
	Error   string
}

// Failed line information structure
type FailedLineInfo struct {
	LineNumber   int
	Coin         string
	DigitalAsset string
	Address      string
	ErrorMessage string
}
