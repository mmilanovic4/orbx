# orbx

orbx is a lightweight CLI utility for quick system tasks.

Some of the available commands include clipboard utilities, a lightweight static web server for serving local directories and system/network information tools.

---

## ⚡ Install

```bash
go install github.com/mmilanovic4/orbx@latest
```

## 📦 Usage

```bash
$ orbx --help

orbx is a lightweight CLI utility for quick system tasks.

Usage:
  orbx [flags]
  orbx [command]

🧰 Utilities
  aes         AES-GCM encryption utilities
  clearclip   Clear system clipboard
  convert     Convert units: length, weight, temperature, storage, time
  countdown   Countdown timer (e.g. 1h30m, 5m, 90s)
  download    Download a file from a URL
  watch       Repeatedly run a command every N seconds

🌐 Network Tools
  dns         Resolve DNS records for a domain
  headers     Show HTTP response status and headers
  ip          Show public and local IP addresses
  ping        HTTP latency check (like ping, but for URLs)
  rdns        Reverse DNS lookup for an IP address
  tcpcheck    Check TCP port connectivity

💻 Developer Tools
  base64      Encode or decode base64
  hash        Generate hash of a string
  hex         Encode or decode hex
  html        Encode or decode HTML entities
  jwt         Decode a JWT token (header and payload, no verification)
  ports       Show processes using network ports
  serve       Start a static file server in current directory
  text        String utilities
  unixts      Unix timestamp utilities
  url         Decode and parse a URL
  uuid        Generate a UUID v4

📦 Content
  xkcd        Fetch latest XKCD comic

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help      help for orbx
  -v, --version   show version

Use "orbx [command] --help" for more information about a command.
```

Stdin support is available across most tools, allowing you to pipe input directly instead of passing it as an argument:

```
echo -n "lorem ipsum" | orbx hash md5
```

## AES-GCM Encryption

```bash
# AES-256 (default)
orbx aes key --out secret.key

# AES-128 or AES-192
orbx aes key 16 --out secret.key
orbx aes key 24 --out secret.key

# Encrypt and save to file
orbx aes encrypt 'Hello from the other side!' --key secret.key --out ciphertext.txt

# Decrypt from file
orbx aes decrypt --key secret.key --file ciphertext.txt
```

## 🔧 Requirements

- Go 1.20+
