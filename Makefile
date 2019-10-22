# Some useful make targets

ASSIGNMENTS = $(shell find ./assignment -name "q-*" -type d)

# Default target for happy markers
all: assignment

# Format go sources
fmt:
	go fmt ./...

# Run unit tests
test: fmt
test:
	go test -coverprofile=testcov.log ./...

# Display test coverage in a browser
cover: test
cover:
	go tool cover -html=testcov.log

# Build assignment binaries
$(ASSIGNMENTS):
	go build -o ./$@/main ./$@

# Build and run all assignment code for happy markers
assignment: fmt
assignment: test
assignment: $(ASSIGNMENTS)
assignment:
	echo "\nAssignment Question 1:\n"
	./assignment/q-1/main
	echo "\nAssignment Question 2:\n"
	./assignment/q-2/main

# Remove any generated files
clean:
	rm $(shell find ./problems -not -name "*.go" -type f) || :
	rm $(shell find ./assignment -not -name "*.go" -and -not -name "*.c" -and -not -name "*.h" -type f) || :
	rm $(shell find ./pdf -not -name "*.tex" -type f) || :
	rm comp-phys.zip testcov.log || :

# Build assignment report
pdf:
	cd ./pdf && pdflatex ./assignment.tex && pdflatex ./assignment.tex

# Build zip for submission
zip:
	rm comp-phys.zip || :
	zip --symlinks -r comp-phys.zip .

.PHONY: all fmt test clean pdf zip $(ASSIGNMENTS)
