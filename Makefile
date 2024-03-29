.PHONY: install
install:
	@for app_dir in cmd/*; do \
  		app_name=$$(basename $$app_dir); \
  		echo "Installing $$app_name..."; \
		go build -o "${GOPATH}/bin/$$app_name" $$app_dir/*.go; \
  	done