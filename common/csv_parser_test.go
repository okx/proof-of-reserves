package common

import (
	"reflect"
	"testing"
)

func TestParseCSVLineQuotedJSONArrayPreservesOwnerColumns(t *testing.T) {
	line := `ETH,3,ETH,0xcontract,0xsig1,0xsig2,DTT_OKC_AA26,"[""pub1"",""pub2""]",0xowner1,0xowner2`
	want := []string{
		"ETH",
		"3",
		"ETH",
		"0xcontract",
		"0xsig1",
		"0xsig2",
		"DTT_OKC_AA26",
		`["pub1","pub2"]`,
		"0xowner1",
		"0xowner2",
	}

	if got := ParseCSVLine(line); !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseCSVLine() = %#v, want %#v", got, want)
	}
}

func TestParseCSVLineCompatibility(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "ordinary fields",
			line: "asset,network,address,message",
			want: []string{"asset", "network", "address", "message"},
		},
		{
			name: "quoted JSON object",
			line: `EOS,EOS,addr,"{""publicKey1"":""key1"",""publicKey2"":""key2""}",,`,
			want: []string{"EOS", "EOS", "addr", `{"publicKey1":"key1","publicKey2":"key2"}`, "", ""},
		},
		{
			name: "legacy unquoted JSON object",
			line: `EOS,EOS,addr,{"publicKey1":"key1","publicKey2":"key2"},,`,
			want: []string{"EOS", "EOS", "addr", `{"publicKey1":"key1","publicKey2":"key2"}`, "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCSVLine(tt.line); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseCSVLine() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseCSVLineMixedQuotedFieldAndLegacyJSONObject(t *testing.T) {
	line := `ETH,"memo,with comma",addr,{"key1":"value1","key2":"value2"},owner1,owner2`
	want := []string{
		"ETH",
		"memo,with comma",
		"addr",
		`{"key1":"value1","key2":"value2"}`,
		"owner1",
		"owner2",
	}

	if got := ParseCSVLine(line); !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseCSVLine() = %#v, want %#v", got, want)
	}
}
