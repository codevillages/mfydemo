HTTPS_PROXY ?= http://127.0.0.1:7890
HTTP_PROXY  ?= http://127.0.0.1:7890
ALL_PROXY   ?= socks5://127.0.0.1:7890

.PHONY: test-user-create
test-user-create:
	HTTPS_PROXY=$(HTTPS_PROXY) HTTP_PROXY=$(HTTP_PROXY) ALL_PROXY=$(ALL_PROXY) \
	GF_GCFG_FILE=./config/config.test.yaml \
	go test ./module_user/internal/integration -run TestCreateUser -v
