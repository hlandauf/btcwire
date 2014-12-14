package btcwire

import "io"
import "fmt"

type NameInfo struct {
	Key string

	Value        string
	Height       int64
	OutPoint     *OutPoint
	Addr         []byte
}

func (ni *NameInfo) BtcEncode(w io.Writer, pver uint32) error {
  err := writeVarString(w, pver, ni.Key)
  if err != nil {
    return fmt.Errorf("failed to serialize name info key: %v", err)
  }

	err = writeVarString(w, pver, ni.Value)
	if err != nil {
    return fmt.Errorf("failed to serialize name info value: %v", err)
	}

	err = writeElement(w, ni.Height)
	if err != nil {
    return fmt.Errorf("failed to serialize name info height: %v", err)
	}

	err = writeOutPoint(w, pver, 1, ni.OutPoint)
	if err != nil {
    return fmt.Errorf("failed to serialize name info outpoint: %v", err)
	}

	err = writeVarBytes(w, pver, ni.Addr)
	if err != nil {
    return fmt.Errorf("failed to serialize name info addr: %v", err)
	}

	return nil
}

func (ni *NameInfo) Serialize(w io.Writer) error {
	return ni.BtcEncode(w, 0)
}

func (ni *NameInfo) BtcDecode(r io.Reader, pver uint32) error {
  k, err := readVarString(r, pver)
  if err != nil {
    return err
  }

  ni.Key = k

  v, err := readVarString(r, pver)
	if err != nil {
		return err
	}

  ni.Value = v

	err = readElement(r, &ni.Height)
	if err != nil {
		return err
	}

  ni.OutPoint = &OutPoint{}
	err = readOutPoint(r, pver, 1, ni.OutPoint)
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

func (ni *NameInfo) IsExpired(height int64) bool {
  return IsNameExpired(ni.Height, height)
}
