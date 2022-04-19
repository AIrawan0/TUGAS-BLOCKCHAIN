//Paket konsensus ngeimplementasikan berbagai engine konsesnus Etherium
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

//ChainHeaderReader mendefinisikan suatu koleksi kecil dari metode yang dibutuhkan untuk mengakses lokal
// blockchain saat verifikasi header
type ChainHeaderReader interface {
	//Condig mengambil chain konfigurasi blockchain
	Config() *params.ChainConfig

	// CurrentHeader mengambil header saat ini dari local chain.
	CurrentHeader() *types.Header

	// GetHeader mengambil block header dari database dengan hash dan angka.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber mengambil block header dari database dengan angka.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash mengambil block header dari database dengan hash miliknya.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd mengambil kesulitan total dari database dengan hash dan angka.
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader mendefinisikan suatu koleksi kecil dari metode yang dibutuhkan untuk mengakses lokal
// blockchain disaat header and/or uncle verification.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock mengambil block dari database dengan hash dan angka.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine adalah algoritma agnostic consensus engine.
type Engine interface {
	// Author mengambil alamat Ethereum akun yang mencetak block yang diberikan
	// yang mungkin akan berbeda dari coinbase header jika engine consensus didasarkan pada signature.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader memeriksa apakah header bersesai dengan aturan consensus yang diberikan engine
	// Memverifikasi tanda dapat dilakukan secara optimal disini, atau dapat dengan tegas menggunakan metode VeryfySeal.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders cukup mirip VerifyHeader, namun memverifikasi sekumpulan headers secara bersamaan.
	// Metodenya mengembalikan channel yang keluar untuk digagalkan operasi dan chanel hasil
	// untuk mengambil verifikasi async.
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles memverifikasi yang diberikan oleh uncle block yang disesuaikan dengan
	// aturan konsesus yang diberikan engine
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Menyiapkan menginisasi lapangan konsesus dari block header sesuai dengan
	// aturan dari engine tertentu. Berubahannya danat dieksekusi inLine.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Mengakhiri dengan menjalankan post transaksi transaksi yang termodifikasi(cnth block rewards)
	// Tetapi tidak menggambungkan blocknya.
	//
	// Note: Block headernya dan kondisi database mungkin tidak terperbaharui untuk mengikuti
	// aturan consensus mana saja yang dapat terselesaikan.(cnth block rewards).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble menjalankan post-transaction manasaja dalam kondisi termodifikasi(cnth. block
	// rewards) dan mennggambungkan block terakhir.
	//
	// Note: Block headernya dan kondisi database mungkin terperbaharui untuk mengikuti
	// aturan consensus mana saja yang dapat terselesaikan.(cnth block rewards)..
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal membuat permintaan segel aru untuk input block yang riberikan dan mendorong
	// hasil kepada chanel yang diberikan.
	//
	// Note: Metode return immmediately and will akan mengirimkan hasil async.
	// Hasil lebih dari satu mungkin juga akan dikembalikan tergantung pada algoritma consesnus.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash mengembalikan hash dari block yang lebih dahulu dari yang tersegel.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty adalah tingkat kesulitan dalam menyesyaikan algoritma. Ia mengembalikan tingkat kesulitan
	// yang harusnya dimiliki oleh block baru.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs mengembalikan atau memberikan keluaran dari RPC API yang disediakan engine consensus.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close Menutup thread background manapun yang dipertahankan oleh engine consesnsur.
	Close() error
}

// PoW adalag engine consesnus yang didasarkan pada proof-of-work.
type PoW interface {
	Engine

	// Hashrate memberikan keluarakn akan hashrate dari PoW engien consesnsus saat ini.
	Hashrate() float64
}