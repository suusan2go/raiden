# Raiden - Simple GitHub Repository Cleaner
Raiden is a simple GitHub Repository Cleaner
- [x] clean old releases
- [ ] clean old branches
- [ ] clean old pull-requests

# Installation
```bash
$ go get github.com/suzan2go/raiden
```

# Usage
## Clean Release
```bash
# clean suzan2go/hoge repo releases created 1 months ago
$ raiden releases clean -r raiden -o suzan2go --months 1
```
