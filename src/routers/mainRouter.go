package routers

import "github.com/julienschmidt/httprouter"

func Router() *httprouter.Router {

	routers := httprouter.New()
	AuthRoute(routers)
	AnimeRoute(routers)

	return routers
}
