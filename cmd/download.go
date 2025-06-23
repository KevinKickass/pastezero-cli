package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [link|uuid[#key]]",
	Short: "Download und Entschlüsselung einer Datei",
	Long: `Lädt eine verschlüsselte Datei von PasteZero herunter und entschlüsselt sie lokal.
Akzeptiert vollständige Links mit #key oder fragt nach dem Key bei Bedarf.

Beispiele:
  pastezero download https://pastezero.de/get/abc123#BASE64KEY
  pastezero download abc123 --key BASE64KEY
  pastezero download --link https://... --output geheim.txt
`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api")
		outPath, _ := cmd.Flags().GetString("output")
		linkFlag, _ := cmd.Flags().GetString("link")
		idFlag, _ := cmd.Flags().GetString("id")
		keyFlag, _ := cmd.Flags().GetString("key")

		// UUID + Key extrahieren
		var id, key string
		if len(args) > 0 {
			id, key = extractUUIDAndKey(args[0])
		} else if linkFlag != "" {
			id, key = extractUUIDAndKey(linkFlag)
		} else if idFlag != "" {
			id = strings.TrimSpace(idFlag)
		}

		if id == "" {
			log.Fatal("Fehlt: UUID oder gültiger Link (--id, --link oder Argument)")
		}
		if key == "" && keyFlag != "" {
			key = keyFlag
		}
		if key == "" {
			fmt.Print("Schlüssel?: ")
			fmt.Scanln(&key)
		}

		// base64 → []byte
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil || len(keyBytes) != 32 {
			log.Fatal("Ungültiger Key: Muss base64-kodiert und 32 Byte groß sein")
		}

		// Abrufen
		url := fmt.Sprintf("%s/file/%s", apiURL, id)
		res, err := http.Get(url)
		if err != nil || res.StatusCode != 200 {
			body, _ := io.ReadAll(res.Body)
			log.Fatalf("Download fehlgeschlagen (%d): %s", res.StatusCode, string(body), url)
		}
		defer res.Body.Close()
		encrypted, err := io.ReadAll(res.Body)
		if err != nil || len(encrypted) < 13 {
			log.Fatal("Ungültige verschlüsselte Datei")
		}

		iv := encrypted[:12]
		ciphertext := encrypted[12:]

		// Entschlüsselung
		block, _ := aes.NewCipher(keyBytes)
		aesgcm, _ := cipher.NewGCM(block)
		plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
		if err != nil {
			log.Fatalf("Entschlüsselung fehlgeschlagen: %v", err)
		}

		// Zerlegen
		buf := bytes.NewBuffer(plaintext)
		nameLen, _ := buf.ReadByte()
		filename := make([]byte, nameLen)
		buf.Read(filename)

		mimeLen, _ := buf.ReadByte()
		mimeType := make([]byte, mimeLen)
		buf.Read(mimeType)

		content := buf.Bytes()
		finalPath := outPath
		if finalPath == "" {
			finalPath = string(filename)
		}

		err = os.WriteFile(finalPath, content, 0600)
		if err != nil {
			log.Fatalf("Speichern fehlgeschlagen: %v", err)
		}

		fmt.Printf("Datei erfolgreich entschlüsselt und gespeichert unter: %s\n", finalPath)
	},
}

func init() {
	downloadCmd.Flags().String("id", "", "UUID der Datei")
	downloadCmd.Flags().String("key", "", "AES-Schlüssel (base64)")
	downloadCmd.Flags().String("link", "", "Kompletter Link (https://pastezero.de/get/<id>#<key>)")
	downloadCmd.Flags().String("output", "", "Zielpfad zum Speichern der entschlüsselten Datei")
	rootCmd.AddCommand(downloadCmd)
}

// Extrahiert UUID und optional Key aus Link oder Argument
func extractUUIDAndKey(input string) (string, string) {
	uuid := ""
	key := ""

	if strings.Contains(input, "/get/") {
		parts := strings.Split(input, "/get/")
		if len(parts) == 2 {
			main := parts[1]
			if strings.Contains(main, "#") {
				split := strings.SplitN(main, "#", 2)
				uuid = split[0]
				key = split[1]
			} else {
				uuid = main
			}
		}
	} else if strings.Contains(input, "#") {
		split := strings.SplitN(input, "#", 2)
		uuid = split[0]
		key = split[1]
	} else {
		uuid = input
	}

	return strings.TrimSpace(uuid), strings.TrimSpace(key)
}