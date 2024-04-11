package crypto

type IEncryptor[T any] interface {
	Encrypt(T) (string, bool)
	Decrypt(string) (T, bool)
}
