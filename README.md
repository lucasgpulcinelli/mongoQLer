# MongoQLer
MongoQLer is a simple graphical application that converts data from an oracle database to a mongoDB one. It was made as the final project for the university discipline of Database Laboratories.

## Made fully by
- [Lucas Eduardo Gulka Pulcinelli](https://github.com/lucasgpulcinelli)

### How to compile the code
First, install the go compiler and git, and on linux you will need the opengl and x11 development packages:
- In debian or ubuntu use `sudo apt install libgl1-mesa-dev xorg-dev golang git`;
- In fedora or other red hat based systems, use `sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel golang git`;
- In nix just use the shell.nix file provided via `nix-shell`;
- In windows, use the download link at [https://go.dev/dl](https://go.dev/dl/) (use the msi installer) or use chocolatey via `choco install golang`. Note that because of the C indirect dependencies, you will need the mingw C compiler as well.

Use `go version` in a terminal and make sure you are using go version 1.20 or later. If your package manager does not provide a recent enough version, check [https://go.dev/doc/install](https://go.dev/doc/install).

Then, use `go run .` to compile and run the application, getting all dependencies to do so.
Note that during the first compilation it will take much longer than for rebuilding.
