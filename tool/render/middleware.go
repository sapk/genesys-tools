package render

//Based on https://github.com/lithammer/go-wiki

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("Serving dump from: %s", options.Dir)

	// Parse base template
	var err error
	options.template, err = template.New("base").Parse(Template)
	if err != nil {
		log.Fatal().Msgf("Error parsing HTML template: %v", err)
	}

	// Check if the wiki folder is a Git repository
	options.git = IsGitRepository(options.Dir)
	if options.git {
		log.Debug().Msg("Git repository found in directory")
	} else {
		log.Debug().Msg("No git repository found in directory")
	}

	http.Handle("/api/diff/", commonHandler(DiffHandler))
	http.Handle("/", commonHandler(WikiHandler))

	log.Info().Msgf("Listening on: http://0.0.0.0:%s", options.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", options.Port), nil)
}

func commonHandler(next http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Warn().Msgf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		t0 := time.Now()
		next.ServeHTTP(w, r)
		log.Debug().Msgf("[%s] %q %v", r.Method, r.URL.String(), time.Now().Sub(t0))
	}

	return http.HandlerFunc(fn)
}
