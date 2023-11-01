BINARY_NAME=chirpy
DATABASE_NAME=database.json

build:
	go build -o ${BINARY_NAME}-darwin main.go

run:build
	./%{BINARY_NAME}

clean:
	echo "Removing database" 
	rm ${DATABASE_NAME}

testDatabase:
	go test internal/database/database_test.go internal/database/database.go
