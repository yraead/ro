---
name: NewFSListener
slug: newfslistener
sourceRef: plugins/fsnotify/source.go#L25
type: plugin
category: fsnotify
signatures:
  - "func NewFSListener(paths ...string)"
playUrl: ""
variantHelpers:
  - plugin#fsnotify#newfslistener
similarHelpers: []
position: 0
---

Creates an observable that monitors file system events for specified paths.

```go
import (
    "fmt"

    "github.com/fsnotify/fsnotify"
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

obs := rofsnotify.NewFSListener("/tmp")

sub := obs.Subscribe(ro.NewObserver(
    func(event fsnotify.Event) {
        switch event.Op {
        case fsnotify.Create:
            fmt.Println("File was created")
        case fsnotify.Write:
            fmt.Println("File was written to")
        case fsnotify.Remove:
            fmt.Println("File was removed")
        }
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
defer sub.Unsubscribe()

// When file events occur:
// File was created
// File was written to
// File was removed
```