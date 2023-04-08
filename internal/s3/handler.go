package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"wildfire-backend/internal/handlers"
	s32 "wildfire-backend/pkg/s3"
)

type handler struct {
	client s32.Client
}

func NewS3Handler(client s32.Client) handlers.Handler {
	return &handler{
		client: client,
	}
}

func (h handler) Register(router *gin.Engine) {
	router.GET("/tiles/", h.GetHotSpots)
}

func (h *handler) GetHotSpots(ctx *gin.Context) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(h.client.Bucket),
		Key:    aws.String(ctx.Request.URL.Path),
	}
	object, err := h.client.Client.GetObject(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	if object.Body != nil {
		defer object.Body.Close()
	}
	if object.ContentLength == nil {
		return
	}
	ctx.Header("Content-Type", *object.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", *object.ContentLength))
	ctx.Status(http.StatusOK)
	io.Copy(ctx.Writer, object.Body)
}
