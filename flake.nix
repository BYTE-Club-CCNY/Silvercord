{
  description = "dev shell for LeetBoard discord bot";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable"; };

  outputs = { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.x86_64-linux.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          python312Packages.python
          python312Packages.pip
          python312Packages.annotated-types
          python312Packages.anyio
          python312Packages.certifi
          python312Packages.colorama
          python312Packages.distro
          python312Packages.h11
          python312Packages.httpcore
          python312Packages.httpx
          python312Packages.idna
          python312Packages.jiter
          python312Packages.openai
          python312Packages.pydantic
          python312Packages.pydantic-core
          python312Packages.python-dotenv
          python312Packages.sniffio
          python312Packages.tqdm
          python312Packages.typing-extensions
          # python312Packages.playwright
          python312Packages.requests
          nodejs_22
          nodemon
          # # playwright packages
          # playwright-driver.browsers
        ];
        shellHook = "	
          source \"$PWD/.env\" || exit 1 ";
      };
    };
}
