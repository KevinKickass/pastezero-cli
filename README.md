# PasteZero CLI

[![Release Build](https://github.com/KevinKickass/pastezero-cli/actions/workflows/release.yml/badge.svg)](https://github.com/KevinKickass/pastezero-cli/actions)
[![Latest Release](https://img.shields.io/github/v/release/KevinKickass/pastezero-cli?sort=semver)](https://github.com/KevinKickass/pastezero-cli/releases)


Kommandozeilentool für sicheren Datei-Upload und -Download mit [PasteZero.de](https://pastezero.de).
Ende-zu-Ende verschlüsselt, zero-trust, plattformübergreifend.

---

## Funktionen

* AES-256-GCM-Verschlüsselung (E2EE)
* Datei-Upload mit zufälligem Einmal-Schlüssel
* Datei-Download per UUID und Key
* Kompatibel mit der https://pastezero.de Weboberfläche
* Plattformübergreifende Binaries (Linux, macOS, Windows)
* Automatischer GitHub Release bei Tag

---

## Installation

### Vorgefertigte Binaries

Siehe [Releases](https://github.com/KevinKickass/pastezero-cli/releases):

```bash
curl -LO https://github.com/KevinKickass/pastezero-cli/releases/download/v0.1.3/pastezero-v0.1.3-linux-amd64.tar.gz
tar -xzf pastezero-v0.1.3-linux-amd64.tar.gz
chmod +x pastezero
./pastezero --help
```

### Selbst bauen

```bash
git clone git@github.com:KevinKickass/pastezero-cli.git
cd pastezero-cli
go build -o pastezero main.go
```

---

## Nutzung

### Upload

```bash
pastezero upload --file geheim.pdf
```

→ Gibt dir einen Link wie:

```
https://pastezero.de/get/<uuid>#<base64-key>
```

### Download

```bash
pastezero download https://pastezero.de/get/abc123#BASE64KEY
```

Oder:

```bash
pastezero download --id abc123 --key BASE64KEY
```

Optionaler Zielpfad:

```bash
pastezero download abc123#KEY --output /tmp/zieldatei.txt
```

---

## Build lokal

```bash
./scripts/build.sh       # erstellt ./build/ + ./release/
```

## GitHub Release (per Tag)

```bash
git tag v0.2.0
git push origin v0.2.0
```

Der Workflow `.github/workflows/release.yml` wird automatisch ausgeführt.

---

## Sicherheit

* AES-256-GCM mit 12-Byte IV und 32-Byte Schlüssel
* Original-Dateiname und MIME-Type verschlüsselt eingebettet
* Keine Klartext-Metadaten
* Keine Passwortspeicherung

---

## Lizenz

GNU General Public License v3.0 – © 2025 [Stage One Solutions](https://stageone.solutions)  
Maintained by [KevinKickass](https://github.com/KevinKickass)
