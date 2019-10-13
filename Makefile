# Some useful make targets

all:
	echo "Not implemented ..."

zip:
	rm comp-phys.zip || :
	zip -r comp-phys.zip .

.PHONY: zip
