package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/profile/:username/", rt.setMyUsername)
	rt.router.POST("/profile/:username/photos/", rt.uploadPhoto)
	rt.router.DELETE("/profile/:username/photos/:photoId", rt.deletePhoto)
	rt.router.GET("/profile/", rt.searchUsers)
	rt.router.GET("/profile/:username/", rt.getUserProfile)
	rt.router.PUT("/profile/:username/banned/:byUsername", rt.banUser)
	rt.router.DELETE("/profile/:username/banned/:byUsername", rt.unbanUser)
	rt.router.PUT("/profile/:username/followed/:byUsername", rt.followUser)
	rt.router.DELETE("/profile/:username/followed/:byUsername", rt.unfollowUser)
	rt.router.GET("/feed/", rt.getMyStream)
	rt.router.GET("/feed/:photoId/", rt.getPhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
