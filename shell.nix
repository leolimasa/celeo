{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
  ]; 

  # Runs a command after shell is started
  shellHook = ''
    export ENV_NAME=celeo
  '';
}
