package routes

import (
	"github.com/go-chi/chi"
	"github.com/musab-olurode/lis_backend/controllers"
	"github.com/musab-olurode/lis_backend/middlewares"
)

func ConfigureRouter() *chi.Mux {
	router := chi.NewRouter()

	// public routes
	router.Post("/auth/login", controllers.Login)
	router.Post("/auth/register", controllers.Register)
	router.Get("/auth/refresh", controllers.RefreshToken)

	router.Get("/posts", controllers.GetPosts)
	router.Get("/posts/{postID}", controllers.GetPost)
	router.Get("/posts/slug/{postSlug}", controllers.GetPostBySlug)

	router.Get("/materials", controllers.GetMaterials)
	router.Get("/materials/{materialID}", controllers.GetMaterial)

	router.Get("/events", controllers.GetEvents)
	router.Get("/events/upcoming", controllers.GetUpcomingEvents)
	router.Get("/events/{eventID}", controllers.GetEvent)

	// authenticated routes
	router.Group(func(router chi.Router) {
		router.Use(middlewares.Authenticated)

		router.Get("/auth/me", controllers.GetLoggedInUser)
		router.Post("/upload", controllers.UploadFile)
		router.Delete("/upload", controllers.DeleteFile)
	})

	// admin routes
	router.Group(func(router chi.Router) {
		router.Use(middlewares.Authenticated, middlewares.Admin)

		router.Post("/posts", controllers.CreatePost)
		router.Put("/posts/{postID}", controllers.UpdatePost)
		router.Delete("/posts/{postID}", controllers.DeletePost)

		router.Post("/materials", controllers.CreateMaterial)
		router.Put("/materials/{materialID}", controllers.UpdateMaterial)
		router.Delete("/materials/{materialID}", controllers.DeleteMaterial)

		router.Post("/events", controllers.CreateEvent)
		router.Put("/events/{eventID}", controllers.UpdateEvent)
		router.Delete("/events/{eventID}", controllers.DeleteEvent)
	})

	return router
}
