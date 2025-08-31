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
  helmquilt <apply|check|diff> [flags]

Flags:
  -c, --config string   path to the config file (default "./helmquilt.yaml")
  -f, --force           force run (ignore lock file)
  -h, --help            help for helmquilt
  -q, --quiet           Silence the logs
  -r, --repack          Repack the chart as a tarball
```

and check the [example](./example/) folder for a sample configuration file.

### Example workflow

1. Create a `helmquilt.yaml` with the chart(s) you want to fetch
2. Run `helmquilt apply` to fetch the chart(s)
3. Make your changes to the chart(s) that got pulled from the previous command
4. Run `helmquilt diff`
5. Use the output of the previous command to create patch files
6. Update the `helmquilt.yaml` config file to include the patches
7. Run `helmquilt apply` to update the lock file with the new checksums

If latter you want to add more changes, just make the changes you need, run `helmquilt diff`, create new patch files, add them to the config and run `helmquilt apply` to update the lock file.
If you keep your patched chart(s) in git, you can also use `git diff` instead of `helmquilt diff` to create the patche files.

In the `example` folder the `example/patches/coredns1.patch` was created using `git diff` and `example/patches/flux.patch` was created using `helmquilt diff`.
