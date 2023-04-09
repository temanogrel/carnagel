package file

import (
	"encoding/json"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

//noinspection GoNameStartsWithPackageName
type fileService struct {
	app *minerva.Application

	publishChannel *amqp.Channel
}

func NewFileService(application *minerva.Application) minerva.FileService {
	service := &fileService{app: application}
	service.setupRabbitmqExchange()

	return service
}

func (service *fileService) setupRabbitmqExchange() {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		service.app.Logger.WithError(err).Fatal("Failed to create amqp channel")
	}

	service.publishChannel = ch
}

func (service *fileService) CleanUp() {
	// Clean up
	service.publishChannel.Close()
}

func (service *fileService) Create(data *minerva.CreateFile) (*minerva.File, error) {

	file := &minerva.File{
		Uuid:       uuid.NewV4(),
		ExternalId: data.ExternalId,

		Type:     ecosystem.FileType(data.Type),
		Hostname: minerva.Hostname(data.Hostname),
		Checksum: data.Checksum,
		Path:     data.Path,

		Size: data.Size,
		Meta: data.Meta,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if file.Type == ecosystem.FileTypeUnknown {
		return nil, errors.New("Invalid file type specified, cannot be of type All")
	}

	if err := service.app.FileRepository.Create(file); err != nil {
		return nil, err
	}

	return file, nil
}

func (service *fileService) Update(data *minerva.UpdateFile) (*minerva.File, error) {

	file, err := service.app.FileRepository.GetByUuid(data.Uuid)
	if err != nil {
		return nil, err
	}

	file.Size = data.Size
	file.Meta = data.Meta
	file.Path = data.Path
	file.Checksum = data.Checksum
	file.Hostname = minerva.Hostname(data.Hostname)
	file.UpdatedAt = time.Now()

	if err = service.app.FileRepository.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

func (service *fileService) ScheduleDeletion(id uuid.UUID) error {

	file, err := service.app.FileRepository.GetByUuid(id)
	if err != nil {
		return err
	}

	if file.PendingDeletion == true {
		return minerva.FileAlreadyPendingDeletionErr
	}

	file.PendingDeletion = true
	file.PendingUpload = false

	service.app.ServerCollection.ServersMtx.RLock()
	server, ok := service.app.ServerCollection.Servers[file.Hostname]
	service.app.ServerCollection.ServersMtx.RUnlock()

	if ok {
		server.DeletionRequests <- id
	} else {
		service.app.Logger.
			WithField("uuid", id.String()).
			Warn("Scheduled file for deletion but not available in server collection")
	}

	if err = service.app.FileRepository.Update(file); err != nil {
		return err
	}

	return nil
}

func (service *fileService) ScheduleDeletionByExternalIdAndType(id ecosystem.ExternalId, fileType ecosystem.FileType) error {
	files, err := service.app.FileRepository.GetByExternalId(id)
	if err != nil {
		return err
	}

	for _, file := range files {
		if fileType == ecosystem.FileTypeUnknown || fileType == file.Type {
			service.ScheduleDeletion(file.Uuid)
		}
	}

	return nil
}

func (service *fileService) ScheduleUpload(id uuid.UUID, name string) error {

	file, err := service.app.FileRepository.GetByUuid(id)
	if err != nil {
		return err
	}

	if file.PendingUpload == true {
		return minerva.FileAlreadyPendingUploadErr
	}

	file.PendingUpload = true
	file.OriginalFilename = name

	// we must update the database before pushing it to the upload requests channel
	if err = service.app.FileRepository.Update(file); err != nil {
		return err
	}

	service.app.ServerCollection.ServersMtx.RLock()
	server, ok := service.app.ServerCollection.Servers[file.Hostname]
	service.app.ServerCollection.ServersMtx.RUnlock()

	if ok {
		server.UploadRequests <- id
	} else {
		service.app.Logger.
			WithField("uuid", id.String()).
			Warn("Scheduled file for upload but not available in server collection")
	}

	return nil
}

func (service *fileService) SetUpstoreHash(id uuid.UUID, hash string) error {
	file, err := service.app.FileRepository.GetByUuid(id)
	if err != nil {
		return err
	}

	file.PendingUpload = false
	file.UpstoreHash = hash

	if err = service.app.FileRepository.Update(file); err != nil {
		return err
	}

	return service.notifyAboutUpload(file)
}

func (service *fileService) Delete(id uuid.UUID) error {

	file, err := service.app.FileRepository.GetByUuid(id)
	if err != nil {
		return err
	}

	return service.app.FileRepository.Delete(file)
}

func (service *fileService) notifyAboutUpload(file *minerva.File) error {
	body, err := json.Marshal(&ecosystem.FileUploadedPayload{
		ExternalId: uint64(file.ExternalId),
		Timestamp:  time.Now(),
	})

	if err != nil {
		return errors.Wrap(err, "Failed to encode fileUploaded jsn")
	}

	payload := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	if err := service.publishChannel.Publish("", ecosystem.AmqpQueueUploaded, false, false, payload); err != nil {
		return errors.Wrapf(err, "Failed to publish to the queue %s", ecosystem.AmqpQueueUploaded)
	}

	service.app.Logger.WithField("recordingId", file.ExternalId).Debug("Dispatched file for publishing")

	return nil
}
