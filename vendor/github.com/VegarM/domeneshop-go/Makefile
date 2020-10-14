generate-domeneshop-client: 
	openapi-generator-cli generate \
		--generator-name go \
		--input-spec ./api/swagger.json \
		--output ./ \
		--package-name domeneshop \
		--git-host github.com \
		--git-user-id VegarM \
		--git-repo-id domeneshop-go

all: generate-domeneshop-client
