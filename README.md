# Code for Computational Physics at Imperial (fall of 2019)

## Build instructions

To build and run the project, use `make` or `make project-run`. To only build the project binary,
use `make project`.

To build and run all assignment code with default data a simple `make assignment` should suffice.

To generate the submitted reports run `make pdf` and to generate the submitted ZIP archive run `make zip`.

## PNG and other Graph Formats

There have been complaints about the graphs being generated as `.pdf` files. When running `make pdf`,
the graphs are generated in `.pdf` format to produce a high-quality report document. If the marker
wishes to inspect the graphs in some other format (e.g. `.png`), they are free to run `./main -graph`,
in the `/project` directory. (Run `./main` to see available options and make sure to build the project
first!)

## Running Tests

To run tests, use `make test`. (To see a code coverage breakdown, use `make cover`.)

## Build prerequisites

This project uses Go. To build the go source, multiple open source compilers are available. See
[the Go website](https://golang.org) for more information on how to obtain a compiler for your
operating system and architecture.

### Installation example macOS

If you are using homebrew on macOs, the installation should be as simple as:
```bash
brew install golang
```

### Installation example Ubuntu

If you are using Ubuntu, you might want to install Go using snap:
```bash
sudo snap install --classic go
```

### Notes on removing the Go installation

Any files generated by Go when building this project will either be located in this folder, or in
`$HOME/go` (unless you manually changed your `GOPATH`, run `go help gopath` for more information).

To remove everything, deinstall Go using your package manager, then delete `$HOME/go`.


## Contained folders

- `/assignment` --- answers to assignment problems
- `/pdf` ---------- LaTeX code for submission PDFs
- `/pkg` ---------- reusable packages
- `/problems` ----- answers to problem sheets
- `/project` ------ project binary code


## License

Unless explicitly stated otherwise, the code in this repository is authored by _Tilman Roeder_ and
licensed under the MIT license (seen `LICENSE`).
