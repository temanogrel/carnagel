package integrity

import (
	"context"
	"os"
	"path/filepath"

	"sync"

	"time"

	"io/ioutil"

	"encoding/hex"

	"fmt"
	"git.misc.vee.bz/carnagel/minerva-bindings/src"
	"git.misc.vee.bz/carnagel/minion/pkg"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/blake2b"
)

type processRequest struct {
	info os.FileInfo
	path string
}

type fileSystemIntegrity struct {
	app *minion.Application
	mtx sync.RWMutex
	log logrus.FieldLogger

	running  bool
	scanning bool
	report   minion.IntegrityReport
}

func NewFileSystemIntegrity(app *minion.Application) minion.FilesystemIntegrity {
	return &fileSystemIntegrity{
		app: app,
		log: app.Logger.WithField("component", "fs-integrity"),
	}
}

func (fs *fileSystemIntegrity) localConsulKey() string {
	return fmt.Sprintf("minion/filesystem-integrity/%s", fs.app.Hostname)
}

func (fs *fileSystemIntegrity) IsRunning() bool {
	fs.mtx.RLock()
	defer fs.mtx.RUnlock()

	return fs.running
}

func (fs *fileSystemIntegrity) IsScanning() bool {
	fs.mtx.RLock()
	defer fs.mtx.RUnlock()

	return fs.scanning
}

func (fs *fileSystemIntegrity) GetReport() minion.IntegrityReport {
	return fs.report
}

func (fs *fileSystemIntegrity) Run(ctx context.Context) {
	fs.mtx.Lock()
	if fs.running {
		fs.log.Error("Filesystem integrity is already running")
		return
	} else {
		fs.running = true
	}
	fs.mtx.Unlock()

	defer func() {
		fs.mtx.Lock()
		fs.running = false
		fs.mtx.Unlock()
	}()

	fs.log.Info("Started filesystem integrity")

	globalValCh, globalErrCh := fs.app.Consul.PollBool(ctx, "minion/filesystem-integrity")
	// Consul.PollBool does not return an error if the key is not found
	// Instead it does not send any values over the channel
	// This means we don't need any particular error handling for this
	localValCh, localErrCh := fs.app.Consul.PollBool(ctx, fs.localConsulKey())

	var innerCtx context.Context
	var cancel context.CancelFunc

	var run bool
	var ok bool
	var err error
	var channelName string

	for {
		select {
		case <-ctx.Done():
			fs.log.Warn("Received ctx.Done() terminating the integrity system")

			// We are current scanning
			if cancel != nil {
				cancel()
			}

			return

		case err, ok = <-localErrCh:
			if !ok {
				fs.channelClosed(cancel, "Local error")

				return
			}

			fs.
				log.
				WithError(err).
				Errorf("Error occurred during the polling of %s", fs.localConsulKey())

			continue

		case err, ok = <-globalErrCh:
			if !ok {
				fs.channelClosed(cancel, "Global error")

				return
			}

			fs.log.WithError(err).Error("Error occurred during the polling of minion/filesystem-integrity")

			continue

		case run, ok = <-localValCh:
			if !ok {
				fs.channelClosed(cancel, "Local value")

				return
			}

			channelName = "local"

		case run, ok = <-globalValCh:
			if !ok {
				fs.channelClosed(cancel, "Global value")

				return
			}

			channelName = "global"
		}

		// Process the updated run value
		log := fs.
			log.
			WithFields(logrus.Fields{
				"run":      run,
				"channel":  channelName,
				"hostname": fs.app.Hostname,
			})
		log.Info("Received a run value")

		// Ignore run if we are already scanning
		if run && fs.IsScanning() {
			log.Warn("Received a start scan notification but we are already scanning")

			continue
		}

		if run {
			log.Info("Starting scan")

			// reset the report
			fs.report = minion.IntegrityReport{
				RWMutex:   &sync.RWMutex{},
				UpdatedAt: time.Now(),
			}

			innerCtx, cancel = context.WithCancel(ctx)

			go fs.scan(innerCtx)
		} else if cancel != nil {
			log.Info("Cancelling scan")
			cancel()
		}
	}
}

