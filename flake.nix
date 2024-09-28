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
          mkdir -p $out/share/xk/userscripts

          # install core
          cp src/bin/xettelkasten $out/bin/xk
          cp src/lib/* $out/lib/xk/
          cp -r src/etc/* $out/etc/xk/

          # install bash userscripts (copy but remove .sh)
          for file in src/userscripts/*.sh; do
              cp "$file" "$out/share/xk/userscripts/$(basename "$file" .sh)"  # Copy and remove the .sh extension
          done
        '';

        meta = with pkgs.lib; {
          description = "xettelkasten core";
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
