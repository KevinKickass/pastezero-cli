package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"bytes"
	"encoding/json"
	"pastezero-cli/internal/crypto"
	"pastezero-cli/internal/api"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Datei verschlüsseln und bei PasteZero hochladen",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, _ := cmd.Flags().GetString("file")
		apiURL, _ := cmd.Flags().GetString("api")

		if inputFile == "" {
			log.Fatal("Erforderlich: --file")
		}

		cfg, err := api.EnsureClientID(apiURL)
		if err != nil {
			log.Fatalf("Client-ID konnte nicht ermittelt werden: %v", err)
		}

		fmt.Println("Datei verschlüsseln...")
		payload, key, iv, err := crypto.EncryptPayload(inputFile)
		if err != nil {
			log.Fatalf("Fehler bei Verschlüsselung: %v", err)
		}
		fullPayload := append(iv, payload...)

		fmt.Println("Datei hochladen...")
		req, _ := http.NewRequest("POST", apiURL+"/upload", io.NopCloser(
			// kein Content-Type nötig
			bytes.NewReader(fullPayload),
		))
		req.Header.Set("X-Client-ID", cfg.ClientID)
		req.Header.Set("X-Client-Signature", cfg.Signature)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Fehler beim Upload: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(res.Body)
			log.Fatalf("Upload fehlgeschlagen: %s", string(b))
		}

		type uploadResp struct {
			Link string `json:"link"`
		}
		var r uploadResp
		json.NewDecoder(res.Body).Decode(&r)

		fmt.Println("Upload erfolgreich abgeschlossen")
		fmt.Println("Download-Link (All-In-One):", "https://pastezero.de" + r.Link + "#" + base64.StdEncoding.EncodeToString(key))
		fmt.Println("-------------------------------------------")
		fmt.Println("Oder getrennt (sicherer):")
		fmt.Println("Download-Link:", "https://pastezero.de" + r.Link)
		fmt.Println("Schlüssel:", base64.StdEncoding.EncodeToString(key))
	},
}

func init() {
	uploadCmd.Flags().StringP("file", "f", "", "Pfad zur Datei")
	rootCmd.AddCommand(uploadCmd)
}
