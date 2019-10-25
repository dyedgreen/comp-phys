# Code for Computational Physics at Imperial (fall of 2019)

## Build instructions

To build and run all assignment code with default data a simple `make` should suffice. (Alternatively
one can run `make assignment`). To generate the submitted reports run `make pdf` and to generate
the submitted ZIP archive run `make zip`.

(Note that the `zip` make target will create a zip with all files in the current directory, to generate
a tidy and small file, use `make clean zip`)


## Build prerequisites

This project uses Go. To build the go source, multiple open source compilers are available. See
[the Go website](https://golang.org) for more information on how to obtain a compiler for your
operating system and architecture.


## Contained folders

- `problems` ----- answers to problem sheets
- `assignment` --- answers to assignment problems
- `pkg` ---------- reusable packages


## License

Unless explicitly stated otherwise, the code in this repository is authored by _Tilman Roeder_ and
licensed under the MIT license (seen `LICENSE`).
