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
  compress    Compress or decompress input using gzip
  convert     Convert units: length, weight, temperature, storage, time
  countdown   Countdown timer (e.g. 1h30m, 5m, 90s)
  download    Download a file from a URL
  entropy     Calculate Shannon entropy of input
  random      Generate a cryptographically secure random string
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
echo -n 'lorem ipsum' | orbx hash md5
```

## Clear Clipboard

> **Note:** `clearclip` is supported on macOS and Linux (X11) only. On Linux, `xclip` must be installed (`apt install xclip`). Windows is not currently supported.

```bash
orbx clearclip
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

## Entropy

Calculates the Shannon entropy of an input, expressed in bits per byte.

> **Note:** By default, `entropy` attempts to decode the input as base64 before calculating — consistent with the output format of the `aes encrypt` command. If your input is not base64-encoded, pass the `--raw` flag to skip decoding and calculate entropy on the raw input directly.

```bash
# Direct input (base64)
orbx entropy 'SGVsbG8='

# Skip base64 decoding, treat input as raw bytes
orbx entropy 'Hello!' --raw

# From file
orbx entropy --file ciphertext.txt

# Good entropy
head -c 10000 /dev/urandom | orbx base64 encode | orbx entropy
```

### Interpreting the result

| bits/byte | Meaning                                            |
| --------- | -------------------------------------------------- |
| ~8.0      | Excellent — output looks random (good ciphertext)  |
| 6.0–7.9   | High entropy — likely compressed or encrypted data |
| 3.0–5.9   | Moderate — structured data, natural language       |
| < 3.0     | Low — highly repetitive or predictable input       |

A well-formed AES-GCM ciphertext should score close to **8.0 bits/byte**.

## Compression

By default, output is base64-encoded for compatibility with other `orbx` commands (e.g. piping into `aes encrypt`). Use `--raw` to write raw gzip bytes, which can be opened with 7-Zip, `gunzip`, or any standard gzip tool.

> **Note:** `--raw` output is fully compatible with `gzip`/`gunzip` for single files. Archives created with `tar -czf` (`.tar.gz`) are not supported — use `tar` directly for directories.

```bash
# Compress a string and save to file (base64 output)
orbx compress 'Hello from the other side!' --out compressed.txt

# Compress a file and save output (base64 output)
orbx compress --file largefile.txt --out compressed.txt

# Decompress a base64-compressed file
orbx compress --file compressed.txt --decode

# Compress to a raw .gz file (compatible with 7-Zip, gunzip, etc.)
orbx compress --file largefile.txt --raw --out archive.gz

# Decompress a raw .gz file
orbx compress --file archive.gz --raw --decode

# Compress and decompress via pipeline
cat largefile.txt | orbx compress | orbx compress --decode

# Compress then encrypt
cat largefile.txt | orbx compress | orbx aes encrypt --key secret.key --out out.enc

# Decrypt then decompress
orbx aes decrypt --file out.enc --key secret.key | orbx compress --decode
```

## 🔧 Requirements

- Go 1.20+
