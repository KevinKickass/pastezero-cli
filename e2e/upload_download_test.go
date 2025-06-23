package e2e

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pastezero-cli/internal/api"
	"pastezero-cli/internal/crypto"
)

func TestUploadDownloadE2E(t *testing.T) {
	tmpfile := filepath.Join(os.TempDir(), "pastezero_test.txt")
	content := []byte("Testinhalt für E2E")
	err := os.WriteFile(tmpfile, content, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile)

	ciphertext, key, iv, err := crypto.EncryptPayload(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	fullPayload := append(iv, ciphertext...)

	cfg, err := api.EnsureClientID("https://api.pastezero.de")
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("POST", "https://api.pastezero.de/upload", bytes.NewReader(fullPayload))
	req.Header.Set("X-Client-ID", cfg.ClientID)
	req.Header.Set("X-Client-Signature", cfg.Signature)
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("Upload fehlgeschlagen: %v", err)
	}
	defer res.Body.Close()

	var data struct {
		Link string `json:"link"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(data.Link, "/get/")
	if len(parts) != 2 {
		t.Fatalf("Ungültiger Link: %s", data.Link)
	}
	uuid := strings.Split(parts[1], "#")[0]

	getURL := "https://api.pastezero.de/file/" + uuid
	getRes, err := http.Get(getURL)
	if err != nil || getRes.StatusCode != 200 {
		t.Fatalf("Download fehlgeschlagen")
	}
	defer getRes.Body.Close()
	encrypted, _ := io.ReadAll(getRes.Body)

	ivBack := encrypted[:12]
	ciphertextBack := encrypted[12:]

	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)
	plaintext, err := aesgcm.Open(nil, ivBack, ciphertextBack, nil)
	if err != nil {
		t.Fatalf("Entschlüsselung fehlgeschlagen: %v", err)
	}

	buf := bytes.NewBuffer(plaintext)
	nameLen, _ := buf.ReadByte()
	filename := make([]byte, nameLen)
	buf.Read(filename)

	mimeLen, _ := buf.ReadByte()
	mime := make([]byte, mimeLen)
	buf.Read(mime)

	result := buf.Bytes()
	if !bytes.Equal(result, content) {
		t.Fatalf("Inhalt stimmt nicht überein")
	}
}
