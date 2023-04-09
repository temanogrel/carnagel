package tasks

//
//import (
//	"os"
//
//	"io/ioutil"
//
//	"context"
//
//	"git.misc.vee.bz/carnagel/encoder/pkg"
//	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
//	"github.com/satori/go.uuid"
//	"github.com/sirupsen/logrus"
//)
//
//type taskService struct {
//	app *encoder.Application
//	log logrus.FieldLogger
//}
//
//func NewTaskService(app *encoder.Application) encoder.TaskService {
//	return &taskService{
//		app: app,
//		log: app.Logger.WithField("service", "TaskService"),
//	}
//}
//
//func (service *taskService) EncodeToMp4(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "EncodeToMp4", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	flvFilePath := ctx.Value("flvFilePath").(string)
//	if flvFilePath == "" {
//		log.WithField("value", "flvFilePath").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	mp4EncodedFile := ctx.Value("mp4EncodedFile").(string)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "mp4EncodedFile").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	meta := ctx.Value("meta").(*encoder.VideoMeta)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "meta").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	// save the old uuid so can delete it later
//	oldUuid := recording.VideoMp4Uuid
//
//	recording.State = "encoding"
//	if err := service.app.AphroditeClient.UpdateRecording(recording); err != nil {
//		log.WithError(err).Error("Failed to update recording state")
//	}
//
//	var err error
//
//	if meta.Size >= service.app.H265Threshold {
//
//		log.Debug("Encoding file to h265")
//		err = service.app.EncoderService.EncodeToH265(ctx, flvFilePath, mp4EncodedFile, meta)
//	} else {
//
//		log.Debug("Encoding file to h264")
//		err = service.app.EncoderService.EncodeToH264(ctx, flvFilePath, mp4EncodedFile, meta)
//	}
//
//	if err != nil {
//		log.WithError(err).Error("Failed to encoded video, recording will be deleted")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	if err := service.setRecordingMeta(ctx, recording, mp4EncodedFile); err != nil {
//		log.WithError(err).Error("Failed to set recording meta")
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	pipeline.Defer(func() {
//		log.Debug("Deleting flv recording from minerva")
//		if err := service.app.MinervaClient.RequestDeletion(oldUuid); err != nil {
//			log.WithError(err).Error("Failed to request the deletion of the old file")
//			return
//		}
//	})
//
//	return nil
//}
//
//func (service *taskService) EncodeHlsToMp4(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "EncodeHlsToMp4", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	panic("not implemented")
//}
//
//func (service *taskService) EncodeMp4ToHls(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "EncodeMp4", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	flvFilePath := ctx.Value("flvFilePath").(string)
//	if flvFilePath == "" {
//		log.WithField("value", "flvFilePath").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	mp4EncodedFile := ctx.Value("mp4EncodedFile").(string)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "mp4EncodedFile").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	hlsVideoFile := encoder.ReplaceFileExtension(flvFilePath, "ts")
//	hlsManifestFile := encoder.ReplaceFileExtension(flvFilePath, "m3u8")
//
//	pipeline.Defer(func() {
//		os.Remove(hlsVideoFile)
//		os.Remove(hlsManifestFile)
//	})
//
//	if err := service.app.EncoderService.EncodeMp42Hls(ctx, mp4EncodedFile, hlsVideoFile, hlsManifestFile); err != nil {
//		log.WithError(err).Error("Failed to convert mp4 to hls")
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	ctx = context.WithValue(ctx, "hlsVideoFile", hlsVideoFile)
//	ctx = context.WithValue(ctx, "hlsManifestFile", hlsManifestFile)
//
//	return nil
//}
//
//func (service *taskService) DecodeMp4(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "DecodeMp4", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	panic("Not implemented")
//}
//
//func (service *taskService) DownloadFromUpstore(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "DownloadFromUpstore", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	panic("Not implemented")
//}
//
//
//
//func (service *taskService) JpegInfo(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "JpegInfo", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	panic("Not implemented")
//}
//
//func (service *taskService) UploadToMinerva(ctx context.Context, pipeline *encoder.SerialPipeline, recording *ecosystem.Recording) error {
//	log := service.log.WithFields(logrus.Fields{"stage": "UploadToMinerva", "pipelineId": pipeline.Id})
//	log.Debug("Running pipeline stage")
//
//	mp4EncodedFile := ctx.Value("mp4EncodedFile").(string)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "mp4EncodedFile").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	hlsVideoFile := ctx.Value("hlsVideoFile").(string)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "hlsVideoFile").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	hlsManifestFile := ctx.Value("hlsManifestFile").(string)
//	if mp4EncodedFile == "" {
//		log.WithField("value", "hlsManifestFile").Warn("Missing context value")
//
//		return encoder.MissingContextValueErr
//	}
//
//	mp4File, err := os.Open(mp4EncodedFile)
//	if err != nil {
//		log.WithError(err).Error("Failed to open encoded file for uploading")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	hlsFile, err := os.Open(hlsVideoFile)
//	if err != nil {
//		log.WithError(err).Error("Failed top open hls encoded file")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	hlsManifest, err := ioutil.ReadFile(hlsManifestFile)
//	if err != nil {
//		log.WithError(err).Error("Failed to open hls manifest file")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	hlsUuid, err := service.app.MinervaClient.Upload(hlsFile, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeRecordingHls, ecosystem.FileMetadata{})
//	if err != nil {
//		log.WithError(err).Error("Failed to upload HLS file")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	mp4Uuid, err := service.app.MinervaClient.Upload(mp4File, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeRecordingMp4, ecosystem.FileMetadata{})
//	if err != nil {
//		log.WithError(err).Error("Failed to upload the encoded file")
//
//		service.remove(ctx, recording)
//
//		return err
//	}
//
//	// Update state
//	recording.State = "encoded"
//	recording.VideoMp4Uuid = mp4Uuid
//	recording.VideoHlsUuid = hlsUuid
//	recording.VideoManifest = string(hlsManifest)
//
//	log.Debug("Updating recording state in aphrodite")
//	if err := service.app.AphroditeClient.UpdateRecording(recording); err != nil {
//		log.WithError(err).Error("Failed to update aphrodite recording entity")
//
//		return err
//	}
//
//	log.Debug("Requesting file upload")
//	if err := service.app.MinervaClient.RequestUpload(mp4Uuid, recording.GetPublishedFilename("mp4")); err != nil {
//		log.WithError(err).Error("Failed to request upload")
//
//		return err
//	}
//
//	return nil
//}
//
//func (service *taskService) generateWordpressCollage(ctx context.Context, recording *ecosystem.Recording, videoFile string) error {
//	target := encoder.ExtendFileName(videoFile, "-wp-collage", "jpg")
//
//	defer os.Remove(target)
//
//	if err := service.app.EncoderService.GenerateWordpressCollage(ctx, videoFile, target, recording.GetPublishedFilename("mp4")); err != nil {
//		return err
//	}
//
//	file, err := os.Open(target)
//	if err != nil {
//		return err
//	}
//
//	defer file.Close()
//
//	id, err := service.app.MinervaClient.Upload(file, ecosystem.ExternalId(recording.Id), ecosystem.FileTypeWordpressCollage, nil)
//	if err != nil {
//		return err
//	}
//
//	recording.WordpressCollageUuid = id
//
//	return nil
//}
//
//func (service *taskService) remove(ctx context.Context, recording *ecosystem.Recording) error {
//	logger := service.app.Logger.WithFields(logrus.Fields{
//		"component":   "encoder",
//		"recordingId": recording.Id,
//	})
//
//	logger.Warn("File does not exist, it will be deleted")
//
//	if recording.InfinityCollageUuid != uuid.Nil {
//		if err := service.app.MinervaClient.RequestDeletion(recording.InfinityCollageUuid); err != nil {
//			logger.WithError(err).Fatal("Failed to delete infinity collage")
//			return err
//		}
//	}
//
//	if recording.WordpressCollageUuid != uuid.Nil {
//		if err := service.app.MinervaClient.RequestDeletion(recording.WordpressCollageUuid); err != nil {
//			logger.WithError(err).Fatal("Failed to delete wordpress collage")
//			return err
//		}
//	}
//
//	for _, id := range recording.Sprites {
//		if err := service.app.MinervaClient.RequestDeletion(id); err != nil {
//			logger.WithError(err).Fatal("Failed to delete sprite")
//			return err
//		}
//	}
//
//	for _, id := range recording.Images {
//		if err := service.app.MinervaClient.RequestDeletion(id); err != nil {
//			logger.WithError(err).Fatal("Failed to delete image")
//			return err
//		}
//	}
//
//	if recording.VideoMp4Uuid != uuid.Nil {
//		if err := service.app.MinervaClient.RequestDeletion(recording.VideoMp4Uuid); err != nil {
//			logger.WithError(err).Error("Failed to schedule deletion for mp4/flv file")
//			return err
//		}
//	}
//
//	if recording.VideoHlsUuid != uuid.Nil {
//		if err := service.app.MinervaClient.RequestDeletion(recording.VideoHlsUuid); err != nil {
//			logger.WithError(err).Error("Failed to schedule deletion for hls file")
//			return err
//		}
//	}
//
//	if err := service.app.AphroditeClient.DeleteRecording(recording.Id); err != nil {
//		logger.WithError(err).Error("Failed to delete from aphrodite")
//		return err
//	}
//
//	logger.Debug("File deleted successfully")
//
//	return nil
//}
