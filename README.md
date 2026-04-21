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
  clearclip   Clear system clipboard
  convert     Convert units: length, weight, temperature, storage, time
  countdown   Countdown timer (e.g. 1h30m, 5m, 90s)
  download    Download a file from a URL

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
  jwt         Decode a JWT token (header and payload, no verification)
  ports       Show processes using network ports
  serve       Start a static file server in current directory
  text        String utilities
  unixts      Unix timestamp utilities
  uuid        Generate a UUID v4

📦 Content
  xkcd        Fetch latest XKCD comic

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  url         Decode and parse a URL

Flags:
  -h, --help      help for orbx
  -v, --version   show version

Use "orbx [command] --help" for more information about a command.
```

## 🔧 Requirements

- Go 1.20+
