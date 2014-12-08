package btcwire_test

import "bytes"
import "github.com/hlandauf/btcwire"
import "testing"
import "encoding/hex"

func TestMerkleBranch(t *testing.T) {
	mb1 := btcwire.MerkleBranch{
		Hashes: []btcwire.ShaHash{
			*newShaHashFromStr("b98db090398ebc4342951f9ba89b3e0110bdc757714b80c695663c9060113639"),
			*newShaHashFromStr("3e0a60195218f27df0edc1d5b008568b2754f8a709eb80e3c1412bdfcb3b7e21"),
		},
		SideMask: 0,
	}

	if mb1.Size() != 2 {
		t.FailNow()
	}

	w := bytes.Buffer{}
	err := mb1.Serialize(&w)
	if err != nil {
		t.FailNow()
	}

	bb := w.Bytes()
	bbRef, err := hex.DecodeString("0239361160903c6695c6804b7157c7bd10013e9ba89b1f954243bc8e3990b08db9217e3bcbdf2b41c1e380eb09a7f854278b5608b0d5c1edf07df2185219600a3e00000000")
	if err != nil {
		t.FailNow()
	}

	if bytes.Compare(bb, bbRef) != 0 {
		t.FailNow()
	}

	r := bytes.NewReader(bb)
	mb2 := btcwire.MerkleBranch{}
	err = mb2.Deserialize(r)
	if err != nil {
		t.FailNow()
	}

	w = bytes.Buffer{}
	err = mb2.Serialize(&w)
	if err != nil {
		t.FailNow()
	}

	bb = w.Bytes()

	if bytes.Compare(bb, bbRef) != 0 {
		t.FailNow()
	}

	h := newShaHashFromStr("d8f244c159278ea8cfffcbe1c463edef33d92d11d36ac3c62efd3eb7ff3a5dbf")
	rh, err := mb2.DetermineRoot(h)
	if err != nil {
		t.FailNow()
	}

	t.Logf("rh: %s", rh.String())

	if rh.String() != "bf0ca48d50405f62cb40fa67c6f9fd9309e9a5fcb2ad05d3976ecb28839b4474" {
		t.FailNow()
	}

	mb3 := btcwire.MerkleBranch{
		Hashes: []btcwire.ShaHash{
			*newShaHashFromStr("d8f244c159278ea8cfffcbe1c463edef33d92d11d36ac3c62efd3eb7ff3a5dbf"),
			*newShaHashFromStr("3e0a60195218f27df0edc1d5b008568b2754f8a709eb80e3c1412bdfcb3b7e21"),
		},
		SideMask: 1,
	}

	h = newShaHashFromStr("b98db090398ebc4342951f9ba89b3e0110bdc757714b80c695663c9060113639")
	rh, err = mb3.DetermineRoot(h)
	if err != nil {
		t.FailNow()
	}

	t.Logf("rh: %s", rh.String())
	if rh.String() != "bf0ca48d50405f62cb40fa67c6f9fd9309e9a5fcb2ad05d3976ecb28839b4474" {
		t.FailNow()
	}

	mb4 := btcwire.MerkleBranch{
		Hashes: []btcwire.ShaHash{
			*newShaHashFromStr("d377b92dd7af8f1b25b2ac96f5ac68d0d8ae0e15fc370f89ea0fa36c3d753266"),
			*newShaHashFromStr("f01b8b33d4737f715303d502cd8dda6b2ea4f9513c169d94b18b5f2fa1a367b7"),
		},
		SideMask: 2,
	}
	h = newShaHashFromStr("d377b92dd7af8f1b25b2ac96f5ac68d0d8ae0e15fc370f89ea0fa36c3d753266")
	rh, err = mb4.DetermineRoot(h)
	if err != nil {
		t.FailNow()
	}

	t.Logf("rh: %s", rh.String())
	if rh.String() != "bf0ca48d50405f62cb40fa67c6f9fd9309e9a5fcb2ad05d3976ecb28839b4474" {
		t.FailNow()
	}
}

func newShaHashFromStr(hexStr string) *btcwire.ShaHash {
	sha, err := btcwire.NewShaHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return sha
}
