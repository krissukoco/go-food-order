dir ?= *

protogen:
	@for i in ./api/proto/*_service; do \
		for d in $$i/*.proto; do \
			f=$$(basename $$d); \
			dirname=$${f%.*}; \
			echo "Generating go files from '$$d' ..."; \
			rm -rf $$i/$$dirname; \
			mkdir -p $$i/$$dirname; \
			protoc -I ./api/proto --go_out=$$i/$$dirname --go-grpc_out=$$i/$$dirname $$d; \
		done; \
	done

openapi:
	@for d in ./api/openapi/$(dir); do \
		packagename="$$(basename $$d)_http_api"; \
		for f in $$d/*.yaml; do \
			out=$${f%.yaml}.oapi.go; \
			echo "Generating OpenAPI Go file for '$$f' ..."; \
			oapi-codegen -o $$out -package $$packagename -generate types,fiber,spec $$f; \
		done; \
	done
	@echo "OpenAPI gen finished"
	@echo ""