package api

import (
	"net/http"
	"strconv"

	"bilibili-analyzer/backend/comment"
	"bilibili-analyzer/backend/task"

	"github.com/gin-gonic/gin"
)

func HandleListHistoryComments(c *gin.Context) {
	history, err := findHistoryByIDOrTaskID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History record not found"})
		return
	}
	rows, err := comment.LoadByHistory(history.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load comments"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func HandleReanalyzeHistory(c *gin.Context) {
	idStr := c.Param("id")
	history, err := findHistoryByIDOrTaskID(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History record not found"})
		return
	}
	reportID, err := task.ReanalyzeFromStore(c.Request.Context(), history.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"report_id": reportID,
		"history_id": history.ID,
		"id":        strconv.FormatUint(uint64(history.ID), 10),
	})
}
