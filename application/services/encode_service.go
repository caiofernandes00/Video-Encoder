package services

import (
	"encoder/application/utils"
	"encoder/domain"
	"os"
	"os/exec"
)

type EncodeUseCase interface {
	Execute() error
}

type EncodeService struct {
	Video *domain.Video
}

func NewEncodeService(v *domain.Video) *EncodeService {
	return &EncodeService{Video: v}
}

func (e *EncodeService) Execute() error {
	targetDir := os.Getenv(utils.LocalStoragePath) + "/" + e.Video.ID

	err := os.Mkdir(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	cmdArgs := encodeCommands(e, targetDir)
	cmd := exec.Command(utils.EncodeCommand, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	utils.PrintOutput(output)
	return nil
}

func encodeCommands(v *EncodeService, targetDir string) []string {
	var cmdArgs []string

	cmdArgs = append(cmdArgs, os.Getenv(utils.LocalStoragePath)+"/"+v.Video.ID+utils.FragFile)
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, targetDir)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	return cmdArgs
}
