Meta information:

- https://golang.org/cmd/go/#hdr-Remote_import_paths
- https://github.com/golang/tools/tree/0a99049195aff55f007fc4dfd48e3ec2b4d5f602/go/vcs
```go
import "golang.org/x/tools/go/vcs"

repoRoot, err := vcs.RepoRootForImportDynamic("pkname", false)
```
- https://github.com/golang/go/tree/master/src/cmd/go/internal/get
- https://sagikazarmark.hu/blog/vanity-import-paths-in-go/
