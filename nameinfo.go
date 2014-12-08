package btcwire

import "io"

type NameInfo struct {
	Key string

	Value        string
	Height       int
	PrevOutPoint OutPoint
	Addr         []byte
}

func (ni *NameInfo) BtcEncode(w io.Writer, pver uint32) error {
	err := writeElement(w, ni.Value)
	if err != nil {
		return err
	}

	err = writeElement(w, ni.Height)
	if err != nil {
		return err
	}

	err = writeOutPoint(w, pver, 1, &ni.PrevOutPoint)
	if err != nil {
		return err
	}

	err = writeVarBytes(w, pver, ni.Addr)
	if err != nil {
		return err
	}

	return nil
}

func (ni *NameInfo) Serialize(w io.Writer) error {
	return ni.BtcEncode(w, 0)
}

func (ni *NameInfo) BtcDecode(r io.Reader, pver uint32) error {
	err := readElement(r, &ni.Value)
	if err != nil {
		return err
	}

	err = readElement(r, &ni.Height)
	if err != nil {
		return err
	}

	err = readOutPoint(r, pver, 1, &ni.PrevOutPoint)
	if err != nil {
		return err
	}

	ni.Addr, err = readVarBytes(r, pver, 4096, "addr")
	if err != nil {
		return err
	}

	return nil
}

func (ni *NameInfo) Deserialize(r io.Reader) error {
	return ni.BtcDecode(r, 0)
}
