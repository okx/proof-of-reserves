package common

import (
	"testing"
)

func TestBETHVerifySignature(t *testing.T) {
	addr := "0x82f5af1eb567301d0f71fe56e4bf8aa6e4ffa00c9ac92a8f67aae056b7f19e18be645f6d46c7862de8f1342fc4786853"
	msg := "hello world"
	sign := "0x98d6aa7da816638b8c27a1b6ff9a46529948960743300814f0413499b04886cdc65f63b273461ebaf42be17f08ca01fa13c964bf5a2d988c5c7fe745b472a90dcc86b522543c70020b5ce525b5ffaf6e4e4a7ce80b40049a509b12c32045541a"
	if err := VerifyBETH(addr, msg, sign); err != nil {
		t.Errorf(err.Error())
	}
}

func TestTRXVerifySignature(t *testing.T) {
	addr := "TEjxQjU3CxkFrSDcPfHwZXSuPpCpdQ27NJ"
	msg := "hello world"
	sign := "0xcd1e3903dc047ea881f7da1647fa3372f37ee6a1cf0726477a20e267408af43f3f9c3a43f7f15e6bf674c9f0776866b6d6a770ce998b29cc03f11f2cb98df5821c"
	if err := VerifyTRX(addr, msg, sign); err != nil {
		t.Errorf(err.Error())
	}

	sign = "0xaddfb6bc248de8de0051d3ea225496091af596a5fffed3ee19a93c827687974d3305f869a86208e03886ec9d1423bb264405b6ef0813b3751080f82bd7a906451c"
	if err := VerifyTRX(addr, msg, sign); err != nil {
		t.Errorf(err.Error())
	}
}

