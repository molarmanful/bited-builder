{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    let
      version = builtins.readFile ./VERSION;
    in

    {
      templates.default = {
        path = ./template;
        description = "bited font project with bited-utils";
      };
    }
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        o = {
          inherit version;
          vendorHash = "sha256-/EcjKt5IBY1tGOFRiL67LovK2y9J+5WaIeCWaNcjrFA=";
        } // self.packages.${system};
      in
      {

        packages = rec {
          default = bited-utils;
          bited-utils = pkgs.callPackage ./. o;
          bitsnpicas = pkgs.callPackage ./bitsnpicas.nix { };
          bited-build = pkgs.callPackage ./bited-build o;
          bited-img = pkgs.callPackage ./bited-img o;
          bited-scale = pkgs.callPackage ./bited-scale o;
          bited-clr = pkgs.callPackage ./bited-clr o;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nil
            nixd
            nixfmt-rfc-style
            statix
            deadnix
            taplo
            go
            gopls
            gotools
            golines
            errcheck
            marksman
            mdformat
            python313Packages.mdformat-gfm
            python313Packages.mdformat-gfm-alerts
          ];
        };

      }
    );
}
