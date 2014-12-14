package btcwire

// Amount to lock in name transactions. This is not (yet) enforced by the
// protocol, but for acceptance to the mempool.
const NameLockedAmount = 100000000/100

const MinFirstUpdateDepth = 12
const MaxNameLength = 255
const MaxNameValueLength = 1023
const MaxNameValueLengthUI = 520
const MempoolHeight = 0x7FFFFFFF

func NamecoinLenientVersionCheck(txHeight int64) bool {
  // TODO: update after soft fork.
  return true
}

func NameExpirationDepth(height int64) int64 {
  // Important: It is assumed in ExpireNames that "n - expirationDepth(n)" is
  // increasing. (This is the update height up to which names expire at height n.)
  if height < 24000 {
    return 12000
  }
  if height < 48000 {
    return height - 12000
  }
  return 36000
}

func IsNameExpired(prevHeight int64, height int64) bool {
  if prevHeight == MempoolHeight {
    return false
  }

  if height == MempoolHeight {
    panic("not supported")
    //height = CHAINACTIVEHEIGHT
  }

  return (prevHeight + NameExpirationDepth(height)) <= height
}
