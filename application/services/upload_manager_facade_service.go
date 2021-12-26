package services

type UploadManagerUseCase interface {
}

type UploadManagerService struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *UploadManagerService {
	return &UploadManagerService{}
}

func (u *UploadManagerService) Execute() error {

}
