package crypto

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"
)

type ZlibCompressor[T any] struct {
	IEncryptor[T]
}

func (z *ZlibCompressor[T]) Encrypt(l T) (string, bool) {
	r, err := json.Marshal(l)
	if err != nil {
		return "", false
	}
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(r))
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes()), true
}

func (z *ZlibCompressor[T]) Decrypt(l string) (T, bool) {
	var r T
	d, err := base64.StdEncoding.DecodeString(l)
	if err != nil {
		return r, false
	}
	b := new(bytes.Buffer)
	zr, err := zlib.NewReader(bytes.NewBuffer(d))
	if err != nil {
		return r, false
	}
	io.Copy(b, zr)
	zr.Close()
	err = json.Unmarshal(b.Bytes(), &r)
	if err != nil {
		return r, false
	}
	return r, true
}
