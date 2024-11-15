{
  lib,
  buildPythonPackage,
  # fetchFromGitHub,
  fetchPypi,
  setuptools,
  lxml,
  python,
}:

buildPythonPackage rec {
  pname = "defusedxml";
  version = "0.7.1";
  format = "setuptools";
  # pyproject = true;

  # src = fetchFromGitHub {
  #   owner = "tiran";
  #   repo = "defusedxml";
  #   rev = "refs/tags/v${version}";
  #   hash = "sha256-X88A5V9uXP3wJQ+olK6pZJT66LP2uCXLK8goa5bPARA=";
  # };

  src = fetchPypi {
    inherit pname version;
    sha256 = "1bb3032db185915b62d7c6209c5a8792be6a32ab2fedacc84e01b52c51aa3e69";
  };

  checkPhase = ''
    ${python.interpreter} tests.py
  '';

  pythonImportsCheck = [ "defusedxml" ];

  meta = with lib; {
    description = "Python module to defuse XML issues";
    homepage = "https://github.com/tiran/defusedxml";
    license = licenses.psfl;
    maintainers = with maintainers; [ fab ];
  };
}