{
  description = "A simple latex-centric zettelkasten written in bash";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }: {
    overlay = (final: prev: {
        xk = self.packages.x86_64-linux.default;
      });
    packages.x86_64-linux = let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      default = pkgs.buildGoModule {
        pname = "xettelkasten";
        version = "1.1.0";

        src = ./.;
        modRoot = ./.;
        vendorHash = null;

        buildInputs = [
          pkgs.bash
          pkgs.pdf2svg
          pkgs.texliveFull
        ];

        buildPhase = ''
          go build -o $out/share/xk/userscripts/genrefs ./src/userscripts-go/cmd/genrefs
          go build -o $out/share/xk/userscripts/gencards ./src/userscripts-go/cmd/gencards
          go build -o $out/share/xk/userscripts/syncanki ./src/userscripts-go/cmd/syncanki
        '';

        installPhase = ''
          mkdir -p $out/bin
          mkdir -p $out/lib/xk
          mkdir -p $out/etc/xk
          mkdir -p $out/share/xk/userscripts

          # core
          cp src/bin/xettelkasten $out/bin/xk
          cp src/lib/* $out/lib/xk/
          cp -r src/etc/* $out/etc/xk/

          # bash userscripts
          cp -r src/userscripts-bash/* $out/share/xk/userscripts
        '';

        meta = with pkgs.lib; {
          description = "xettelkasten core";
          license = licenses.mit;
        };
      };
    };
  };
}
