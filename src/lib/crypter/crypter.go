package crypter

// Crypter 加密器
type Crypter interface {
	Encode(src []byte) []byte
	Decode(src []byte) []byte
}
