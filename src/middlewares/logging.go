package middlewares

import (
	"history_anime/src/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func Logging(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		log := logger.New()
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Request received")

		next(w, r, params)
	}
}
