{
  pname,
  version,
  nerd ? false,
  release ? false,

  lib,
  stdenvNoCC,
  bited-build,
  zip,
  ...
}:

stdenvNoCC.mkDerivation {
  inherit pname version;
  src = ./.;

  buildPhase = ''
    runHook preBuild
    rm -rf out
    ${bited-build}/bin/bited-build ${lib.optionalString nerd "--nerd"}
    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share/fonts $out/share/consolefonts
    cp -r out/. $out/share/fonts
    ${lib.optionalString release ''
      ${zip}/bin/zip -r $out/share/fonts/${pname}_${version}.zip $out/share/fonts
    ''}
    mv out/*.psfu.gz $out/share/consolefonts
    runHook postInstall
  '';
}
