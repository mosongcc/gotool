package gobjid

// Hex 生成16进制字符串ObjectID，长度24位
func Hex() string {
	return NewObjectID().Hex()
}
