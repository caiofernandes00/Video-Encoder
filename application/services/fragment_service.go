package services

import (
	"encoder/application/utils"
	"encoder/domain"
	"os"
	"os/exec"
)

type FragmentUseCase interface {
	Execute() error
}

type FragmentService struct {
	Video *domain.Video
}

func NewFragmentService(v *domain.Video) *FragmentService {
	return &FragmentService{Video: v}
}

func (f *FragmentService) Execute() error {
	source := os.Getenv(utils.LocalStoragePath) + "/" + f.Video.ID + utils.Mp4Format
	target := os.Getenv(utils.LocalStoragePath) + "/" + f.Video.ID + utils.FragmentCommand

	cmd := exec.Command(utils.FragmentCommand, source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	utils.PrintOutput(output)
	return nil
}