func (fs *fileSystemIntegrity) channelClosed(cancel context.CancelFunc, channelName string) {
	fs.log.Error(fmt.Sprintf("%s channel has been closed", channelName))

	// Check if we're currently scanning
	if cancel != nil {
		cancel()
	}
}

func (fs *fileSystemIntegrity) scan(ctx context.Context) {
	fs.mtx.Lock()
	fs.scanning = true
	fs.mtx.Unlock()

	queue := make(chan *processRequest)
	defer close(queue)

	for i := 0; i <= 5; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return

				case request, ok := <-queue:
					if !ok {
						return
					}

					fs.process(ctx, request)
				}
			}
		}()
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {

		// This will
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Probably best to just handle the error
		if err != nil {
			return err
		}

		// Nothing to do here
		if info == nil {
			return nil
		}

		queue <- &processRequest{
			info: info,
			path: path,
		}

		return nil
	}

	if err := filepath.Walk(fs.app.DataDir, walkFunc); err != nil {
		fs.log.WithError(err).Error("Error received while walking data path")
	} else {
		// Update consul key/value store
		_, err := fs.app.Consul.API().KV().Put(&api.KVPair{
			Key:   fs.localConsulKey(),
			Value: []byte("false"),
		}, nil)

		if err != nil {
			fs.log.WithError(err).Error("Error setting consul value")
		}
	}

	fs.mtx.Lock()
	fs.scanning = false
	fs.mtx.Unlock()
}

func (fs *fileSystemIntegrity) process(ctx context.Context, request *processRequest) {
	log := fs.log.WithFields(logrus.Fields{
		"filePath":       request.path,
		"filePermission": request.info.Mode(),
	})

	log.Debug("Processing file")

	fs.report.Lock()
	fs.report.Scanned += 1
	fs.report.UpdatedAt = time.Now()
	fs.report.Unlock()

	if request.info.IsDir() {
		if err := os.Chmod(request.path, 0755); err != nil {
			log.WithError(err).Error("Failed to set folder permissions")
		}

		return
	}

	if err := os.Chmod(request.path, 0644); err != nil {
		log.WithError(err).Error("Failed to set file permissions")
		return
	}

	minervaRequest := &pb.GetRequest{
		Hostname: fs.app.Hostname,
		Path:     request.path,
	}

	resp, err := fs.app.FileClient.Get(ctx, minervaRequest)
	if err != nil {
		log.WithError(err).Error("Failed get file from minerva")
		return
	}

	if resp.Status == pb.StatusCode_FileNotFound {

		fs.report.Lock()
		fs.report.Orphaned += 1
		fs.report.Unlock()

		log.Warn("Deleting file that does not exist in minerva")
		os.Remove(request.path)
		return
	}

	data, err := ioutil.ReadFile(request.path)
	if err != nil {
		log.WithError(err).Errorf("Failed to read file")
		return
	}

	checksum := blake2b.Sum256(data)
	checksumAsString := hex.EncodeToString(checksum[:])

	if resp.File.Checksum == "" {

		setChecksumRequest := &pb.UpdateRequest{
			Uuid:     resp.File.Uuid,
			Hostname: resp.File.Hostname,
			Path:     resp.File.Path,
			Size:     resp.File.Size,
			Meta:     resp.File.Meta,
			Checksum: checksumAsString,
		}

		setChecksumResponse, err := fs.app.FileClient.Update(ctx, setChecksumRequest)
		if err != nil {
			log.WithError(err).Error("Failed to call minerva to set file checksum")
			return
		}

		if setChecksumResponse.Status != pb.StatusCode_Ok {
			log.
				WithField("statusCode", setChecksumResponse.Status).
				Error("Unexpected response code from minerva")

			return
		}
	} else if resp.File.Checksum != checksumAsString {
		log.WithFields(logrus.Fields{
			"uuid":        resp.File.Uuid,
			"externalId":  resp.File.ExternalId,
			"oldChecksum": resp.File.Checksum,
			"newChecksum": checksumAsString,
		}).Error("Checksum does not match")

		return
	}

	log.Debug("Successfully processed")
}
