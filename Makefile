-include artifacts/make/adr.mk
-include artifacts/make/go.mk

artifacts/make/%.mk:
	curl -sf https://dogmatiq.io/makefiles/fetch | bash /dev/stdin $*
