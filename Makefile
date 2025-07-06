GENERATED_FILES += docs/adr/README.md

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

docs/adr/README.md: .adr-dir $(filter-out docs/adr/README.md,$(wildcard docs/adr/*.md))
	adr generate toc -i docs/adr/README.intro.md > "$@"

lint::
	vale --no-wrap --no-global .

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
