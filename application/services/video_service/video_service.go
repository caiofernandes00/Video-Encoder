package video_service

import (
	"encoder/application/repositories"
	"encoder/domain"
)

type VideoUseCase interface {
	InsertVideo() error
}

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService(video *domain.Video, videoRepository repositories.VideoRepository) *VideoService {
	return &VideoService{Video: video, VideoRepository: videoRepository}
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)

	if err != nil {
		return err
	}

	return nil
}
