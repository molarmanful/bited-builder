# FIXME

A cool font waiting to be built.

## Getting Started

This project's default font name is `FIXME`, which you should replace with your
own font's name. The project's default author is
[ghost](https://github.com/ghost), which you should also replace with your name.

> [!TIP]
> A tool like [rgr](https://github.com/acheronfail/repgrep) can make this
> process easier!

### Updating

To update all inputs:

```
nix flake update
```

To only update bited-utils:

```
nix flake update bited-utils
```

### Licensing

The project comes with SIL Open Font License 1.1 by default. Feel free to alter
this to your liking.

### Versioning

By default, Nix flake reads versions via the `VERSION` file.
`.github/workflows/pub.yml` bumps the version based on Github releases. These
can all be modified to fit your preferred versioning scheme.
