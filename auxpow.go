package btcwire

import "io"
import "fmt"

// TODO
const MaxCoinbaseTxSize = 100000
const MaxBranchHashes = 32
const MaxBranchSize = 4 + MaxBranchHashes*HashSize
const MaxAuxPowSize = MaxCoinbaseTxSize + HashSize + MaxBranchSize*2 + MaxBlockHeaderHeaderPayload

type MerkleBranch struct {
	Hashes   []ShaHash
	SideMask uint32
}

func (mb *MerkleBranch) Size() uint {
	return uint(len(mb.Hashes))
}

func (mb *MerkleBranch) BtcEncode(w io.Writer, pver uint32) error {
	var err error

	err = writeVarInt(w, pver, uint64(len(mb.Hashes)))
	if err != nil {
		return err
	}

	for i := range mb.Hashes {
		err = writeElement(w, &mb.Hashes[i])
		if err != nil {
			return err
		}
	}

	err = writeElement(w, mb.SideMask)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MerkleBranch) Serialize(w io.Writer) error {
	return mb.BtcEncode(w, 0)
}

func (mb *MerkleBranch) BtcDecode(r io.Reader, pver uint32) error {
	n, err := readVarInt(r, pver)
	if err != nil {
		return err
	}

	if n > 0x02000000 {
		return fmt.Errorf("size too large")
	}

	mb.Hashes = make([]ShaHash, n)

	for i := uint64(0); i < n; i++ {
		err = readElement(r, &mb.Hashes[i])
		if err != nil {
			return err
		}
	}

	err = readElement(r, &mb.SideMask)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MerkleBranch) Deserialize(r io.Reader) error {
	return mb.BtcDecode(r, 0)
}

func (mb *MerkleBranch) SerializeSize() int {
	n := VarIntSerializeSize(uint64(len(mb.Hashes))) + HashSize*len(mb.Hashes) + 4
	return n
}

func reverseBytes(b []byte) {
	L := len(b)
	for i := 0; i < L/2; i++ {
		b[i], b[L-i-1] = b[L-i-1], b[i]
	}
}

// Determine the root hash for the Merkle tree formed from the Merkle branch
// and the component hash specified.
func (mb *MerkleBranch) DetermineRoot(component *ShaHash) (h *ShaHash, err error) {
	//log.Printf("MerkleBranch: DetermineRoot (component=%s)", component.String())
	//log.Printf("MerkleBranch contains %d hashes (0x%08x):", len(mb.Hashes), mb.SideMask)

	m := mb.SideMask
	h = component
	hbuf := make([]byte, HashSize*2)

	if component == nil {
		panic("component must be specified")
	}

	for i := range mb.Hashes {
		//log.Printf("  %s", mb.Hashes[i].String())

		if (m & 1) != 0 {
			copy(hbuf[0:HashSize], mb.Hashes[i][:])
			copy(hbuf[HashSize:HashSize*2], h[:])
		} else {
			copy(hbuf[0:HashSize], h[:])
			copy(hbuf[HashSize:HashSize*2], mb.Hashes[i][:])
		}
		h, err = NewShaHash(DoubleSha256(hbuf))
		if err != nil {
			return
		}
		m = m >> 1
	}

	return
}

func (mb *MerkleBranch) HasRoot(component *ShaHash, root *ShaHash) bool {
	r, err := mb.DetermineRoot(component)
	if err != nil {
		return false
	}
	return r.IsEqual(root)
}

type AuxPowHeader struct {
	CoinbaseTx        MsgTx
	ParentBlockHash   ShaHash // vestigal
	CoinbaseBranch    MerkleBranch
	BlockChainBranch  MerkleBranch
	ParentBlockHeader BlockHeaderHeader
}

// CAuxPow
//   CMerkleTx
//     CTransaction
//       nVersion       int
//       vin            vector<CTxIn>
//       vout           vector<CTxOut>
//       nLockTime      unsigned int
//     hashBlock        uint256
//     vMerkleBranch    vector<uint256>
//     nIndex           int
//   vChainMerkleBranch vector<uint256>  // } These are the Merkle branch?
//   nChainIndex int                     // }
//   parentBlock CBlock (header only)

func (h *AuxPowHeader) BtcEncode(w io.Writer, pver uint32) error {
	err := h.CoinbaseTx.BtcEncode(w, pver)
	if err != nil {
		return err
	}

	err = writeElement(w, &h.ParentBlockHash)
	if err != nil {
		return err
	}

	err = h.CoinbaseBranch.BtcEncode(w, pver)
	if err != nil {
		return err
	}

	err = h.BlockChainBranch.BtcEncode(w, pver)
	if err != nil {
		return err
	}

	err = h.ParentBlockHeader.Serialize(w)
	if err != nil {
		return err
	}

	return nil
}

func (h *AuxPowHeader) BtcDecode(r io.Reader, pver uint32) error {
	err := h.CoinbaseTx.BtcDecode(r, pver)
	if err != nil {
		return err
	}

	err = readElement(r, &h.ParentBlockHash)
	if err != nil {
		return err
	}

	err = h.CoinbaseBranch.BtcDecode(r, pver)
	if err != nil {
		return err
	}

	err = h.BlockChainBranch.BtcDecode(r, pver)
	if err != nil {
		return err
	}

	err = h.ParentBlockHeader.Deserialize(r)
	if err != nil {
		return err
	}

	return nil
}

func (h *AuxPowHeader) Deserialize(r io.Reader) error {
	return h.BtcDecode(r, 0)
}

func (h *AuxPowHeader) SerializeSize() int {
	n := HashSize + blockHeaderHeaderLen
	n += h.CoinbaseTx.SerializeSize()
	n += h.CoinbaseBranch.SerializeSize()
	n += h.BlockChainBranch.SerializeSize()
	return n
}
