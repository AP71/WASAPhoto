package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/profile/:username/username", rt.wrap(rt.setMyUsername))
	rt.router.POST("/profile/:username/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/profile/:username/photos/:photoId", rt.wrap(rt.deletePhoto))
	rt.router.GET("/profile/", rt.wrap(rt.getUsers))
	rt.router.GET("/profile/:username/", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/profile/:username/banned/:byUsername", rt.wrap(rt.banUser))
	rt.router.DELETE("/profile/:username/banned/:byUsername", rt.wrap(rt.unbanUser))
	rt.router.PUT("/profile/:username/followed/:byUsername", rt.wrap(rt.followUser))
	rt.router.DELETE("/profile/:username/followed/:byUsername", rt.wrap(rt.unfollowUser))
	rt.router.GET("/feed/", rt.wrap(rt.getMyStream))
	rt.router.GET("/feed/:photoId/", rt.wrap(rt.getPhoto))
	rt.router.PUT("/feed/:photoId/likes/:username", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/feed/:photoId/likes/:username", rt.wrap(rt.unlikePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
