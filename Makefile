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
	go test -coverprofile=testcov.log ./pkg/... ./assignment/comply/...

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
	mkdir -p ./pdf/images/
	echo "\nAssignment Question 1:\n"
	./assignment/q-1/main
	echo "\nAssignment Question 2:\n"
	./assignment/q-2/main
	echo "\nAssignment Question 3:\n"
	cd ./assignment/q-3/ && ./main
	cp ./assignment/q-3/plot.pdf ./pdf/images/assignment-q-3.pdf
	echo "\nAssignment Question 4:\n"
	cd ./assignment/q-4/ && ./main
	cp ./assignment/q-4/plot.pdf ./pdf/images/assignment-q-4.pdf
	echo "\nAssignment Question 5:\n"
	cd ./assignment/q-5/ && ./main
	cp ./assignment/q-5/plot-a.pdf ./pdf/images/assignment-q-5-a.pdf
	cp ./assignment/q-5/plot-b.pdf ./pdf/images/assignment-q-5-b.pdf
	cp ./assignment/q-5/plot-c.pdf ./pdf/images/assignment-q-5-c.pdf

# Remove any generated files
clean:
	rm $(shell find ./problems -not -name "*.go" -type f) || :
	rm $(shell find ./assignment -not -name "*.go" -and -not -name "*.c" -and -not -name "*.h" -type f) || :
	rm $(shell find ./pdf -not -name "*.tex" -and -no -name "*.bib" -type f) || :
	rm comp-phys.zip testcov.log || :

# Build assignment report
pdf: assignment
pdf:
	cd ./pdf && pdflatex ./assignment.tex && bibtex ./assignment.aux && pdflatex ./assignment.tex

# Build zip for submission
zip:
	rm comp-phys.zip || :
	zip --symlinks -r comp-phys.zip .

.PHONY: all fmt test clean pdf zip $(ASSIGNMENTS)
