# Some useful make targets

all:
	echo "Not implemented ..."

fmt:
	go fmt ./pkg...
	go fmt ./problems/...

zip:
	rm comp-phys.zip || :
	zip --symlinks -r comp-phys.zip .

.PHONY: zip
