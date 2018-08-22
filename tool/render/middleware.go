package render

//Based on https://github.com/lithammer/go-wiki

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var options struct {
	Dir       string
	Port      string
	CustomCSS string
	template  *template.Template
	git       bool
}

func Serve(folder, port string) error {
	options.Port = port
	options.Dir = folder
	logrus.Infoln("Serving dump from", options.Dir)

	// Parse base template
	var err error
	options.template, err = template.New("base").Parse(Template)
	if err != nil {
		logrus.Fatalln("Error parsing HTML template:", err)
	}

	// Check if the wiki folder is a Git repository
	options.git = IsGitRepository(options.Dir)
	if options.git {
		logrus.Debugln("Git repository found in directory")
	} else {
		logrus.Debugln("No git repository found in directory")
	}

	http.Handle("/api/diff/", commonHandler(DiffHandler))
	http.Handle("/", commonHandler(WikiHandler))

	logrus.Infof("Listening on: http://0.0.0.0:%s", options.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", options.Port), nil)
}

func commonHandler(next http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logrus.Warnf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		t0 := time.Now()
		next.ServeHTTP(w, r)
		logrus.Debugf("[%s] %q %v", r.Method, r.URL.String(), time.Now().Sub(t0))
	}

	return http.HandlerFunc(fn)
}
