// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcwire

import (
	"bytes"
	"io"
	"time"
)

// BlockVersion is the current latest supported block version.
const BlockVersion = 2

// Version 4 bytes + Timestamp 4 bytes + Bits 4 bytes + Nonce 4 bytes +
// PrevBlock and MerkleRoot hashes.
const MaxBlockHeaderHeaderPayload = 16 + (HashSize * 2)

// BlockHeaderHeader defines information about a block and is used in the bitcoin
// block (MsgBlock) and headers (MsgHeaders) messages.
type BlockHeaderHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version int32

	// Hash of the previous block in the block chain.
	PrevBlock ShaHash

	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot ShaHash

	// Time the block was created.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.
	Timestamp time.Time

	// Difficulty target for the block.
	Bits uint32

	// Nonce used to generate the block.
	Nonce uint32
}

// blockHeaderHeaderLen is a constant that represents the number of bytes for a block
// header excluding any AuxPow part.
const blockHeaderHeaderLen = 80

// AuxPow
const auxPowFlag int32 = (1 << 8)
const chainIDStart int32 = (1 << 16)
const chainIDEnd int32 = (1 << 30)

func (h *BlockHeaderHeader) AuxPow() bool {
	return (h.Version & auxPowFlag) != 0
}

func (h *BlockHeaderHeader) SetAuxPow(auxpow bool) {
	if auxpow {
		h.Version |= auxPowFlag
	} else {
		h.Version &= ^auxPowFlag
	}
}

func (h *BlockHeaderHeader) ChainID() uint32 {
	return uint32(h.Version / chainIDStart)
}

func (h *BlockHeaderHeader) BlockVersion() int32 {
	return h.Version & 0xFF
}

// BlockSha computes the block identifier hash for the given block header.
func (h *BlockHeaderHeader) BlockSha() (ShaHash, error) {
	// Encode the header and run double sha256 everything prior to the
	// number of transactions.  Ignore the error returns since there is no
	// way the encode could fail except being out of memory which would
	// cause a run-time panic.  Also, SetBytes can't fail here due to the
	// fact DoubleSha256 always returns a []byte of the right size
	// regardless of input.
	var buf bytes.Buffer
	var sha ShaHash
	_ = writeBlockHeaderHeader(&buf, 0, h)
	_ = sha.SetBytes(DoubleSha256(buf.Bytes()[0:blockHeaderHeaderLen]))

	// Even though this function can't currently fail, it still returns
	// a potential error to help future proof the API should a failure
	// become possible.
	return sha, nil
}

// Deserialize decodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (h *BlockHeaderHeader) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of readBlockHeader.
	return readBlockHeaderHeader(r, 0, h)
}

// Serialize encodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (h *BlockHeaderHeader) Serialize(w io.Writer) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of writeBlockHeader.
	return writeBlockHeaderHeader(w, 0, h)
}

// readBlockHeader reads a bitcoin block header from r.  See Deserialize for
// decoding block headers stored to disk, such as in a database, as opposed to
// decoding from the wire.
func readBlockHeaderHeader(r io.Reader, pver uint32, bh *BlockHeaderHeader) error {
	var sec uint32
	err := readElements(r, &bh.Version, &bh.PrevBlock, &bh.MerkleRoot, &sec,
		&bh.Bits, &bh.Nonce)
	if err != nil {
		return err
	}
	bh.Timestamp = time.Unix(int64(sec), 0)

	return nil
}

// writeBlockHeader writes a bitcoin block header to w.  See Serialize for
// encoding block headers to be stored to disk, such as in a database, as
// opposed to encoding for the wire.
func writeBlockHeaderHeader(w io.Writer, pver uint32, bh *BlockHeaderHeader) error {
	sec := uint32(bh.Timestamp.Unix())
	err := writeElements(w, bh.Version, &bh.PrevBlock, &bh.MerkleRoot,
		sec, bh.Bits, bh.Nonce)
	if err != nil {
		return err
	}

	return nil
}

// BlockHeader defines information about a block and is used in the bitcoin
// block (MsgBlock) and headers (MsgHeaders) messages.
type BlockHeader struct {
	BlockHeaderHeader

	// AuxPow
	AuxPowHeader AuxPowHeader
}

const MaxBlockHeaderPayload = MaxBlockHeaderHeaderPayload + MaxAuxPowSize

func (h *BlockHeader) SerializeSize() int {
	L := blockHeaderHeaderLen
	if h.AuxPow() {
		L += h.AuxPowHeader.SerializeSize()
	}
	return L
}

// Deserialize decodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (h *BlockHeader) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of readBlockHeader.
	return readBlockHeader(r, 0, h)
}

// Serialize encodes a block header from r into the receiver using a format
// that is suitable for long-term storage such as a database while respecting
// the Version field.
func (h *BlockHeader) Serialize(w io.Writer) error {
	// At the current time, there is no difference between the wire encoding
	// at protocol version 0 and the stable long-term storage format.  As
	// a result, make use of writeBlockHeader.
	return writeBlockHeader(w, 0, h)
}

// readBlockHeader reads a bitcoin block header from r.  See Deserialize for
// decoding block headers stored to disk, such as in a database, as opposed to
// decoding from the wire.
func readBlockHeader(r io.Reader, pver uint32, bh *BlockHeader) error {
	err := readBlockHeaderHeader(r, pver, &bh.BlockHeaderHeader)
	if err != nil {
		return err
	}

	if bh.AuxPow() {
		err = bh.AuxPowHeader.BtcDecode(r, pver)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeBlockHeader writes a bitcoin block header to w.  See Serialize for
// encoding block headers to be stored to disk, such as in a database, as
// opposed to encoding for the wire.
func writeBlockHeader(w io.Writer, pver uint32, bh *BlockHeader) error {
	err := writeBlockHeaderHeader(w, pver, &bh.BlockHeaderHeader)
	if err != nil {
		return err
	}

	if bh.AuxPow() {
		err = bh.AuxPowHeader.BtcEncode(w, pver)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewBlockHeader returns a new BlockHeader using the provided previous block
// hash, merkle root hash, difficulty bits, and nonce used to generate the
// block with defaults for the remaining fields.
func NewBlockHeader(prevHash *ShaHash, merkleRootHash *ShaHash, bits uint32,
	nonce uint32) *BlockHeader {

	// Limit the timestamp to one second precision since the protocol
	// doesn't support better.
	return &BlockHeader{
		BlockHeaderHeader: BlockHeaderHeader{
			Version:    BlockVersion,
			PrevBlock:  *prevHash,
			MerkleRoot: *merkleRootHash,
			Timestamp:  time.Unix(time.Now().Unix(), 0),
			Bits:       bits,
			Nonce:      nonce,
		},
	}
}
