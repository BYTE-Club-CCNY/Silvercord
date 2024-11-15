{ buildPythonPackage
, fetchPypi
, anthropic
, defusedxml
, langchain-core
, poetry-core
}:


buildPythonPackage rec {
  pname = "langchain_anthropic";
  version = "0.2.4";
  pyproject = true;
  doCheck = false;

  src = fetchPypi {
    inherit pname version;
    # hash = "sha256-MuesUeGHTEfhogST519b/Iiw/+r18a7WCRVH4a5Eu4U=";
    hash = "sha256-A4LUx7UjaDm3A/e3Kz4G3ku1vpkQSxk/cZrb40xJVis=";
  };

  nativeBuildInputs = [
    anthropic
    defusedxml
    langchain-core
    poetry-core
  ];
}