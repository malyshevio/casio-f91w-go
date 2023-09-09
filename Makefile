.DEFAULT_GOAL := apps

PROJECT_URL := github.com/MawKKe/casio-f91w-go

# ------
#  Apps
# ------

out/hourly:
	go build -o $@ ./cmd/hourly

out/demo2:
	go build -o $@ ./cmd/demo2

apps: out/hourly out/demo2

.PHONY: apps out/hourly out/demo2

# ----------------
#  Go Maintenance
# ----------------

build:
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

fix-imports:
	find . -type f -iname "*.go" -exec goimports -w {} +

# ---------
#  Testing
# ---------

test:
	go test ./...

test-verbose:
	go test -v ./...

tmp:
	mkdir -p $@

tmp/coverage.data: | tmp
	go test -coverprofile=$@

.PHONY: tmp/coverage.data

tmp/coverage.html: tmp/coverage.data
	go tool cover -html=$< -o $@

coverage: tmp/coverage.html
	@echo "---"
	@echo "Open $^ in browser to view coverage info"
	@echo "---"

# ------
#  Misc
# ------

clean:
	go clean -x ./...
	rm -rf out
	rm -rf tmp

git_latest_version_tag := git describe --tags --match "v[0-9]*" --abbrev=0

# Make sure the tags are published and pushed to the public remote!
sync-package-proxy:
	GOPROXY=proxy.golang.org go list -m ${PROJECT_URL}@$(shell ${git_latest_version_tag})

.PHONY: build test fmt vet clean coverage sync-package-proxy
