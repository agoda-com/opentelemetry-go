.PHONY: test-coverage

ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

TOOLS = $(CURDIR)/.tools

$(TOOLS):
	@mkdir -p $@
$(TOOLS)/%: | $(TOOLS)
	cd $(TOOLS_MOD_DIR) && \
	$(GO) build -o $@ $(PACKAGE)

GOCOVMERGE = $(TOOLS)/gocovmerge
$(TOOLS)/gocovmerge: PACKAGE=github.com/wadey/gocovmerge

.PHONY: tools
tools: $(GOCOVMERGE)


test-coverage: | $(GOCOVMERGE)
		@set -e; \
		for dir in $(ALL_COVERAGE_MOD_DIRS); do \
		  (cd "$${dir}" && go test -coverprofile=coverage.out ./... && \
		  go tool cover -html=coverage.out -o coverage.html); \
		done; \
		$(GOCOVMERGE) $$(find . -name coverage.out) > coverage.txt
