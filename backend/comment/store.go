package comment

import (
	"fmt"
	"strconv"
	"time"

	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

type StoredComment struct {
	HistoryID uint   `json:"history_id"`
	VideoID   string `json:"video_id"`
	CommentID string `json:"comment_id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Likes     int    `json:"likes"`
}

func PersistForHistory(historyID uint, comments []bilibili.Comment, videoByKey map[string]string) error {
	if database.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	if historyID == 0 || len(comments) == 0 {
		return nil
	}
	rows := make([]models.RawComment, 0, len(comments))
	for _, c := range comments {
		id := strconv.FormatInt(c.RPID, 10)
		if id == "0" {
			id = fmt.Sprintf("local-%d-%d", historyID, len(rows))
		}
		videoID := ""
		if videoByKey != nil {
			videoID = videoByKey[id]
		}
		rows = append(rows, models.RawComment{
			HistoryID:   historyID,
			VideoID:     videoID,
			CommentID:   fmt.Sprintf("%d:%s", historyID, id),
			Content:     c.Content.Message,
			Author:      c.Member.Uname,
			Likes:       c.Like,
			ReplyCount:  c.Count,
			PublishTime: time.Unix(c.Ctime, 0),
		})
	}
	return database.DB.CreateInBatches(&rows, 100).Error
}

func LoadByHistory(historyID uint) ([]StoredComment, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var rows []models.RawComment
	if err := database.DB.Where("history_id = ?", historyID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]StoredComment, 0, len(rows))
	for _, r := range rows {
		out = append(out, StoredComment{
			HistoryID: r.HistoryID,
			VideoID:   r.VideoID,
			CommentID: r.CommentID,
			Content:   r.Content,
			Author:    r.Author,
			Likes:     r.Likes,
		})
	}
	return out, nil
}
