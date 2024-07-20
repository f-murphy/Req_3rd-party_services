package task

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}

	if _, ok := tasks[taskID]; ok {
		c.JSON(http.StatusOK, gin.H{
			"response": tasks[taskID],
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}
}
