build:
	hugo --gc --minify
	mkdir -p functions
	cd src/uuid && go build -o ../../functions/uuid ./...
	cd src/s1 && go build -o ../../functions/s1 ./...
	cd src/clipboard && go build -o ../../functions/clipboard ./...

