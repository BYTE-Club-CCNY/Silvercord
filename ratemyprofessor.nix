{ buildPythonPackage,
# fetchPypi, 
fetchFromGitHub
}:
buildPythonPackage rec {
  pname = "RateMyProfessorAPI";
  version = "1.3.6";
  # format = "setuptools";
  # pyproject = true;
  # doCheck = false;


  src = fetchFromGitHub {
    owner = "nobelz";
    repo = "ratemyprofessorapi";
    rev = "9bd5fb359fe61ac3461ad7fbade5221e5d9b4b4d";
    sha256 = "sha256-fDFmDh9t+c6fbCxPjoVq7Ci5sQQVMplZ2UBf0gxmqTM="; # TODO
  };
  # src = fetchPypi {
  #   inherit pname version;
  #   # hash = "sha256-MuesUeGHTEfhogST519b/Iiw/+r18a7WCRVH4a5Eu4U=";
  #   hash = "sha256-Vitt69aT0uJ+vMPsafHfUvN+oFHwI+iNIfuVMV1Ht8E=";
  # };
}