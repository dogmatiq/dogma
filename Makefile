GENERATED_FILES += docs/adr/README.md

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

docs/adr/README.md: .adr-dir $(filter-out docs/adr/README.md,$(wildcard docs/adr/*.md))
	adr generate toc -i docs/adr/README.intro.md > "$@"

.vale/vale.touch: .vale.ini
	vale sync
	@touch $@

lint:: .vale/vale.touch
	vale --no-wrap --no-global --glob='!.*/**' .

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
