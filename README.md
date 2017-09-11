# Raiden - Simple GitHub Repository Cleaner
**Documentation:** [![GoDoc](https://godoc.org/github.com/suusan2go/raiden?status.svg)](https://godoc.org/github.com/suusan2go/raiden)

Raiden is a simple GitHub Repository Cleaner
- [x] clean old releases by name or created time
- [ ] clean old branches
- [ ] clean old pull-requests

# Installation
```bash
$ go get github.com/suusan2go/raiden
```

# Usage
## Clean Release
```bash
# clean suusan2go/hoge repo releases created 1 months ago
$ raiden releases clean -r raiden -o suusan2go --months 1
```
