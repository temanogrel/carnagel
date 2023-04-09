package file

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"io/ioutil"

	"encoding/hex"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/blake2b"
)

type fileService struct {
	app *minion.Application
}

func NewFileService(application *minion.Application) minion.FileService {
	return &fileService{application}
}

func (service *fileService) Transfer(ctx context.Context, id uuid.UUID, targetHost string) error {

	data, err := service.app.MinervaClient.GetDataByUuid(id)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve information about file")
	}

	// Check the external host name matches
	if data.Hostname != service.app.Hostname {
		return errors.Errorf("File hostname %s does not match with server hostname %s", data.Hostname, service.app.Hostname)
	}

	fp, err := os.Open(data.Path)
	if err != nil {
		return errors.Wrap(err, "Failed to open the file")
	}

	defer fp.Close()

	err = service.app.MinervaClient.Transfer(ctx, id, ecosystem.ExternalId(data.ExternalId), fp, targetHost)
	if err != nil {
		return errors.Wrap(err, "Failed to transfer file to remote server")
	}

	if err := os.Remove(data.Path); err != nil {
		return errors.Wrap(err, "Failed to delete the file locally")
	}

	return nil
}

func (service *fileService) HandleTransfer(ctx context.Context, id uuid.UUID, file io.Reader, extension string) error {

	data, err := service.app.MinervaClient.GetDataByUuid(id)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve information about file")
	}

	dir, err := service.generateRandomDir()
	if err != nil {
		return errors.Wrap(err, "Failed to generate random directory")
	}

	name := fmt.Sprintf("%s%s", uuid.NewV4().String(), extension)
	path := filepath.Join(dir, name)

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open the target file")
	}

	defer fp.Close()

	if _, err := io.Copy(fp, file); err != nil {
		return errors.Wrap(err, "Failed to write to the target file")
	}

	request := &pb.UpdateRequest{
		Path:     path,
		Uuid:     data.Uuid,
		Meta:     data.Meta,
		Size:     data.Size,
		Hostname: service.app.Hostname,
	}

	response, err := service.app.FileClient.Update(ctx, request)
	if err != nil {
		if e := os.Remove(path); err != nil {
			return errors.Wrap(e, "Failed to delete transferred file after request to relocate file failed")
		}

		return errors.Wrap(err, "Failed to relocate file in minerva")
	}

	if response.Status != pb.StatusCode_Ok {
		return errors.Errorf("Unexpected response code %d received", response.Status)
	}

	return nil
}

func (service *fileService) HandleUpload(
	file io.Reader,
	originalFileName string,
	externalId uint64,
	fileType pb.FileType) (uuid.UUID, error) {

	dir, err := service.generateRandomDir()
	if err != nil {
		return uuid.Nil, err
	}

	name := fmt.Sprintf("%s%s", uuid.NewV4().String(), filepath.Ext(originalFileName))
	path := filepath.Join(dir, name)

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "Failed to open the target file")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "Failed to read the entire file")
	}

	checksum := blake2b.Sum256(data)

	defer fp.Close()

	written, err := fp.Write(data)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "Failed to write to the target file")
	}

	created := &pb.CreateRequest{
		Hostname:   service.app.Hostname,
		Checksum:   hex.EncodeToString(checksum[:]),
		ExternalId: externalId,
		Type:       fileType,
		Path:       path,
		Size:       uint64(written),
	}

	resp, err := service.app.FileClient.Create(context.TODO(), created)
	if err != nil {

		// Delete the file so we don't end up with a bunch of shit on the servers
		os.Remove(path)

		return uuid.Nil, errors.Wrap(err, "Failed to create the file in minerva")
	}

	if resp.Status != pb.StatusCode_Ok {

		// Delete the file so we don't end up with a bunch of shit on the servers
		os.Remove(path)

		return uuid.Nil, errors.New(fmt.Sprintf("Unexpected response code %d", resp.Status))
	}

	return uuid.FromStringOrNil(resp.Uuid), nil
}

func (service *fileService) generateRandomDir() (string, error) {

	path := service.app.DataDir

	for i := 0; i <= 3; i++ {
		path = filepath.Join(path, randSeq(2))
	}

	err := os.MkdirAll(path, 0755)
	return path, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
