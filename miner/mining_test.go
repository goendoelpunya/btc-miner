package miner

import (
	"btcnetwork/block"
	"btcnetwork/common"
	"btcnetwork/storage"
	"encoding/hex"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestMine(t *testing.T) {
	var (
		version      = uint32(0x20000000) //0x20000000
		preBlockHash = "00000000690b85ab50322252a64d2cc0262e6fd282e81e293e1b88614ee8bfd2"
		merkleRoot   = "33e0cda9e3ccf618b672e85eed0d2af92bc33a80504702df226822d87883e58e"
		timestamp    = uint32(1590280470)
		bits         = uint32(486604799)
		wantNonce    = uint32(1867461409)
	)

	txids := []string{
		"b96667b89cc035f0d07edce0a2df5a7cfa6aedd21206c35403854b6c05a6fd7f",
		"c8c5b519188bead8177cdb7186d94b731b7d054444883a2449f3e99739dc79ac",
		"b5d0b90a30cf3b022cdf72950941f1993e795bc8901b0c203f8d130a5a15eb76",
		"02da2fd57c8ccf46f1017497fd21d31a9a3da371e22551c191ac5d983f9fe83c",
		"3ef2774594ae1272d5ed44f5629bb3584fa5957f216b1e321dfe711c2f6f02fd",
		"1ac096435bb2ad795f4dffe9093d40645b4adab9d4fa96c3cd8eba52c04c9760",
		"3e98b1bddc5a1d9a5be07b2fdee069958600febb669c8524f858f58926bb8300",
		"f365c162884b722a58f06a76ceb9bffde4372a0c7c630e403479561ee2552604",
		"8cbfe031adb3d05b16dd74f1aaddc213ded417cc9ae498550d0e8974bba4ac86",
		"7e383ab294ca9312232fa0d1db3247f0aff8e9b3882522d421ed91f7349f210c",
		"d6e47a57e0955bf0f40b2ca929bf63f803b768e61001bffff6961c89aadb34f8",
		"bfa9e2015b84aa29a471e88666ef3278cf5f4ab6c16c19d1d25a188cf030560d",
		"fae4ab3f1d7b93162e2e83320b93d41684e2054f0df5b2e0d7ee9301a0b65d1c",
		"d59ce999a8b4c5312055c2e7265365046f74309f5787918475ec3738f5454543",
		"95ed32254c63220255f22f215bffb33f5b4a81534c3980b55d3c3fca0d6c1110",
		"5649b629da16ca90de6496d0deaa0ab25f58d9b4d3fc67f38c645ec42fdef3da",
		"80cea25e7c43eb7c1b4ba492669a726615b32cf89010889523fe72f76dfce92b",
		"a8d558d88ea94a37328587e9fa7f1a77ff5b27b61f944912e43af087d30c8efe",
		"3e1e2f5990a14c766aabcf051796ac8a23fcbdd7f63fee1bf88641434ab4805f",
		"59f9a73537cb3a0b7c4e7ad8c6f99a17023e4592709c5186a8e31743075721f7",
		"05417dfcfc276233555d1eae574275ded6fef22009e8c29a1a8588ec430bac42",
		"294246885b9b7c69f0d8890552201c50ba403dc0c735e8570016ecca946c3d7d",
		"9a35251f24b3a37e7ef530dd5b29b29b92635074fe6deb49bb11330bdb14d246",
		"f3aebfb8a519f91918e37501df9d1cae5db598418d213a93fdfbc0dfd8b8415b",
		"e7cee2e263040ed6898780f506210ba1f0a69c8728ddf53c4ddde045f343055d",
		"7dfd922de85f7eca6aa329f52b0c406e96d6f5988ab88619cbfa51c6308b666a",
		"7c82f2390dc4a1c780f575da66164a57d587e495da18974dce053f95de33385e",
		"32a2fbb54eda5cc1e0e612eba1420dc95e7384f58453c769eaadb96381c0c561",
		"5ed5929528bd09a3f59100c5fd434474381805f5d4a6004ec539296f56c73ab3",
		"e3f94b8e55ea3a5c1e40d571db686f9f7dbab491a29297042a0e566128cc39da",
		"fba9cfd68791f36692f39bda3a90303463ae409331b21fd102c3af21fae15366",
		"18d60717e56d2fe195ae5f7d12230d18635a938488fcb1af6c6c2b892db5216c",
		"27fd8e719c092d7d460729483acb7405918486ca02201918d770828218640db8",
		"f62a86e1ee9ec468955e543830c6cbb34b4ec0518535b3e22211689a1471f1b1",
		"88a67fc7f3ee86839f84e652f7522659c976730c0a29004c429c9c8169c5278c",
		"f6d45fa4ad52ba504d29ce43816cfe2119d620bd36c6ba9d7de499f3098c1790",
		"c927f729a759f04b42a91f37dcacf2e5ea17f757ab24a16048ca5a87992e3293",
		"c2ea3a25ab7d4c6edf313ac8d09b4dc7735d32a3993aa1407cd7287eaa9395c3",
		"3ec0dbf44aaa960ed347f36f7350d2989e949d656ff80296ff4e188ad3324f9b",
		"b93f453a941f0f59bf24835aff386eb346c40dd5ccd9d9728dc15bdf8c3575e9",
		"640c72fac4eada1f5e65d023573e119a737a87c80991eade566deaa4b574e8a2",
		"4914c3d72105afd81f7b92976de160b61d60f226e3c1af867e9d9026daef93ce",
		"7b39b04f9c936da476c2d7452ac04e2d1a58fd6e5ed4c8b4c96a25bc9037b4d3",
		"1b35d5a3aa03dabfc42c1185171111cddbf3f892639e7486d3a077e0cabdd7d4",
		"ec029723ac3b93fd573c818b84400019950674aadef720d133776fce7f3c0bbd",
		"ef484f8c6c7b931c98055f7325c7c1203a853110738b9fe745b3b1fdeb94cc58",
		"d013e8c6226ff57e1b8e860046746696c83447902531590f2c3e7bba2259195f",
		"603a4283d8144935e8b99293dffbabb2f0323fa9543e65e23724e2d04a9c8356",
		"df850b9088e39966aae2e00c68a62fa4e9ccb97e4c4e563ddf8703685120db65",
		"05fe259557b6ff037f386a3750ccb34fe3cc2a68d825a2d93671d080f6012e40",
		"722e706985bc81ab8abf2789264a719fbb204f4a35efe7dcbc9d3556fb80e419",
		"e1bb31aa6c91f971ba2a0b5ae50713927c07cd4b5f2474846fee7035937b3b7f",
		"287854205d38d145f24d740eab5d291410aa71fd62b3538f58add09293a67274",
		"d7eabe8c35db0a70a9f5dccb6b8ed76fd4de561933baf4ca7ff41e95dd766028",
		"63443075e5359cc73668e80750522e34f91c0f32261265e7d913820f7f051643",
		"2197ca6b336acadb30b36957aa6a2a39508ca1d9e74dddd8a6fe852d4825e300",
		"8b53c54cffda0ed0b5729e465a270c5f43891c3adc65a0cdba31147fafce9609",
		"c30c7a2832b6749be6f7a7299b5777e50c48d819e4d7cd953e048c85590c5934",
		"bb7136aed1296514fe8ee063e134c898368395b07cf07bcee1316e792c9e613b",
		"7daf28b3f821538a240d985b99fdb1eea57b098d1a947a08dbfc2ac1397b3e41",
		"573831d8a6a75aa858a3042bfdd64bff66a9abbd4d21fdd735ff4113dd8a3f47",
		"2e3b1cc43190d3eb7d07766253ca1b5007b4b9f2718a4323177cadc3940e8a84",
		"11780fe799ce7e6f9cac31fc55ef4d9edb66daaf0fe525b1c9876ae848a81f4f",
		"2451ec13308fde72211dcec01ba8c8f0c35db7b78683b70cf6187ded4e30726a",
		"0c6558ac3c6ddc3c12057bf485ab60520e857ae5aac9645a1ffb5a2844f7907e",
		"4463fa4290de39f1a96ae298a77247804588d3032ad56a33a144a54e623fe78e",
		"e3296fc5d9f5a0aa60b7edfdfe3f3636e3e50fa7b273746f06460a4f1938947e",
		"c811702222e551e0db28a482a8d172e2251307a759654bf52cbf99ec6829409c",
		"f7c2d12c738217e1513b924912e1743c19ecd5b034db27c4adbf089bc13046b2",
		"a6d5d5e84bbd5fff4ec58769eb46197e0cb5261792dbc4e1d6957fa9fd62f0b4",
		"f28e6906bcb20c04083935a30a438c29c5952e26c5b859a77c6737da3698cccd",
		"af17ca7b91ab57f6231802d42f7a4fccc82ef58005280e32c746b412d5c77ad4",
		"06aa3b281cff48c7d041e441e490847fadcf94c0908eb6a9f8cd9b8f00038bd4",
		"1b7a35072751393e44ea426525bcc4392fc7d7ec491e96baaae1c8f39fe133dd",
		"f78a29da4817763d026968470a1f44a0003770debc00a38ba4406abb9cc2c3ef",
		"6cc39611c035307456035fbe268a1d26636de93320825afa985d090a4366b5b6",
		"c097b0784d2349dd1d99cc5e88800ac32ac199f19daee07f575b19dc93e37dfb",
		"23c067e18453f21d7b1e8cca5e6efe25734beb1cab58c0ff15bae2dfbdcbfb59",
		"8ddb834abf6210f390988af6d4d0a69f2c0ab300a94fcf48f44f662c68916e03",
		"a68b15a24a160cf2852e714a23d6034997ba71ac08a18d34bcb67796a39b2f0e",
		"6e8acbe33f5dcc9faecc682610ecef694c1e2cf860ed924eed93aed09dacac15",
		"f526facefadba15fb2fd43faa8656f30b9306775cb14a87a8e46caf9a077191b",
		"8d6fd2de22d4dc73d1cc805032de63368077220b4ad268974c6c3e85cf09f031",
		"d5ab9e7ec395818df9b8ec851abdf8716bfd1af636702283cd9e7a33fb175d32",
		"ea87edd68aa3ec065cc14f972f29b4a50700114650c2ee69d74a43e289aafe4b",
		"355e194f03fa5e8998e39716751aa34b4d3f9c7fcbf8e0fbbb91ef004969244c",
		"fb8d0167b3c0bdc988931e8c2e7e099d967b5e523e3f88c5946027d9ec33c06b",
		"e6d0d3f1ff65b38b9ae2f4fa88a2bda03bd767a332313de35819a88f4f878570",
		"a0f7c4dde1a83fa2cd759cf2908e06824e3b824f08c3e6815e24c2ce444f907f",
		"cdfa23c378bc7e1b3e2135f3ac96f03756336e6d13bcd73c1fa67c7327239b7f",
		"866cac79865069c0bde2f57becc7463ad2222ef4ee5b89b04659a73a597f098e",
		"67c5a7aee0dea7367b8691a3f62fb60f40d606ec19bc55ce260bb83ddf68db93",
		"67c9033445b8a3c6b843c806d57466979a5b01ecbc5ea3ec79271787592ef595",
		"30961c096a610a4a166e7e98b6c2a7d939eecfff282f2be95cebc06f6893379a",
		"9891f82c0179e4eed33b0ee7c01f8671f84f340ebe1cb4fdc512076f4a1703a7",
		"2f58c8fbc28ab7ab8d34a956f8a36c94f4645bbab360ee653c7cfb18468235e4",
		"a15a00bb6dd649fe1e4dfbbf513884ae6b06168d2f1494ae529715d792dfaaef",
		"ac5debb733982803e001cffd51a83d7f485fe5c1a0908e9c26b901d0607b6ff7",
		"45f8c51bb54e5fa069adf48987baaa67f60f3b9f1f96ffb12b2b1029fb4517fe",
	}

	gotMerkleRoot, err := block.ConstructMerkleRoot(txids)
	if err != nil {
		t.Error(err)
		return
	}
	if gotMerkleRoot.Value != merkleRoot {
		t.Error("merkle root error")
		t.Error("want: ", merkleRoot)
		t.Error("got: ", gotMerkleRoot)
		return
	}

	//求区块头hash
	var header = block.Header{
		BlockVersion: int32(version),
		Timestamp:    timestamp,
		Bits:         bits,
	}
	var buf []byte
	if buf, err = hex.DecodeString(preBlockHash); err != nil {
		t.Error(err)
		return
	}

	buf = common.ReverseBytes(buf)
	copy(header.PreHash[:], buf)

	if buf, err = hex.DecodeString(merkleRoot); err != nil {
		t.Error(err)
		return
	}
	buf = common.ReverseBytes(buf)
	copy(header.MerkleRootHash[:], buf)

	var gotNonce uint32
	if gotNonce, err = mine(&header, wantNonce-1000); err != nil {
		t.Error(err)
		return
	}
	if gotNonce != wantNonce {
		t.Error("nonce error")
		t.Error("want: ", wantNonce)
		t.Error("got:  ", gotNonce)
		return
	}
}

func TestMining(t *testing.T) {
	cfg := &common.Config{}
	cfg.DataDir = "F:/go/src/btcnetwork/data"
	cfg.MinerAddr = "SQmHEbXs5qhDt5mqeibX6MJnKpitfz9EHQ"
	storage.Start(cfg)
	Start(cfg)
	Mining()
	storage.Stop()
}

func TestInteger2bytes(t *testing.T) {
	var testcases = []struct {
		data int32
		arr  []byte
	}{
		{1, []byte{0x1}},
		{255, []byte{0xff}},
		{365, []byte{0x6d, 0x01}},
		{65533, []byte{0xfd, 0xff}},
		{38557882, []byte{0xba, 0x58, 0x4c, 0x02}},
		{616926126, []byte{0xae, 0x8b, 0xc5, 0x24}},
	}

	for _, oneCase := range testcases {
		buf := Integer2bytes(oneCase.data)
		if !reflect.DeepEqual(buf, oneCase.arr) {
			t.Error("wrong convert")
			t.Error("data:", oneCase.data)
			t.Errorf("want:%v,get:%v", oneCase.arr, buf)
			return
		}
	}
}

func TestBigIn(t *testing.T) {
	var data [8]byte
	data[7] = 0x01
	bigNum := new(big.Int).SetBytes(data[:])
	t.Log(bigNum.Uint64())
	t.Log(hex.EncodeToString(bigNum.Bytes()))
}

func TestLoop(t *testing.T) {
	var i uint32
	for i = uint32(0); i != math.MaxUint32; i++ {
		//common.Sha256AfterSha256()
	}
	t.Log(i)
}
