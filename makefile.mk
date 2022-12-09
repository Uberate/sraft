# CHECK CODES
.PHONY: check
check:
	cargo check

# TEST CODES

.PHONY: test
test: check
	cargo test