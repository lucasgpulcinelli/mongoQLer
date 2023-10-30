{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [ go glfw pkg-config xorg.libXcursor xorg.libX11 xorg.libXrandr xorg.libXinerama xorg.libXi xorg.libXxf86vm ];
}

