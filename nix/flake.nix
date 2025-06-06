{
  description = "Dev shell for Silvercord Discord Bot";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable"; };

  outputs = { self, nixpkgs, ... }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
    python = pkgs.python312.override {
      self = python;
      packageOverrides = pyfinal: pyprev: {
        langchain_anthropic = pyfinal.callPackage ./langchain-anthropic.nix {defusedxml = pyfinal.callPackage ./defusedxml.nix {} ;};
        ratemyprofessor = pyfinal.callPackage ./ratemyprofessor.nix {};
      };
    };
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
        python312Packages.beautifulsoup4
        python312Packages.langchain
        python312Packages.langchain-community
        python312Packages.langchain-openai
        python312Packages.langchain-chroma
        python312Packages.langchain-ollama
        python312Packages.requests
	      python312Packages.anthropic
	      python312Packages.openai
        (python.withPackages (python-pkgs: [
          # select Python packages here
          python-pkgs.langchain_anthropic
          python-pkgs.ratemyprofessor
        ]))
        nodejs_22
        nodemon
      ];
#         shellHook = "	
#          source \"$PWD/.env\" || exit 1 ";
    };
  };
}
