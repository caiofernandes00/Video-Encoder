package download_service

import (
	"encoder/application/utils"
	"encoder/domain"
	"os/exec"
)

type FragmentUseCase interface {
	Execute(sourceMp4File string, targetFragFile string) error
}

type FragmentService struct {
	Video *domain.Video
}

func NewFragmentService(v *domain.Video) *FragmentService {
	return &FragmentService{Video: v}
}

func (f *FragmentService) Execute(sourceMp4File string, targetFragFile string) error {
	println("Source mp4 file: " + sourceMp4File)
	println("Target frag file: " + targetFragFile)
	cmd := exec.Command(utils.FragmentCommand, sourceMp4File, targetFragFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	utils.PrintOutput(output)
	return nil
}
