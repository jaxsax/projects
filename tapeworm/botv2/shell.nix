{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go
    gotools
    gopls
    go-outline
    gocode
    gopkgs
    gocode-gomod
    godef
    golint
    go-migrate
    wire
    sqlite
    yq
    litecli
    nodejs-18_x
    nodePackages.typescript-language-server
  ];
}
