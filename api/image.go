package api

import (
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

type listImageRequest struct {
	Key string `form:"key"`
}
func (server *Server) listImage(ctx *gin.Context) {

	var req listImageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	images, err := server.docker.listImages(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if (req.Key == "" ) {
		ctx.JSON(http.StatusOK, images)
	} else {
		var searchRes = []types.ImageSummary{}
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if (strings.Contains(tag,req.Key)){
					searchRes = append(searchRes, image)
				}
			}
		}
		ctx.JSON(http.StatusOK, searchRes)
	}
}

type pullImageRequest struct {
	ImageName string `json:"image_name" binding:"required"`
}

func (server *Server) pullImage(ctx *gin.Context) {
	var req pullImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.docker.pullImage(ctx, req.ImageName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type pullImageWithAuthRequest struct {
	ImageName string `json:"image_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (server *Server) pullImageWithAuth(ctx *gin.Context) {
	var req pullImageWithAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.docker.pullImageWithAuth(ctx, req.ImageName, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type tagImageRequest struct {
	ImageId string `json:"image_id" binding:"required"`
	Target  string `json:"target" binding:"required"`
}

func (server *Server) tagImage(ctx *gin.Context) {
	var req tagImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.docker.tagImage(ctx, req.ImageId, req.Target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type saveImageRequest struct {
	Images []string `json:"images" binding:"dive"`
}

func (server *Server) saveImage(ctx *gin.Context) {
	var req saveImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.docker.saveImages(ctx, req.Images)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func (server *Server) loadImage(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.docker.loadImages(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type pushImageRequest struct {
	ImageName string `json:"image_name" binding:"required"`
}

func (server *Server) pushImage(ctx *gin.Context) {
	var req pushImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.docker.pushImage(ctx, req.ImageName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type pushImageWithAuthRequest struct {
	ImageName string `json:"image_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (server *Server) pushImageWithAuth(ctx *gin.Context) {
	var req pushImageWithAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.docker.pushImageWithAuth(ctx, req.ImageName, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}
