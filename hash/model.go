package hash

type hash interface {
	getHash(string) string
}