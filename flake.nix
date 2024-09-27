{
  description = "A simple latex centric zettelkasten written in bash";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }: {
    packages.x86_64-linux = let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      default = pkgs.stdenv.mkDerivation {
        pname = "xettelkasten";
        version = "1.0.0";

        src = ./.;

        buildInputs = [ pkgs.bash ];

        installPhase = ''
          mkdir -p $out/bin
          mkdir -p $out/lib/xk
          mkdir -p $out/etc/xk
          mkdir -p $out/share/xk

          cp src/bin/xettelkasten $out/bin/xk
          cp src/lib/* $out/lib/xk/
          cp -r src/etc/* $out/etc/xk/
          cp -r src/share/* $out/share/xk/
        '';

        meta = with pkgs.lib; {
          description = "A simple latex centric zettelkasten written in bash";
          license = licenses.mit;
        };
      };
    };

    devShells.x86_64-linux = let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in pkgs.mkShell {
      buildInputs = [ pkgs.bash ];
    };
  };
}
