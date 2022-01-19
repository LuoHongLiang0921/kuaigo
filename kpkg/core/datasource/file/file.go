package file

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kfile"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kgo"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// fileDataSource file provider.
type fileDataSource struct {
	path        string
	dir         string
	enableWatch bool
	changed     chan struct{}
	logger      *klog.Logger
}

// NewDataSource returns new fileDataSource.
func NewDataSource(path string, watch bool) *fileDataSource {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		klog.Panic(context.TODO(), "new datasource", klog.Any("err", err))
	}

	dir := kfile.CheckAndGetParentDir(absolutePath)
	ds := &fileDataSource{path: absolutePath, dir: dir, enableWatch: watch}
	if watch {
		ds.changed = make(chan struct{}, 1)
		kgo.Go(ds.watch)
	}
	return ds
}

// ReadConfig ...
func (fp *fileDataSource) ReadConfig() (content []byte, err error) {
	return ioutil.ReadFile(fp.path)
}

// Close ...
func (fp *fileDataSource) Close() error {
	close(fp.changed)
	return nil
}

// IsConfigChanged ...
func (fp *fileDataSource) IsConfigChanged() <-chan struct{} {
	return fp.changed
}

func (fp *fileDataSource) getContext() context.Context {
	return context.TODO()
}

// Watch file and automate update.
func (fp *fileDataSource) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Fatal(context.TODO(), "new file watcher", klog.FieldMod("file datasource"), klog.Any("err", err))
	}

	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				klog.Debug(fp.getContext(), "read watch event",
					klog.FieldMod("file datasource"),
					klog.String("event", filepath.Clean(event.Name)),
					klog.String("path", filepath.Clean(fp.path)),
				)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if event.Op&writeOrCreateMask != 0 && filepath.Clean(event.Name) == filepath.Clean(fp.path) {
					log.Println("modified file: ", event.Name)
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				// log.Println("error: ", err)
				klog.Error(fp.getContext(), "read watch error", klog.FieldMod("file datasource"), klog.Any("err", err))
			}
		}
	}()

	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
