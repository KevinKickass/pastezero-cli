# PasteZero CLI

Command-Line-Tool für sicheren Datei-Upload und -Download mit [PasteZero.de](https://pastezero.de)  
Ende-zu-Ende-Verschlüsselung (E2EE) im Zero-Trust-Modell – vollständig clientseitig.

---

## Features

- Upload beliebiger Dateien mit starker AES-256-GCM-Verschlüsselung
- Download verschlüsselter Dateien über UUID + Key
- Automatische Client-ID-Registrierung beim ersten Start
- Volle Kompatibilität zur Web-Oberfläche (https://pastezero.de)
- Plattformübergreifende Binaries (Linux, macOS, Windows)

---

## Installation

### Vorgefertigte Binaries

Lade dein Binary von GitHub Releases herunter:

```bash
curl -LO https://github.com/KevinKickass/pastezero-cli/releases/download/v0.1.0/pastezero-v0.1.0-linux-amd64.tar.gz
tar -xzf pastezero-v0.1.0-linux-amd64.tar.gz
chmod +x pastezero
./pastezero --help
Selbst kompilieren
bash
Kopieren
Bearbeiten
git clone https://github.com/KevinKickass/pastezero-cli
cd pastezero-cli
go build -o pastezero main.go
Nutzung
Upload
bash
Kopieren
Bearbeiten
pastezero upload --file geheim.pdf
→ Gibt einen Download-Link wie
https://pastezero.de/get/<uuid>#<base64-key>
Dieser Link ist alles, was der Empfänger benötigt.

Download
bash
Kopieren
Bearbeiten
pastezero download https://pastezero.de/get/abc123#BASE64KEY
Alternativ:

bash
Kopieren
Bearbeiten
pastezero download abc123 --key BASE64KEY
Oder interaktiv:

bash
Kopieren
Bearbeiten
pastezero download https://pastezero.de/get/abc123
# → CLI fragt nach dem Key
Mit Zielpfad:

bash
Kopieren
Bearbeiten
pastezero download abc123#BASE64KEY --output /tmp/zieldatei.txt
Build
bash
Kopieren
Bearbeiten
./scripts/build.sh       # erstellt Binaries + .tar.gz/.zip in ./release/
./scripts/release.sh     # (optional) GitHub Release hochladen (setzt Tag voraus)
Sicherheit
AES-256-GCM, 12-Byte IV, 32-Byte Schlüssel

Symmetrische One-Time-Keys pro Upload

Dateiname und MIME-Type verschlüsselt eingebettet

Keine Passwörter, kein Backend-Zugriff auf Inhalte

Kompatibel mit Browser-Client

Lizenz
MIT License – © 2025 Stage One Solutions
Repo maintained by KevinKickass