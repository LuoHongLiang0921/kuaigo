package governor

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kstring"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)
var startTime= time.Now()
func init() {
	//HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
	//	encoder := json.NewEncoder(w)
	//	if r.URL.Query().Get("pretty") == "true" {
	//		encoder.SetIndent("", "    ")
	//	}
	//	encoder.Encode(conf.Traverse("."))
	//})

	HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Sub(startTime) < time.Second*5 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	HandleFunc("/runtime", func(w http.ResponseWriter, r *http.Request) {
		runtimeStates := map[string]string{
			"Goroutine":         strconv.Itoa(runtime.NumGoroutine()),
		}
		_ = jsoniter.NewEncoder(w).Encode(runtimeStates)
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
			"name":         pkg.GetAppName(),
			"appID":        pkg.GetAppID(),
			"appMode":      pkg.GetAppMode(),
			"appVersion":   pkg.GetAppVersion(),
			"tabbyVersion": pkg.GetTabbyVersion(),
			"buildUser":    pkg.GetBuildUser(),
			"buildHost":    pkg.GetBuildHost(),
			"buildTime":    pkg.GetBuildTime(),
			"startTime":    pkg.GetStartTime(),
			"hostName":     pkg.GetHostName(),
			"goVersion":    pkg.GetGoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}