func TestVerifyUtxoCoinSignature(t *testing.T) {
	args := []struct {
		coin   string
		addr   string
		msg    string
		sign1  string
		sign2  string
		script string
	}{
		{
			coin:   "BTC",
			addr:   "1DcT5Wij5tfb3oVViF8mA8p4WrG98ahZPT",
			msg:    "I am an OKX address",
			sign1:  "IA1jDx3zkn4J4F6mCVU68Vm7TwNf+bCsp+hKo3LwV/Y+PlZEoNsajnAHqd/FrEmv5/VAGz7pPiWPOXjmCLRfxIM=",
			sign2:  "",
			script: "",
		},
		{
			coin:   "BTC",
			addr:   "3Fs7C97NmvhWUZ2pSjth9YbTxMv4sk9nHi",
			msg:    "I am a OKX address",
			sign1:  "H2vshvcYTGrUw0XG1AundmbivdrhTWUOTqcXKhN+MqbaEfVYfGkgDhEumiJoEJFhlzuma6bBpg4pXNUHoTENOPI=",
			sign2:  "H1eFA8Y2woAnDqxamcLDMVDr4Jd8g6PiagExWCyzvZNZU8xZ2TKV2RNcbXArRgUfniLzgJFvzmBEUC6vgM5bd7A=",
			script: "522103447bead626f13c79de937c0879b64172e5984456a47350b44e8bd23a02e6895e2103864969c155d42c5f61999bcaafeadfc8574b033142f03b5bf3025c6794570b952103304fa164de84f710e44a563f5038d355d6a36a1d7f25695cba884f0b4b6d184653ae",
		},
		{
			coin:   "BTC",
			addr:   "bc1qpypsu8sytw959yu53dk48eaq9saxumwegzwd4anava9qe40k6gfqyrsxaq",
			msg:    "I am a OKX address",
			sign1:  "IKfykjJJSywz2g/KGXvWwE1aphUvryiBaNQYK0m4Ain7QkNhM7VjwV964DPn4dvpOGyzhQYSAqLIz1BOBXbIQcY=",
			sign2:  "H0FnZYEXYbmrnh8sreeUO7wL8BPPKKSPDyYvZbjn/tOCVpxnZLN3yL8lyyPHkl3NttL7WDHVx/jxG5HBLbR7T5k=",
			script: "522102b514e7ccc2845d3f1ca7181dacab0d1ac277616e753547922f82cc0cdfb5c691210318c30ad87d44e7c8b940b47e6963aabfda0581a1e1aac59019b9b1589179aa7a2103d7a534927a03b195a0082d53ab15145bbac8964ffb09d54869da2e59ea1b100553ae",
		},
		{
			coin:   "BSV",
			addr:   "1Hgc1DnWHwfXFxruej4H5g3ThsCzUEwLdD",
			msg:    "hello world",
			sign1:  "IDtG3XPLpiKOp4PjTzCo/ng8gm4MFTTyHeh/DaPC1XYsYaj5Jr4h8dnxmwuJtNkPkH40rEfnrrO8fgZKNOIF5iM=",
			sign2:  "",
			script: "",
		},
		{
			coin:   "DOGE",
			addr:   "9vbpNnyNRZpSWzptDwxGhw2Vny2yJ4W9V2",
			msg:    "I am an OKX address",
			sign1:  "H/Oog7EXWzpoR7CA4yV6k5IYH+aQreIxHuMmBTH53ucWVKLn8F7lLJNNPwe67ElPYr49Ox2PpXyExr+W9+pYIgw=",
			sign2:  "H3LGI7d7ZIMSRf+S9TrfdUrTM+tPA21MRnLyANqaH3LwWzSZAZtxJo4dGZs8JSL7dPkM9NWtMg7bTlLN8pIU0iE=",
			script: "5221037c660ee71005b5e991068021448ed61a650ef018f56c4614b28ae4618169107c2103aa2fdc3e4a5207c68d452dc42f615dee425eea6c2b7ef61f8a677fbff076fa9a2103d23a7924f45f3288816b91fdaba88ea688020b364bf55e506d21f8c8787dd71853ae",
		},
		{
			coin:   "DASH",
			addr:   "7XswsaSt9HtGfrLvboEEzhsmPBSSky1mnu",
			msg:    "I am an OKX address",
			sign1:  "INc+CF8jxNHTyJvnWzbsUe98mRA9eywBVuK1Dd1YTzLgTKiUa8dwfror+px/hlK/hHQ0R7saUi28ijQRCwYJWA0=",
			sign2:  "IFMFCCk781/sT/oHoUxL5oZpqQZlImQrhJ3odvkVWARcaK4EAl3XPEA/IVgxh6VaSLmMP/Wno1dxWdmxe2gcs3U=",
			script: "522102f5c9ab0dd178eb44cd6baae7c1698ae23caec399d58da93a32509665113152742103729f817997a0442e6a39ee0c15f0cd3a17e55ef768ae6f055b2c323a1cba9eb121021e6568b58ef452791f56d0ba2f2ea1200c730f114ab88701833a0a2cf77ef09e53ae",
		},
		{
			coin:   "LTC",
			addr:   "MH7Qru4KFC6tLSiSq7bTWUDyiXERwzdqKp",
			msg:    "I am an OKX address",
			sign1:  "ID7R92j4VjjajZGAqRWyHVEAGLfbJZnIKD7uZVyMPYNqf3Cp0BxeFSuZ8c+weifZ8M3DBE8F2ftDV7s9XYTLPus=",
			sign2:  "H2wcYS+Osohhi06wr78H4Gn8zsB9Oqw9aTRo7OEdTJNaWNGwoo+6zwsywdfUXk6IgxY7hreo0I/s8fIjHJceFdQ=",
			script: "522103c4a04606ac7b32924f0fc0b8aadbf925c139738701a09138bc7c42cb9a3609da21034b8b29e6f460e6293423173863a9cca093bbbd865a611abd984fd674d2537e5b210350ade45e21816652e9812bb39f55e0e8cad02a0884c3c024ab642a43cd8094a353ae",
		},
		{
			coin:   "LTC",
			addr:   "ltc1qzvcgmntglcuv4smv3lzj6k8szcvsrmvk0phrr9wfq8w493r096ssm2fgsw",
			msg:    "I am an OKX address",
			sign1:  "IGLtPax4D8T/a5YIfTjqT4hoy8j71gG+yFeKPUhH4xKTcgEFxLXrJF+2ZLgPoqKGBt+QSNH6l7RK7bIZzVIvzBA=",
			sign2:  "Hzs4gdaoFHT1hfVqyi2tvoIhGiUTP/vVETd2QPtg8Uu/eGzec6oK5DvVIK0oUMC8G0dfWgrEcJM4IE7MIHS4oBc=",
			script: "5221024ae22b49865207865e57b5056abb1e9f6480a6621a246b7066ab84c3b7f9d2b5210255f79ad111076c25b3f6d5e2e99250b6150b1d1b6af10f860f41f63c68859704210247eb8b75b8e50a8a33a7d9ada81f3ecaa05228bfeec3274bc6520637f2eacbf553ae",
		},
		{
			coin:   "LTC",
			addr:   "LectSLekhkxEQBHofNFHa3pVHndnb4Z9fN",
			msg:    "I am an OKX address",
			sign1:  "ILEkfsIaFagDnMq0pY9qS4yIAjTFlpDknLCVm5dhMHbnZn0SXC+Op9jQfD4wbu+Rpxs1N5ut6c2349GEM8AFx0A=",
			sign2:  "",
			script: "",
		},
		{
			coin:   "BCHN",
			addr:   "3EniMkZiTiCcKsyKGxJbizi5zgPLkkpdYB",
			msg:    "I am an OKX address",
			sign1:  "Hzm4HO21jFamx/+AiVgH5+8MNbFLIpyYIgVUdEc/eDGBVaaNpOgmvigZjcFDikzw97BUMUTGLL00TT5PKSAX/0M=",
			sign2:  "IDgJJKiWo38ujfYIeYx8yiHBmlmzCc7c2yIrUb6Xqep0FnnHMyBhySdnBan8rjs+csO9K6GkGsNNs4p6QY9iIg8=",
			script: "52210318af2bec99f512a736c88fa30a5711c16652591510b9d388af553f3644372af121035bfe589cb7086355f9ab44655c8d700f8e7ba2d90116658682130587eb35b0b5210306913a183d6b4cb7c38a99fdcfd3be24fb171b50de065c8ed0e759ee09d7267453ae",
		},
	}

	for _, tt := range args {
		err := VerifyUtxoCoin(tt.coin, tt.addr, tt.msg, tt.sign1, tt.sign2, tt.script)
		if err != nil {
			t.Errorf("coin: %s, addr: %s, msg: %s, sign1: %s, sign2: %s, script: %s", tt.coin, tt.addr, tt.msg, tt.sign1, tt.sign2, tt.script)
			t.Errorf(err.Error())
		}
	}
}

