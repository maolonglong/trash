{
  description = "trash - Move FILE(s) to the trash";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nur.url = "github:nix-community/NUR";
    maolonglong-nur.url = "github:maolonglong/nur-packages";
    maolonglong-nur.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    nur,
    maolonglong-nur,
    ...
  }: let
    trashVersion =
      if (self ? shortRev)
      then self.shortRev
      else "dev";
    supportSystems = [
      "x86_64-darwin"
      "i686-darwin"
      "aarch64-darwin"
      "armv7a-darwin"
    ];
    overlays = [
      (final: prev: {
        nur = import nur {
          nurpkgs = prev;
          pkgs = prev;
          repoOverrides = {
            maolonglong = import maolonglong-nur {pkgs = prev;};
          };
        };
      })
    ];
  in
    flake-utils.lib.eachSystem supportSystems (
      system: let
        pkgs = import nixpkgs {inherit system;};
        lib = pkgs.lib;
      in rec {
        packages.trash = pkgs.buildGoModule {
          pname = "trash";
          version = trashVersion;
          src = lib.cleanSource self;
          vendorHash = "sha256-3PnXB8AfZtgmYEPJuh0fwvG38dtngoS/lxyx3H+rvFs=";
          ldflags = ["-s" "-w"];
          checkFlags = ["-skip=TestRm"];
        };
        packages.default = packages.trash;

        apps.trash = flake-utils.lib.mkApp {
          drv = packages.trash;
        };
        apps.default = apps.trash;
      }
    )
    // flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system overlays;};
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs =
            (with pkgs; [
              just
              go
              golines
              gosimports
            ])
            ++ (with pkgs.nur.repos.maolonglong; [
              gofumpt
              skywalking-eyes
            ]);
        };
      }
    );
}
