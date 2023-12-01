package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type Docker struct {
	cli *client.Client
}

func NewDocker(c *client.Client) *Docker {
	return &Docker{
		cli: c,
	}
}

func (d *Docker) listImages(ctx *gin.Context) (images []types.ImageSummary, err error) {
	images, err = d.cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (d *Docker) pullImage(ctx *gin.Context, imageName string) (err error) {
	out, err := d.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	response, err := io.ReadAll(out)
	if err != nil {
		return err
	}
	err = dockerResponseCheck(ctx, response)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

func (d *Docker) pullImageWithAuth(ctx *gin.Context, imageName, username, passwod string) (err error) {
	authConfig := registry.AuthConfig{
		Username: username,
		Password: passwod,
	}
	var encodedJSON []byte
	encodedJSON, err = json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, err := d.cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}
	response, err := io.ReadAll(out)
	if err != nil {
		return err
	}
	err = dockerResponseCheck(ctx, response)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

func (d *Docker) tagImage(ctx *gin.Context, source, target string) (err error) {
	err = d.cli.ImageTag(ctx, source, target)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) saveImages(ctx *gin.Context, images []string) (err error) {
	saveResponse, err := d.cli.ImageSave(ctx, images)
	if err != nil {
		return err
	}
	defer saveResponse.Close()

	// TODO: can OOM?
	response, err := io.ReadAll(saveResponse)
	if err != nil {
		return err
	}

	fileName := strconv.FormatInt(time.Now().Unix(), 10)+".tar"
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(response)))
	ctx.Writer.Write(response)

	return nil
}

func (d *Docker) loadImages(ctx *gin.Context, file multipart.File) (err error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	input := bytes.NewReader(content)
	imageLoadResponse, err := d.cli.ImageLoad(ctx, input, true)
	if err != nil {
		return err
	}
	response, err := io.ReadAll(imageLoadResponse.Body)
	if err != nil {
		return err
	}
	err = dockerResponseCheck(ctx, response)
	if err != nil {
		return err
	}

	defer imageLoadResponse.Body.Close()
	return nil
}

// not auth cant access resource
func (d *Docker) pushImage(ctx *gin.Context, imageName string) (err error) {

	authConfig := registry.AuthConfig{
		Username: "docker",
		Password: "",
	}
	var encodedJSON []byte
	encodedJSON, err = json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	resp, err := d.cli.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}

	response, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	err = dockerResponseCheck(ctx, response)
	if err != nil {
		return err
	}

	defer resp.Close()
	return nil
}

func (d *Docker) pushImageWithAuth(ctx *gin.Context, imageName, username, passwod string) (err error) {
	authConfig := registry.AuthConfig{
		Username: username,
		Password: passwod,
	}
	var encodedJSON []byte
	encodedJSON, err = json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	resp, err := d.cli.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})

	response, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	err = dockerResponseCheck(ctx, response)
	if err != nil {
		return err
	}

	defer resp.Close()
	return nil
}

func dockerResponseCheck(ctx *gin.Context, response []byte) (err error) {
	fmt.Println("---[docker log start]---")
	fmt.Println(string(response))
	fmt.Println("---[docker log  end ]---")
	re, _ := regexp.Compile(`"error":(.*)}`)
	result := re.FindAllString(string(response), -1)
	if len(result) > 0 {
		return errors.New(strings.TrimSuffix(result[0], "}"))
	} else {
		return nil
	}
}

func (d *Docker) removeImgae(ctx *gin.Context, imageId string) (imageDeletes []types.ImageDeleteResponseItem,err error) {
	imageDeletes, err =  d.cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return imageDeletes,err
}

func (d *Docker) removeImgaes(ctx *gin.Context, imageIds []string) (imageDeletes []types.ImageDeleteResponseItem,err error) {
	for _, imageId := range imageIds {
		res, err := d.cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			return imageDeletes,err
		}
		imageDeletes = append(imageDeletes, res...)
	}
	return imageDeletes,err
}