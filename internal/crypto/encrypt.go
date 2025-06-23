package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func EncryptPayload(filePath string) ([]byte, []byte, []byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	stat, _ := file.Stat()
	filename := []byte(stat.Name())
	mimetype := []byte("application/octet-stream") // optional: mime guessing

	var buf bytes.Buffer
	buf.WriteByte(byte(len(filename)))
	buf.Write(filename)
	buf.WriteByte(byte(len(mimetype)))
	buf.Write(mimetype)
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, nil, nil, err
	}

	plaintext := buf.Bytes()
	key := make([]byte, 32)
	iv := make([]byte, 12)
	rand.Read(key)
	rand.Read(iv)

	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)
	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)

	return ciphertext, key, iv, nil
}
