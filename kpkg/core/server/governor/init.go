package governor

import (
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/kpkg"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kstring"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(conf.Traverse("."))
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(kstring.PrettyJSONBytes(conf.Traverse(".")))
	})

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":         kpkg.Name(),
			"appID":        kpkg.AppID(),
			"appMode":      kpkg.AppMode(),
			"appVersion":   kpkg.AppVersion(),
			"tabbyVersion": kpkg.TabbyVersion(),
			"buildUser":    kpkg.BuildUser,
			"buildHost":    kpkg.BuildHost(),
			"buildTime":    kpkg.BuildTime(),
			"startTime":    kpkg.StartTime(),
			"hostName":     kpkg.HostName(),
			"goVersion":    kpkg.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}
