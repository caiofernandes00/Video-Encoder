package services

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/application/repositories"
	"encoder/application/utils"
	"encoder/domain"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

const (
	FragmentCommand = "mp4fragment"
	EncodeCommand   = "mp4dash"
	Mp4             = ".mp4"
	Frag            = ".frag"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
}

func NewVideoService(video *domain.Video, repo repositories.VideoRepository) VideoService {
	return VideoService{Video: video, VideoRepository: repo}
}

func (v *VideoService) Download(bucketname string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketname)
	obj := bkt.Object(v.Video.FilePath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has bees stored", v.Video.ID)

	closeConnections(reader, file)
	return nil
}

func (v *VideoService) Fragment() error {
	source := os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + Mp4
	target := os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + Frag

	cmd := exec.Command(FragmentCommand, source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func (v *VideoService) Encode() error {
	targetDir := os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID

	err := os.Mkdir(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	cmdArgs := encodeCommands(v, targetDir)
	cmd := exec.Command(EncodeCommand, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func (v *VideoService) Finish() error {

	err := os.Remove(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + Mp4)
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+Mp4, " with error: ", err)
		return err
	}

	err = os.Remove(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + Frag)
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+Frag, " with error: ", err)
		return err
	}

	err = os.RemoveAll(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID)
	if err != nil {
		log.Println("error removing", v.Video.ID, "folder with error: ", err)
		return err
	}

	log.Println("files have been removed: ", v.Video.ID)
	return nil
}

func removeVideoFile(filename string) error {
	err := os.Remove(os.Getenv(utils.LocalStoragePath) + "/" + filename)
	if err != nil {
		return err
	}

	return nil
}

func encodeCommands(v *VideoService, targetDir string) []string {
	var cmdArgs []string

	cmdArgs = append(cmdArgs, os.Getenv(utils.LocalStoragePath)+"/"+v.Video.ID+Frag)
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, targetDir)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	return cmdArgs
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s", string(out))
	}
}

func closeConnections(reader *storage.Reader, file *os.File) {
	defer reader.Close()
	defer file.Close()
}
