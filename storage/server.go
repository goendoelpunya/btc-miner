package storage

import (
	"btcnetwork/common"
	"btcnetwork/p2p"
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
)

var (
	ErrBlockNotFound          = errors.New("block not found")
	ErrBlockUnserializeFailed = errors.New("block unserialize failed")
	ErrTxNotFound             = errors.New("tx not found")
	ErrUtxoNotFound           = errors.New("UTXO not found")
	stop                      = false
)

func Store(newBlock *p2p.BlockPayload) {
	//防止向已经关闭的newBlock通道写入数据
	if !stop {
		defaultBlockMgr.newBlock <- *newBlock
	}
}

func StoreSync(newBlock *p2p.BlockPayload) error {
	return defaultBlockMgr.updateDBs(newBlock)
}

func Start(cfg *common.Config) {
	startBlockMgr(cfg)
	startTxMgr(cfg)
	startUtxoMgr(cfg)
}

func Stop() {
	stop = true
	stopBlockMgr()
	stopTxMgr()
	stopUtxoMgr()
}

func BlockFromHash(hash [32]byte) (*p2p.BlockPayload, error) {
	log.Debug(hex.EncodeToString(hash[:]))
	buf, err := defaultBlockMgr.DBhash2block.Get(hash[:], nil)
	if err != nil {
		log.Error(err)
		return nil, ErrBlockNotFound
	}
	blk := p2p.BlockPayload{}
	if err = blk.Parse(buf); err != nil {
		log.Error(err)
		return nil, ErrBlockUnserializeFailed
	}
	return &blk, nil
}

func HasBlockHash(hash [32]byte) bool {
	has, err := defaultBlockMgr.DBhash2block.Has(hash[:], nil)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return has
}

func LatestBlockHeight() uint32 {
	if defaultBlockMgr.IsEmpty() { //如果是空的就返回创世区块哈希
		return 0
	}
	var buf []byte
	var err error
	if buf, err = defaultBlockMgr.DBlatestblock.Get(LatestBlockKey, nil); err != nil {
		log.Error(err)
		panic(err)
	}
	return binary.LittleEndian.Uint32(buf)
}

func LatestBlockHash() [32]byte {
	if defaultBlockMgr.IsEmpty() { //如果是空的就返回创世区块哈希
		return defaultBlockMgr.genesisBlockHash()
	}
	var buf []byte
	var err error
	if buf, err = defaultBlockMgr.DBlatestblock.Get(LatestBlockKey, nil); err != nil {
		log.Error(err)
		panic(err)
	}
	var hash [32]byte
	if buf, err = defaultBlockMgr.DBheight2hash.Get(buf, nil); err != nil {
		log.Error(err)
		panic(err)
	}
	copy(hash[:], buf)
	return hash
}

// todo:根据区块高度找出区块数据
func BlockFromHeight(hash [32]byte) (*p2p.BlockPayload, error) {
	return nil, ErrBlockNotFound
}

// todo:根据交易交易id找出交易数据
func Tx(txid [32]byte) (*p2p.TxPayload, error) {
	return nil, ErrTxNotFound
}

// 根据PreOut组成的key找出交易输出数据
//func Utxo(key [36]byte) (*p2p.TxOutput, error) {
//	txout,err := utxo(key)
//	if err != nil {
//		log.Error(err)
//		return nil,ErrUtxoNotFound
//	}
//	return txout, nil
//}