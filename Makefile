# Some useful make targets

ASSIGNMENTS = $(shell find ./assignment -name "q-*" -type d)

all: assignment

fmt:
	go fmt ./...

test: fmt
test:
	go test ./...

$(ASSIGNMENTS):
	go build -o ./$@/main ./$@

assignment: fmt
assignment: test
assignment: $(ASSIGNMENTS)
assignment:
	echo "\nAssignment Question 1:\n"
	./assignment/q-1/main
	echo "\nAssignment Question 2:\n"
	./assignment/q-2/main

clean:
	rm $(shell find ./problems -not -name "*.go" -type f) || :
	rm $(shell find ./assignment -not -name "*.go" -and -not -name "*.c" -and -not -name "*.h" -type f) || :
	rm $(shell find ./pdf -not -name "*.tex" -type f) || :

pdf:
	cd ./pdf && pdflatex ./assignment.tex && pdflatex ./assignment.tex

zip:
	rm comp-phys.zip || :
	zip --symlinks -r comp-phys.zip .

.PHONY: all fmt test clean pdf zip $(ASSIGNMENTS)
