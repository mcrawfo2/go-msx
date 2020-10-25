# MSX Resource Module

MSX Resource manages locating and accessing files from the source, staging, and runtime filesystems in a consistent fashion.

## Filesystem

To correctly use resources, it is first important to understand the resource filesystem, and how it is used to locate files during development and inside containers.

The resource filesystem contains one or more of the following layers, if found:

- **production** - rooted in the Docker image at `/var/lib/${app.name}`
- **staging** - rooted at `dist/root/var/lib/${app.name}` underneath the **source** root
- **source** - rooted at the folder containing the repository's `go.mod`

The resource filesystem will attempt to locate each of these folders and if found, will search it for your resource references.

## Resource References

The primary data type of the MSX resource module is the resource reference.  It represents the resource file subpath.  All resource paths use the forward-slash (`/`) as the path component separator. 

Two types of paths can be used:

- **relative** - No leading forward-slash (`data/my-resource.json`): File path is relative to the code file consuming the reference.
- **absolute** - Leading forward-slash (`/internal/migrate/resource.json`): File path is relative to the resource filesystem root.

### Obtaining a Single Resource Reference

To work with a resource you must first create a reference to it using the `resource.Reference` function:

```go
func processMyResource(ctx context.Context) error {
  myResourceRef := resource.Reference("data/my-resource.json")
}
```

This returns a `resource.Ref` object pointing to the specified path.

### Obtaining Multiple Resource References

To retrieve multiple resource references using a glob pattern you can call the `resource.References` function:

```
func processMyResources(ctx context.Context) error {
  myResourceRefs := resource.References("data/*.json")
}
```

This returns a `[]resource.Ref` slice with an entry for each matching resource.

## Consuming Resources

Once you have obtained one or more resource references, you can access their contents using one of its methods.

### JSON

To read in the contents of the resource as JSON and unmarshal it to an object, use the `Unmarshal()` method:

```go
var myResourceContents MyResourceStruct
err := resource.Reference("data/my-resource.json").Unmarshal(&myResourceContents)
```

### Bytes

To read in the contents of the resource as a  `[]byte`, use the `ReadAll()` method:

```go
data, err := resource.Reference("data/my-resource.json").ReadAll()
```

### http.File

To open the file and return an `http.File`, use the `Open()` method:

```go
file, err := resource.Reference("data/my-resource.json").Open()
```

Note that `http.File` also meets the requirements of the `io.ReadCloser` interface, and can therefore be used with `io`.

