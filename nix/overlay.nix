self: super:
{
  buildGoModule = super.buildGo126Module;

  gnutls = super.gnutls.overrideAttrs (old: {
    configureFlags = (old.configureFlags or [ ]) ++ [ "--disable-doc" ];
    outputs = builtins.filter (o: o != "devdoc" && o != "man") (old.outputs or [ "out" ]);
    # This override is never in cache.nixos.org, so gnutls always builds from
    # source, and its gnulib tests fail on some build hosts.
    doCheck = false;
  });

  libseccomp = super.libseccomp.overrideAttrs (x: {
    doCheck = false;
    dontDisableStatic = true;
  });
}
