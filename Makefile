# Some useful make targets

ASSIGNMENTS = $(shell find ./assignment -name "q-*" -type d)

all: assignment

fmt:
	go fmt ./pkg/...
	go fmt ./problems/...
	go fmt ./assignment/...

test:
	go test ./pkg/...

$(ASSIGNMENTS):
	go build -o ./$@/main ./$@

assignment: fmt
assignment: test
assignment: $(ASSIGNMENTS)
assignment:
	echo "\nAssignment outputs:\n"
	./assignment/q-1/main

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
