package apollo

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"github.com/LuoHongLiang0921/kuaigo/pkg/flag"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/url"
	"os"

	"github.com/philchia/agollo/v4"
)

// DataConfigApollo defines apollo scheme
const DataConfigApollo = "apollo"

type apolloConfigSource struct {
	client      agollo.Client
	namespace   string
	propertyKey string
	changed     chan struct{}
}

func RegisterConfigHandler() {
	manager.Register(DataConfigApollo, func() conf.ConfigSource {
		var (
			configAddr = flag.String("config")
		)
		if configAddr == "" {
			configAddr = os.Getenv("APOLLO_SERVER_ADDR")
			if configAddr == "" {
				klog.KuaigoLogger.Panic("new apollo configSource, configAddr is empty")
				return nil
			}
		}
		// configAddr is a string in this format:
		// apollo://ip:port?appId=XXX&cluster=XXX&namespaceName=XXX&key=XXX&accesskeySecret=XXX&insecureSkipVerify=XXX&cacheDir=XXX
		urlObj, err := url.Parse(configAddr)
		if err != nil {
			klog.KuaigoLogger.Panic("parse configAddr error", klog.FieldErr(err))
			return nil
		}
		apolloConf := agollo.Conf{
			AppID:              urlObj.Query().Get("appId"),
			Cluster:            urlObj.Query().Get("cluster"),
			NameSpaceNames:     []string{urlObj.Query().Get("namespaceName")},
			MetaAddr:           urlObj.Host,
			InsecureSkipVerify: true,
			AccesskeySecret:    urlObj.Query().Get("accesskeySecret"),
			CacheDir:           ".",
		}
		if urlObj.Query().Get("insecureSkipVerify") == "false" {
			apolloConf.InsecureSkipVerify = false
		}
		if urlObj.Query().Get("cacheDir") != "" {
			apolloConf.CacheDir = urlObj.Query().Get("cacheDir")
		}

		return NewConfigSource(&apolloConf, urlObj.Query().Get("namespaceName"), urlObj.Query().Get("key"), urlObj.Query().Get("verbose") == "true")
	})
}

// NewConfigSource creates an apolloConfigSource
func NewConfigSource(conf *agollo.Conf, namespace string, key string, verbose bool) conf.ConfigSource {
	client := agollo.NewClient(conf, agollo.WithLogger(&agolloLogger{Verbose: verbose}))
	ap := &apolloConfigSource{
		client:      client,
		namespace:   namespace,
		propertyKey: key,
		changed:     make(chan struct{}, 1),
	}
	ap.client.Start()
	ap.client.OnUpdate(
		func(event *agollo.ChangeEvent) {
			ap.changed <- struct{}{}
		})
	return ap
}

// ReadConfig reads config content from apollo
func (ap *apolloConfigSource) ReadConfig() ([]byte, error) {
	//value := ap.client.GetString(ap.propertyKey, agollo.WithNamespace(ap.namespace))
	value := ap.client.GetContent(agollo.WithNamespace(ap.namespace))
	return []byte(value), nil
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ap *apolloConfigSource) IsConfigChanged() <-chan struct{} {
	return ap.changed
}

// Close stops watching the config changed
func (ap *apolloConfigSource) Close() error {
	ap.client.Stop()
	close(ap.changed)
	return nil
}

type agolloLogger struct {
	// V 是否打印appolo
	Verbose bool
}

// Infof ...
func (l *agolloLogger) Infof(format string, args ...interface{}) {
	if l.Verbose {
		klog.KuaigoLogger.Infof(format, args...)
	}
}

// Errorf ...
func (l *agolloLogger) Errorf(format string, args ...interface{}) {
	klog.KuaigoLogger.Errorf(format, args...)
}
