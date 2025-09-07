# Helmquilt

Helmquilt is a tool for managing helm patches. It allows you to pull helm charts from oter locations, including other git repositories and apply patches to them.
Instead of forking the whole chart you can now just maintain a set of patches and repackage external charts to your needs.

Like with [quilt](https://www.man7.org/linux/man-pages/man1/quilt.1.html), the key philosophical concept is that your primary working material is patches.

Some times, as a platform engineer, you need to manage a set of costomisations on top of upstream helm charts, this can make keeping them up-to-date hard.
Helmquilt allows you to fetch the upstream helm chart at a specified version and apply your changes with ease, as long as your patches still apply, you can update the chart by just pulling a newer version and applying your changes again.
This becomes specially useful in environments where mono-repos are used and you can't really fork an upstream helm chart. This allows you to keep a modified copy of a helm chart and still be able to pull upstream changes.
The tool does suppot `git` as a source, so, in theory, it could be used with anything and not just helm charts...

## Usage

see `helmquilt -h`

```console
$ helmquilt -h
helmquilt is a tool for managing helm package patches

Usage:
  helmquilt [command]

Available Commands:
  apply       Apply the helmquilt config, fetch charts and apply patches
  check       Read the helmquilt config and lock files and check if everything is up-to-date
  completion  Generate the autocompletion script for the specified shell
  diff        Check if changes were made and return the differences
  help        Help about any command

Flags:
  -h, --help    help for helmquilt
  -q, --quiet   Silence the logs

Use "helmquilt [command] --help" for more information about a command.
```

and check the [example](./example/) folder for a sample configuration file.

### Example workflows

#### Start managing some charts and patches

1. Create a `helmquilt.yaml` with the chart(s) you want to fetch
2. Run `helmquilt apply` to fetch the chart(s)
3. Make your changes to the chart(s) that got pulled from the previous command
4. Run `helmquilt diff` and validate the output
5. Run `helmquilt diff` again, this time with the `-w` flag so it saves the patches as a file and adds them to the config.

If latter you want to add more changes, just make the changes you need, run `helmquilt diff`, create new patch files, add them to the config and run `helmquilt apply` to update the lock file.
If you keep your patched chart(s) in git, you can also use `git diff` instead of `helmquilt diff` to create the patch files.

In the `example` folder the `example/patches/coredns1.patch` was created using `git diff` and `example/patches/flux.patch` was created using `helmquilt diff`.

#### Update existing modified charts to a new upstream version

If you already have some modified charts that you'd like to upgrade without loosing your changes, you can try:

1. Create a `helmquilt.yaml` file matching the charts and versions you already have.
2. Run `helmquilt diff`, this shold give you the diffs with your changes.
3. Save the diffs into patch files, you can use the `-w` option in `helmquilt diff`
4. Update the `helmquilt.yaml` config file with the versions you want to update to.
5. Run `helmquilt apply` to update the lock file with the new checksums

As long as your patches still apply to the new versions, you shold get the new versions of the charts with your changes on top.

#### Combine multiple patches into one

If you've been using the tool for a while and have accumulated many changes in separate patch files, you can try the following to combine all changes into a single patch.

1. run `helmquilt apply`, just to be sure everythig is up-to-date
2. delete the patch files and remove them from the `helmquilt.yaml` config
3. without any additional changes, run `helmquilt diff` with the `-w` option

This shold generate a single path file (per chart) with all the combined changes
