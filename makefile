clean:
	@rm -rf app/target

dev:
	@echo "Building dev version"
	@cd app/src; go build -o ../target/polling-dev main.go

build: clean
	@mkdir -p app/target
	@rsync -a app/src/* app/target --exclude *_test.go
	@cd app/target && GOOS=linux go build -o polling

package-only:
	@cd app/target && zip -qr ../source.zip .
	@rm -rf app/target

package: build
	@cd app/target && zip -qr ../source.zip .
	@rm -rf app/target
