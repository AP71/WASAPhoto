package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	// Login
	rt.router.POST("/session", rt.doLogin)

	// Profile operations
	rt.router.PUT("/profiles/:username/username", rt.wrap(rt.setMyUsername))
	rt.router.POST("/profiles/:username/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/profiles/:username/photos/:photoId", rt.wrap(rt.deletePhoto))

	// Interactions with other users
	rt.router.GET("/profiles/", rt.wrap(rt.getUsers))
	rt.router.GET("/profiles/:username/", rt.wrap(rt.getUserProfile))
	rt.router.GET("/profiles/:username/banned/:byUsername", rt.wrap(rt.banStatus))
	rt.router.PUT("/profiles/:username/banned/:byUsername", rt.wrap(rt.banUser))
	rt.router.DELETE("/profiles/:username/banned/:byUsername", rt.wrap(rt.unbanUser))
	rt.router.GET("/profiles/:username/followed/:byUsername", rt.wrap(rt.followStatus))
	rt.router.PUT("/profiles/:username/followed/:byUsername", rt.wrap(rt.followUser))
	rt.router.DELETE("/profiles/:username/followed/:byUsername", rt.wrap(rt.unfollowUser))

	// Feed operations
	rt.router.GET("/feed/", rt.wrap(rt.getMyStream))
	rt.router.GET("/feed/:photoId/", rt.wrap(rt.getPhoto))
	rt.router.GET("/feed/:photoId/likes/:username", rt.wrap(rt.likeStatus))
	rt.router.PUT("/feed/:photoId/likes/:username", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/feed/:photoId/likes/:username", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/feed/:photoId/comments/:username", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/feed/:photoId/comments/:username", rt.wrap(rt.uncommentPhoto))
	rt.router.GET("/feed/:photoId/comments/", rt.wrap(rt.getComments))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
