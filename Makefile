.PHONY: run test fmt

build:
	@echo "==> Building..."
	@sh -c "'$(CURDIR)/scripts/makefile.sh' build"

test:
	@echo "==> Testing..."
	@sh -c "'$(CURDIR)/scripts/makefile.sh' test -v -cover -tags=\'test local\' -count 1 ./..."

test_light:
	@echo "==> Testing..."
	@sh -c "'$(CURDIR)/scripts/makefile.sh' test -tags=\'test local\' -count 1 ./..."

fmt:
	@echo "==> Formatting go sources..."
	@sh -c "'$(CURDIR)/scripts/makefile.sh' fmt"

# Arguments have priority
#FOO=hoge
#export FOO
#testecho:
#	@echo $(FOO)

