package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"created_at" valid:"-"`
	job        []*Job    `valid:"-" gorm:"ForeignKey:VideoID"`
}

func NewVideo(resourceId string, filePath string) (*Video, error) {
	video := Video{
		ResourceID: resourceId,
		FilePath:   filePath,
	}
	video.prepare()

	err := video.Validate()
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (video *Video) prepare() {
	video.ID = uuid.NewV4().String()
	video.CreatedAt = time.Now()
}

func (video *Video) Validate() error {
	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil

}
