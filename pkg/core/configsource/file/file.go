package file

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kfile"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kgo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// fileConfigSource file provider.
type fileConfigSource struct {
	path        string
	dir         string
	enableWatch bool
	changed     chan struct{}
	logger      *klog.Logger
}

// ConfigSourceFile defines file scheme
const ConfigSourceFile = "file"

// RegisterConfigHandler
// 	@Description  注册file
func RegisterConfigHandler() {
	manager.Register(ConfigSourceFile, func() conf.ConfigSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			configAddr = os.Getenv("CONFIG_FILE_ADDR")
			if configAddr == "" {
				klog.KuaigoLogger.Panic("new file configSource, configAddr is empty")
				return nil
			}
		}
		if !watchConfig {
			watchConfig = os.Getenv("TABBY_CONFIG_WATCH") == "true"
		}
		return NewConfigSource(configAddr, watchConfig)
	})
	manager.DefaultScheme = ConfigSourceFile
}

// NewConfigSource returns new fileConfigSourc.
func NewConfigSource(path string, watch bool) *fileConfigSource {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		klog.KuaigoLogger.Panic("new configsource", klog.Any("err", err))
	}

	dir := kfile.CheckAndGetParentDir(absolutePath)
	ds := &fileConfigSource{path: absolutePath, dir: dir, enableWatch: watch}
	if watch {
		ds.changed = make(chan struct{}, 1)
		kgo.Go(ds.watch)
	}
	return ds
}

// ReadConfig ...
func (fp *fileConfigSource) ReadConfig() (content []byte, err error) {
	return ioutil.ReadFile(fp.path)
}

// Close ...
func (fp *fileConfigSource) Close() error {
	close(fp.changed)
	return nil
}

// IsConfigChanged ...
func (fp *fileConfigSource) IsConfigChanged() <-chan struct{} {
	return fp.changed
}

func (fp *fileConfigSource) getContext() context.Context {
	return context.TODO()
}

// Watch file and automate update.
func (fp *fileConfigSource) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		klog.KuaigoLogger.Fatal("new file watcher", klog.FieldMod("file configsource"), klog.Any("err", err))
	}

	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				klog.KuaigoLogger.Debug("read watch event",
					klog.FieldMod("file configsource"),
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
				klog.KuaigoLogger.Error("read watch error", klog.FieldMod("file configsource"), klog.Any("err", err))
			}
		}
	}()

	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
