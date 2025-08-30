package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
)

type DownloadRequestHelper struct {
	repos *Repos
}

func NewDownloadRequestHelper(repos *Repos) *DownloadRequestHelper {
	return &DownloadRequestHelper{repos: repos}
}

func (h *DownloadRequestHelper) CheckAndCreateDownload(ctx context.Context, userID int64, md5, title string) (bool, error) {
	canDownload, err := h.repos.DownloadRequest.CheckAndCreateDownload(ctx, userID, md5, title)
	if err != nil {
		return false, fmt.Errorf("failed to check and create download: %w", err)
	}

	return canDownload, nil
}

func (h *DownloadRequestHelper) GetDownloadStatus(ctx context.Context, userID int64) (map[string]interface{}, error) {
	canDownload, err := h.repos.DownloadRequest.CanDownload(ctx, userID)
	if err != nil {
		return nil, err
	}

	remaining, err := h.repos.DownloadRequest.GetRemainingDownloads(ctx, userID)
	if err != nil {
		return nil, err
	}

	count, err := h.repos.DownloadRequest.GetDailyDownloadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	nextReset := utils.GetNextUTCMidnight()
	timeUntilReset := utils.GetTimeUntilNextUTCMidnight()

	return map[string]interface{}{
		"can_download":        canDownload,
		"downloads_used":      count,
		"downloads_remaining": remaining,
		"daily_limit":         10,
		"next_reset":          nextReset.Format(time.RFC3339),
		"time_until_reset":    timeUntilReset.String(),
	}, nil
}

func (h *DownloadRequestHelper) GetUserDownloadHistory(ctx context.Context, userID int64, limit, offset int) ([]model.DownloadRequest, error) {
	return h.repos.DownloadRequest.GetUserDownloadHistory(ctx, userID, limit, offset)
}

func (h *DownloadRequestHelper) CanDownloadBook(ctx context.Context, userID int64, bookHash string) (bool, error) {
	// Check if user has a valid download job for this book
	job, err := h.repos.DownloadJob.GetJobByUserAndBook(ctx, userID, bookHash)
	if err != nil {
		return false, err
	}

	// Check if job is completed
	if job.Status != model.DownloadStatusCompleted {
		return false, nil
	}

	// Check daily download limit
	canDownload, err := h.repos.DownloadRequest.CheckAndCreateDownload(ctx, userID, bookHash, "")
	if err != nil {
		return false, err
	}

	return canDownload, nil
}

func (h *DownloadRequestHelper) GetBookDownloadStatus(ctx context.Context, userID int64, bookHash string) (*model.DownloadJob, error) {
	// Get the download job for this user and book
	job, err := h.repos.DownloadJob.GetJobByUserAndBook(ctx, userID, bookHash)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (h *DownloadRequestHelper) GetDownloadsByMD5(ctx context.Context, md5 string) ([]model.DownloadRequest, error) {
	return h.repos.DownloadRequest.GetDownloadsByMD5(ctx, md5)
}

func (h *DownloadRequestHelper) GetDailyDownloadStats(ctx context.Context, userID int64, days int) (map[string]int, error) {
	return h.repos.DownloadRequest.GetDailyDownloadStats(ctx, userID, days)
}

func (h *DownloadRequestHelper) CleanupOldRequests(ctx context.Context, daysOld int) error {
	return h.repos.DownloadRequest.CleanupOldRequests(ctx, daysOld)
}
