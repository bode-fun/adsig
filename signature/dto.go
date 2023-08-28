package signature

type signaturesForAccount struct {
	Signatures []signatureForAccount
}

type signatureForAccount struct {
	Name string
	Path string
}
