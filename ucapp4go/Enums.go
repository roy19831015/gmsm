package ucapp4go

type SymmType int

const (
	AES256 SymmType = iota - 1
	SM4
	DES
	DESede
	SM1
)

//type HashType int
//const (
//	SHA256 HashType = iota - 1
//	SHA1
//	SM3
//	MD5
//)

type AsymmAlgType int

const (
	RSA AsymmAlgType = iota - 1
	SM2
)
