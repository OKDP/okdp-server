/*
 *    Copyright 2025 okdp.io
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package controllers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	_api "github.com/okdp/okdp-server/api/openapi/v3/_api/pods"
	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/services"
)

type IPodController struct {
	podService *services.PodService
}

func PodController() *IPodController {
	return &IPodController{
		podService: services.NewPodService(),
	}
}

func (r IPodController) GetLogs(c *gin.Context, clusterID, namespace, pod, container string, params _api.GetLogsParams) {
	accept := c.GetHeader("Accept")
	isDownload := params.Download != nil && *params.Download
	isSSE := strings.Contains(accept, "text/event-stream")
	var tailLines *int64
	if !isDownload {
		def := int64(100)
		tailLines = &def
	}
	if params.TailLines != nil {
		v := int64(*params.TailLines)
		tailLines = &v
	}
	stream, err := r.podService.StreamLogs(clusterID, namespace, pod, container, tailLines, isSSE)
	if err != nil {
		c.AbortWithStatusJSON(err.Status, err)
		return
	}
	defer stream.Close()

	switch {
	case isSSE:
		r.streamPodLogs(c, stream)
	case isDownload:
		r.downloadPodLogs(c, pod, container, stream)
	default:
		r.readPodLogs(c, stream)
	}
}

func (r IPodController) streamPodLogs(c *gin.Context, stream io.ReadCloser) {
	log.Debug("Streaming pod logs using Server-Sent Events (SSE) ...")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "streaming not supported",
		})
		return
	}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(c.Writer, "data: %s\n\n", line)
		flusher.Flush()
	}
	if err := scanner.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to stream logs",
			"details": err.Error(),
		})
	}
}

func (r IPodController) downloadPodLogs(c *gin.Context, pod, container string, stream io.ReadCloser) {
	log.Debug("Downloading pod logs as a plain text file ...")
	filename := fmt.Sprintf("%s-%s-logs.log", pod, container)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Status(http.StatusOK)
	if _, copyErr := io.Copy(c.Writer, stream); copyErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to write logs",
			"details": copyErr.Error(),
		})
	}
}

func (r IPodController) readPodLogs(c *gin.Context, stream io.ReadCloser) {
	log.Debug("Returning logs as JSON array ...")
	var logs []string
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to read logs",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, logs)
}