func TestVerifyEvmCoinSignature(t *testing.T) {
	args := []struct {
		coin string
		addr string
		msg  string
		sign string
	}{
		{
			coin: "ETH",
			addr: "0x52b311c52436789f3754bd199bf3886b8ccbab4c",
			msg:  "I am an OKX address",
			sign: "0x98767aedf0ed8bad7413e7c2e6b134ae6baaf5d913c9a8e2659b93922edfbca90cf5fc97e6385aec280a2b7dcdf7d2a95e91f0d99632ab7ed0c167e5628d3d841c",
		},
		{
			coin: "BABYDOGE-BSC",
			addr: "0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547",
			msg:  "hello world",
			sign: "0x9c271461e5876fac4e5a02aee7a877831a91cee6a24b75cafd8650ac72b2a5e5147e2e90558d4e38d113ff54e734f041687f41268d55ff7850791e1e2833dc061b",
		},
		{
			coin: "ETH",
			addr: "0x8c3cb9665833fd9f79eb14cba16d82bbab6f22d8",
			msg:  "I am an OKX address",
			sign: "0xc75173a3ca53bcfeb7b2bfc16aed036191436085fe1a5c846f7021ae2baf5f81646b5089822399b6ee076bb59974ec6bd425954ca97bff084a74efeca0c8c8c61c",
		},
		{
			coin: "ETH",
			addr: "0xa28062bd708ce49e9311d6293def7df63f2b0816",
			msg:  "I am an Okcoin address",
			sign: "0x462950c4dbbc0f2fb36002ba7e5c2a98dfae7d89203f4dbf152e03304edb444d670c8bbacb78072c4dc1184db245401c73ee1395715afbb8dd600ba6a63e3abc1b",
		},
		{
			coin: "FIL-EVM",
			addr: "f410fa7sh5u6fvd7vt625dx2akhbu3jt7yvkhgdp55ki",
			msg:  "hello world",
			sign: "0x9c271461e5876fac4e5a02aee7a877831a91cee6a24b75cafd8650ac72b2a5e5147e2e90558d4e38d113ff54e734f041687f41268d55ff7850791e1e2833dc061b",
		},
		{
			coin: "FIL-EVM",
			addr: "f410f3py2r3kh22yh7ona3fyp2tyzhxo4p2okiwuqjza",
			msg:  "I am an OKX address",
			sign: "0x05168f6ede2213bebae69256ad78cfcf4f3a99d4b11bacd4a24ec2171becb34e2f690c515e3b1227a52e7e7c4dd99dd0f307c167f768b2201bd82b69af5a112b1c",
		},
		{
			coin: "ETH-LINEA",
			addr: "0xd3d7dbe73bbdd5a5c7a49ca322763c4d400fc240",
			msg:  "I am an OKX address",
			sign: "0x75f1d8203bef472f4f5f2627b186eb6d700544633690f4f368ea3d6d069040ac378b84765bf4fecfc4d62d807158817a0d6f5d2ac13121aba578ff0f6880067b1c",
		},
		{
			coin: "BASE",
			addr: "0xfd92f4e91d54b9ef91cc3f97c011a6af0c2a7eda",
			msg:  "I am an OKX address",
			sign: "0x20c683960499e200a78d0f3011f821aea901e6cb26371f0d4edd8b64fbf0efb80b147346a8acfb03e037a40ad83d2ef189af3ca9f32f07e99ba644aac1d951d71b",
		},
		{
			coin: "OKB-X1-USDT",
			addr: "0x2d22604d6bbf51839c404aef5c65443e424e0945",
			msg:  "I am an OKX address",
			sign: "0x6fc4257d86eac2355cfe3998c6524e5f84af3fed014d03810b4ffca2248d076d2093afc5f4ec7bfea166fcacda81e1cf9f0f707806f9fcfd7d2e3163e456c1981b",
		},
		{
			coin: "OKB-X1",
			addr: "0x74be7cb4fadcda2e9df4dd7fa4dc3a069888838e",
			msg:  "I am an OKX address",
			sign: "0xa3632aa7bf9f83a77e85c35cee44354bed5a7b243e3814425eb2e80f809355c44f28fd1b8d53522a2b6d4387039ab47c6c4dd8b971603469b6b0ba1f9845b7b11b",
		},
	}
	for _, tt := range args {
		err := VerifyEvmCoin(tt.coin, tt.addr, tt.msg, tt.sign)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestVerifyEd25519CoinSignature(t *testing.T) {
	args := []struct {
		coin string
		addr string
		pub  string
		msg  string
		sign string
	}{
		{
			coin: "SOL",
			addr: "7bzoTJhZmpU1vQVjN63fQ3iVYmWCVgQh1sYSqsjuapU9",
			pub:  "0x621d398b19304995ee140c21afc544d62382d387b5c08dfd096b475a304339ea",
			msg:  "hello world",
			sign: "0x282c737229f72d03275ac7bc5955da027d693d90dd9e6d4c2aafcc4f272de3be0be08637552027abb68e4d2818e060846b002e490d8bbe36e6dae8b2508fd40a",
		},
		{
			coin: "DOT",
			addr: "15WXogcgXnHsZ1FeuNc6cg34i8R6JCXgWrDLNLERLLesJ7bf",
			pub:  "0xc776bfbeeeb0b1ddd1ce6cccf55ce795f5306bf63de37d72e5af50b3be23ce49",
			msg:  "hello world",
			sign: "0x996528cea9ae0ef66a0c1782cf281726f3e167906e9eb61161558c482f5b92b5439e0662b2b72906cd2a155903fed4b739652a4ddd97689618793e90d8f2d608",
		},
		{
			coin: "APTOS",
			addr: "0x327dd297dfacf7c2d8207aaa23c0f0e8bcaf4c1612febbf63b9f7376810b8ec8",
			pub:  "0x61f579fc779146304353027b425a216d8015889c5f3b715ad26135b862f3bf84",
			msg:  "hello world",
			sign: "0xe5eea05d4156e1aef7867739b86f560b3b6a14a9525b53b436b5ff16ce8ca9490d4e5586ddc469b43453cf9796d87a4c4d3ead5d8dd3a2e88026713ae866e30b",
		},
		{
			coin: "TONCOIN-NEW",
			addr: "EQA5rifVSCc8qQfpCXvq4zJGJPsA0EPCDoWdtg234INftsWj",
			pub:  "0x3d2696e3d5cbc9047b338e6a56552db1d43ca6e063bc7aa667b18005984372d2",
			msg:  "hello world",
			sign: "0xaa406900fdf658e793850d7d47798fa501098db4a6697ac460c1d2800152f40174d2705f1ec87b1a0b34434647b0efed2b7b70569bc00e8bbc3561c372aacc0b",
		},
		{
			coin: "APTOS",
			addr: "0x4062a7375a1ab2bb8bc33d2c1cb2f6b0afa1cef98ea7a5afec7afc96675fe1d",
			pub:  "0x04d1885cf7d4add68ebf2953c1e3406c5ee0e9b9d6c07782ec6bae9fa4bd9335",
			msg:  "I am an OKX address",
			sign: "0xc084c15b8fcfbbf7fcc6b3a8131c885e40eccebfd53472ade25c2ba3d57f2bc1f5089d5a29d68a812477013ea948d06ef77f5894040f7026f323ffd2cdc0640b",
		},
		{
			coin: "SUI",
			addr: "0xdf7c04c9bebf4b35bd8c66a92469f0b66cf77ce8586b6262709a897659f4e772",
			pub:  "4b4c3343b6c3efcd579cd6d6877d5a953166f19825d595296fec109d1f5e0950",
			msg:  "I am an OKX address",
			sign: "0xbbccff4544b0178700273db07846c25541123fbaa70e1e0102b66b0dcfb29f2f9aad23de02195128e0e634a4cab6890bf5db1063d454a69df811040ba298be0f",
		},
		{
			coin: "ADA",
			addr: "addr1v80axvuddhj0mhlyzk9ness59ev752qvpw53fvf5xjy24yqhvxq3j",
			pub:  "0x6adf88e2c6c0773abdf059e16d1800af03a44e6031dbd66e27fa31256db08b69",
			msg:  "I am an OKX address",
			sign: "0xac22bb399d0f93c5e17e43a8f1dd3b145e4c439546b7aa41555f434f304b05020e16ae220dd5721b8cd7f280416217ac5d0905ecb9792c712d13b4ed7b005208",
		},
	}
	for _, tt := range args {
		err := VerifyEd25519Coin(tt.coin, tt.addr, tt.msg, tt.sign, tt.pub)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestVerifyEcdsaCoinSignature(t *testing.T) {
	args := []struct {
		coin string
		addr string
		msg  string
		sign string
	}{
		{
			coin: "FIL",
			addr: "f1lzszobfjwres2otlbitgpbeo6ha72sujwsdjy5i",
			msg:  "hello world",
			sign: "0xcde439bce3471220be8d58eb09a35a8c11320f2cbaeb6714816972b044d059542acbd3618f90cbc5d4baa208105a0c0b0ebe98d4353f4ba1e2abda6a362103c81c",
		},
		{
			coin: "OKT",
			addr: "ex1a0ugda5r0hc3vrzu9wyfkx22vz3g2y2paegqvc",
			msg:  "hello world",
			sign: "0xf4d27cc1407e186ae8cf5c3c4ac8c4bb7d20dda7b5de2b1b212e660cb1115f0b5c8a545ee4b18da7e8e78a6d03e094411ff82f0421dc96e364febdd4bd8b86d41c",
		},
		{
			coin: "OKT",
			addr: "0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547",
			msg:  "hello world",
			sign: "0x0dc53fbecb12f7e14e6eabcb5c9c4e03373318a271d984d56d984ab6b7b9494a73544be662c164490d1d958ae20caf309f0d8003847f1b68c14ead516dfaa83b1c",
		},
		{
			coin: "CFX",
			addr: "cfx:aameksd3gwvmtduc861ym2uzkfaawu9566k0jnte55",
			msg:  "hello world",
			sign: "0x963782c81868cd018211f8cb1ef9eb3a3dc460fd6bb6f6fd46022200c68fdebb3e1125c14379deceb8aca37e94401e0902516d7e7ccde9b5e7c1f40f6d1958a61b",
		},
		{
			coin: "ELF",
			addr: "2XNagboftecQgKRtgG8W5zpdRiZWinUfsfZm62a5NmmEVoZG7X",
			msg:  "hello world",
			sign: "0x933d483f718750d43841cabff4884650221fb0425a72ac98839286927851651645b07275d9f0083622f7925300b0959762c01d8b7610818e89da7315fd0f567d1c",
		},
		{
			coin: "LUNC",
			addr: "terra1hf7afhf4y6wlxqvr7lx4pmct5gunczmnh9emsg",
			msg:  "hello world",
			sign: "0x2fb7d2afb07123b7b7f843d601382f0d13535d3cc620db9bdd062f9d7ed0a6ec00f0fe1a075f882f3679e2c518fb009ebac62cf972f871ad2f40fc2cd85da53f1c",
		},
		{
			coin: "OKB-OKC20",
			addr: "0xeb196a61f9a1e35bf5053b65aaa57c5541dcba86",
			msg:  "I am an OKX address",
			sign: "0xe8df58ec46822f86a0a2fb547260ac55caeeb256916a8c2aabcc01cbdfc13ff264992f2127f3e1cc8e45bf936947c50c8ea097602712e6868526d7fccd9273bc1c",
		},
		{
			coin: "OKT",
			addr: "0x4ce08ffc090f5c54013c62efe30d62e6578e738d",
			msg:  "I am an OKX address",
			sign: "0xa181d622f9a1d1aac327c026a46d11c95a44fb8994a07d232a15b79c12d225a7059db547c0a79b67605f68930d8d0f93a9939589c4fc70041e66322efc61a2421b",
		},
		{
			coin: "LTCK-OKC20",
			addr: "0xeb196a61f9a1e35bf5053b65aaa57c5541dcba86",
			msg:  "I am an OKX address",
			sign: "0xe8df58ec46822f86a0a2fb547260ac55caeeb256916a8c2aabcc01cbdfc13ff264992f2127f3e1cc8e45bf936947c50c8ea097602712e6868526d7fccd9273bc1c",
		},
		{
			coin: "LTCK-OKC20",
			addr: "0x4a11078a99b118bbfee78a5c187d98d264360433",
			msg:  "I am an OKX address",
			sign: "0xee5e609b94842cca32b0f482351bbd2e76e6730dfb146337a077fd03dc6219b62ed1d4e4daf243cb2ba50a12ba47e5f76b8d22e91e8a9c7d228b22a2b236a0ad1b",
		},
	}
	for _, tt := range args {
		err := VerifyEcdsaCoin(tt.coin, tt.addr, tt.msg, tt.sign)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestGuessUtxoCoinAddressType(t *testing.T) {
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
		{
			address: "1GhLyRg4zzFixW3ZY5ViFzT4W5zTT9h7Pc",
			wants:   "P2PKH",
		},
		{
			address: "ecash:qq8d5lh8c78sraajk2ndeqvgqjhdu58zny7etakvlm",
			wants:   "P2PKH",
		},
		{
			address: "ecash:pzgm4hmxk35vkuphlz8v8lprsmppruf2a5l75ru30k",
			wants:   "P2SH",
		},
		{
			address: "bitcoincash:qq8d5lh8c78sraajk2ndeqvgqjhdu58zny7etakvlm",
			wants:   "P2PKH",
		},
		{
			address: "bitcoincash:pzgm4hmxk35vkuphlz8v8lprsmppruf2a5l75ru30k",
			wants:   "P2SH",
		},
	}
	for _, tt := range args {
		res := GuessUtxoCoinAddressType(tt.address)
		if tt.wants != res {
			t.Errorf("Want %s, Got %s", tt.wants, res)
		}
	}
}

func TestVerifyStarkCoinSignature(t *testing.T) {
	args := []struct {
		coin string
		addr string
		msg  string
		sign string
		pub  string
	}{
		{
			coin: "STARKNET",
			addr: "0x0668a58aec3151bdf0f6f9360ea89453733d5c65f45cafe2d1591d8cddbd2395",
			msg:  "hello world",
			sign: "07abc5982853352d98763940726e0ab31f8c25fa6f1fc9951edad49cbfab6297018d6b42f497cefdcdb70aab45604612f03598a085bf984778fa413c7fab1b75",
			pub:  "0x346262ffa4ec2f40feb9ae81e416af7cca9fcfa8871f1f9169e6dccd63aa667",
		},
	}
	for _, tt := range args {
		err := VerifyStarkCoin(tt.coin, tt.addr, tt.msg, tt.sign, tt.pub)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestVerifyTemplate(t *testing.T) {
	cases := []struct {
		verifyType string
		coin       string
		addr       string
		msg        string
		sign       string
		pubKey     string
		// only for utxo
		sign1  string
		sign2  string
		script string
	}{
		{
			verifyType: "evm",
			coin:       "ETH",
			addr:       "0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0xe2d351cd4617fea1cf3ad6b907f7897384b6a01cdd7db3fb23c1967bc40cae527e95690bebf9527e0882ffad3bac76e3424c91378a019997d970cd63935200861c",
		},
		{
			verifyType: "trx",
			coin:       "TRX",
			addr:       "TEjxQjU3CxkFrSDcPfHwZXSuPpCpdQ27NJ",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0xb52562f78a945b7a1fcfbbc4bc19c3dc80f0ed10ba633f37266ad25dbe9100910e1a5712f2c90e61b05d10ba7c43d0fbba2a06c48d9b7425b8229c3ea617ffa31b",
		},
		{
			verifyType: "ecdsa",
			coin:       "OKT",
			addr:       "0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x7a819fef9e4383bf165b3173da04e0b3fde7eb501f8711f29bc9a14d90a37ba337e394f9544d4d90ee65d19cfbd437b038be4d2a4aadf9da3973f9eacbc6b2001c",
		},
		{
			verifyType: "ecdsa",
			coin:       "FIL",
			addr:       "f1lzszobfjwres2otlbitgpbeo6ha72sujwsdjy5i",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0xbc7cbd93ae5231a94f5289cef2d161d134994a2fc626ab1863ec5a10ec89dc902f0cf695fd99308e402185d4184fc5a002074d6a31b6666b09634c0df85220b31c",
		},
		{
			verifyType: "ecdsa",
			coin:       "CFX",
			addr:       "cfx:aameksd3gwvmtduc861ym2uzkfaawu9566k0jnte55",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x3201f07c8762b1336c58c13f586ce329a53bf59c493c80cde5854d02dda16b4109aa7e4437671d53f72bf6e8340b9370fa863cd505f845a0c964eda7e8696fed1b",
		},
		{
			verifyType: "edd25519",
			coin:       "SUI",
			addr:       "0x4c6072a97c304be864d9e187ded3ea33e83e447672054dd69cc1e2de7914b250",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x870faf27df65e6b0a8f397b1461374d4b1843f53c0e71407316f8b0bba3e065a57742fc59924caec8f0f30be0d95a2051a73c8f0e54ba21ec1a4b8e7b7fa6e02",
			pubKey:     "0xd035239ebfe6e6192249a8c028cdc845aed0a469f119e200c7e2f306e5067094",
		},
		{
			verifyType: "edd25519",
			coin:       "APTOS",
			addr:       "0x327dd297dfacf7c2d8207aaa23c0f0e8bcaf4c1612febbf63b9f7376810b8ec8",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0xe4d968ddf33366e3e085294bf80077fdb2e89243f4e4f41779e531640504ae1e3b52f5d3e1a1c8fd838a4ce07bad511d7d9438ea7515fe9629d6ec81124cce01",
			pubKey:     "0x61f579fc779146304353027b425a216d8015889c5f3b715ad26135b862f3bf84",
		},
		{
			verifyType: "edd25519",
			coin:       "SOL",
			addr:       "7bzoTJhZmpU1vQVjN63fQ3iVYmWCVgQh1sYSqsjuapU9",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x81c709d6f44647a7e34c318e63c4f7e34617fa0420c9b58b071adfac3f7fb3e2ed7da6fe0b59dc7318cf471bbbe0abc873362b7cce2f269ef1f7e9ecd5716704",
			pubKey:     "0x621d398b19304995ee140c21afc544d62382d387b5c08dfd096b475a304339ea",
		},
		{
			verifyType: "edd25519",
			coin:       "TONCOIN-NEW",
			addr:       "EQA5rifVSCc8qQfpCXvq4zJGJPsA0EPCDoWdtg234INftsWj",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x48d46ee1557c4850dff0622cb3b960a667e72df10868c905df1c17bd1dc8ebc42691bb8df2814c56933cdfc67a6ba7626ae5faf4ee6190c4a030a43179e0c502",
			pubKey:     "0x3d2696e3d5cbc9047b338e6a56552db1d43ca6e063bc7aa667b18005984372d2",
		},
		{
			verifyType: "ecdsaPub",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x3a7c6d3f6e6885c07d21bb9eeb266cdf1400ab5c1379664ca9c74c6eaf4dcbdd03895a64253f7c5da9da504548dc87bc634e571042bc74bcfb18a005d3c592531b",
			pubKey:     "0243eb3f77eac97882ca3a6c038830e199f1b02cbab4f8664fd08c54911b0a05e6",
		},
		{
			verifyType: "edd25519",
			coin:       "PI",
			addr:       "GDMA3OUZHWTQ4RZBYEXR4VEZ7YGQNXOAMRSYRHDM5GSXH3UZQF5VGPXN",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x71e6adb98dce49b523594f3bc0f96ab7eda84c00ef3ddc7e1eff8d18c5dd537ae409f6cc7c0f9d2f98d2bcfeb2885c5f3856bdf867ee9b68a5e16c12be7be50d",
			pubKey:     "0xd80dba993da70e4721c12f1e5499fe0d06ddc06465889c6ce9a573ee99817b53",
		},
		{
			verifyType: "edd25519",
			coin:       "ADA",
			addr:       "addr1vykqe7mmq83jsfxu0relg46vasd5d4gr4yg5pchg5v8tuvck0jkfn",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0xefb9a223ad794dac0a791caa29909603cca6c780a0e184c69e1a13a9a4c1c0b60ae83671de4178ab7ba83131bda9210d84ffe2a4556d1891add02fdcd7181a06",
			pubKey:     "0x7afe9bac605a184fe32d9dfb582b9d837841f5ff99916b73dd3f705d19bbba50",
		},
		{
			verifyType: "edd25519",
			coin:       "DOT",
			addr:       "15WXogcgXnHsZ1FeuNc6cg34i8R6JCXgWrDLNLERLLesJ7bf",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x2d6149f4facd2966a3e8c229e2fad837ad35fe017d1a64c53464c325fd0fcb9726bef8f21004b689cd884f3749ce07ef08f4dceff02fd9aa96d48bd466985400",
			pubKey:     "0xc776bfbeeeb0b1ddd1ce6cccf55ce795f5306bf63de37d72e5af50b3be23ce49",
		},
		{
			verifyType: "edd25519",
			coin:       "NEAR",
			addr:       "b1e2af21d50c8940aaebe9650bd3e09eb49f3945c6ec2d917de0d64cdfaf3fba",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x0583df7d07e242e1240a97b6322bfa1df7203a09d10e8a219442ebb17064ae978024df47a58551f61afd9b10774688c0c217aa319eddeb9285b252006edb3c07",
			pubKey:     "0xb1e2af21d50c8940aaebe9650bd3e09eb49f3945c6ec2d917de0d64cdfaf3fba",
		},
		{
			verifyType: "ecdsaPub",
			coin:       "AVAX-X",
			addr:       "X-avax1v88xtvkfmas5467c0zcdhsdxv42nm8d0gjmp5q",
			msg:        "OKC_DTT_AUP2025",
			sign:       "0x80c84e2a993899f12882d4319061b8ce1d7319742e4939cc3d93caca349f3b0247f2a3d5f9bc5063a62d2cb168ca43d1cb788e5d7f7e1770fb04d71f7649ed501b",
			pubKey:     "022e6cb3bc597b11f64aa46b40a6c6aa4853557ef66b53c6a49a9540b25e0afe1b",
		},
		{
			verifyType: "utxo",
			coin:       "BTC",
			addr:       "393maTY7rQScy4SmYE1XSXUSgK73byhgfA",
			msg:        "OKC_DTT_AUP2025",
			sign1:      "IDRyTg89A1N6quCCnj3jLht3Dboar6rSh/7YlhxbVN1RfVOYti0fyRcWXplLQjbPgDChthLssytrbDFVDcLA750=",
			sign2:      "IJ+Wf/va772DV2ETAGxN7mo2VECUM9kQBMpUNvAJBnrgM8/o2C3+ID7Z6E5QZTS+VLCZ8Bgtym0rxcLfFLl0B+s=",
			script:     "5221027adce0bd3080066ab90c68199ff73128b3ff8c847d15d9e4c6e88fb4c6e6486b210273fa0df3ffceeda23b0074d9fe83d9ee3a209fad6e4546fdec5ede39abcbb70d2102a38ce748c5a1e1889f0d72ecd6f2130f5f73a11e01fff9f0d22796e40217571953ae",
		},
		{
			verifyType: "utxo",
			coin:       "BSV",
			addr:       "1Hgc1DnWHwfXFxruej4H5g3ThsCzUEwLdD",
			msg:        "OKC_DTT_AUP2025",
			sign1:      "IJ+Wf/va772DV2ETAGxN7mo2VECUM9kQBMpUNvAJBnrgM8/o2C3+ID7Z6E5QZTS+VLCZ8Bgtym0rxcLfFLl0B+s=",
			sign2:      "",
			script:     "76a914b6ff746a94f1ed230af7c891bc68a8dcc45bdf4f88ac",
		},
		{
			verifyType: "utxo",
			coin:       "DOGE",
			addr:       "9yo2KJc1vUKWsRpExMfwgf6pNtV5gG3vqW",
			msg:        "OKC_DTT_AUP2025",
			sign1:      "H7UaLfjDn/VCq9dXpaQ7c5bwvrLi3D+3+wh1xBu7jKkrInQosKZC6c0oaR9Kl01T0WIm70h9rA2KkyivTZB4fFw=",
			sign2:      "IBRzablLv1KsJCyZnGE1DZLUP+o2gghRspnZbzmE56egf64Od0MmvFdBtfqajGyLoZwhkXjp62nm+I5U4LtJg5w=",
			script:     "5221027adce0bd3080066ab90c68199ff73128b3ff8c847d15d9e4c6e88fb4c6e6486b210273fa0df3ffceeda23b0074d9fe83d9ee3a209fad6e4546fdec5ede39abcbb70d2102a38ce748c5a1e1889f0d72ecd6f2130f5f73a11e01fff9f0d22796e40217571953ae",
		},
		{
			verifyType: "utxo",
			coin:       "LTC",
			addr:       "MFFutLx5oXJ3mZife6zsGAir11hVbmUF99",
			msg:        "OKC_DTT_AUP2025",
			sign1:      "HzQdk7Jt9N4J6xQxFQSF1Ts/m5cNLVIgYgv75eg8KjZ9fyxQgpS3f8Aywp8yd5g/IpAe36eMgBNadSbQkjATnmI=",
			sign2:      "H1Ni/G3cPbjLTzn5Bmm1WDu5wPKLwCmzKzXVzwzz3EJCOREBWusePO3ZYvAoVALuVBurKHhWSLTdO6sZw5A43oM=",
			script:     "5221027adce0bd3080066ab90c68199ff73128b3ff8c847d15d9e4c6e88fb4c6e6486b210273fa0df3ffceeda23b0074d9fe83d9ee3a209fad6e4546fdec5ede39abcbb70d2102a38ce748c5a1e1889f0d72ecd6f2130f5f73a11e01fff9f0d22796e40217571953ae",
		},
	}

	for _, tt := range cases {
		t.Run(tt.coin, func(t *testing.T) {
			var err error
			switch tt.verifyType {
			case "evm":
				err = VerifyEvmCoin(tt.coin, tt.addr, tt.msg, tt.sign)
			case "ecdsa":
				err = VerifyEcdsaCoin(tt.coin, tt.addr, tt.msg, tt.sign)
			case "trx":
				err = VerifyTRX(tt.addr, tt.msg, tt.sign)
			case "edd25519":
				err = VerifyEd25519Coin(tt.coin, tt.addr, tt.msg, tt.sign, tt.pubKey)
			case "ecdsaPub":
				err = VerifyEcdsaCoinWithPub(tt.msg, tt.sign, tt.pubKey)
			case "utxo":
				err = VerifyUtxoCoin(tt.coin, tt.addr, tt.msg, tt.sign1, tt.sign2, tt.script)
			}
			if err != nil {
				t.Errorf("%s verify failed: %v", tt.coin, err)
			}
		})
	}
}
