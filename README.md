# Helmquilt

Helmquilt is a tool for managing helm patches. It allows you to pull helm charts from oter locations, including other git repositories and apply patches to them.
Instead of forking the whole chart you can now just maintain a set of patches and repackage external charts to your needs.

Like with [quilt](https://www.man7.org/linux/man-pages/man1/quilt.1.html), the key philosophical concept is that your primary working material is patches.

Some times, as a platform engineer, you need to manage a set of costomisations on top of upstream helm charts, this can make keeping them up-to-date hard.
Helmquilt allows you to fetch the upstream helm chart at a specified version and apply your chagnes with ease, as long as your patches still apply, you can update the chart by just pulling a newer version and applying your chagnes again.

## Usage

see `helmquilt -h`

```console
$ helmquilt -h
helmquilt is a tool for managing helm package patches

Usage:
  helmquilt <apply|check> [flags]

Flags:
  -c, --config string   path to the config file (default "./helmquilt.yaml")
  -f, --force           force run
  -h, --help            help for helmquilt
  -r, --repack          Repack the chart as a tarball
```

and check the [example](./example/) folder for a sample configuration file.
