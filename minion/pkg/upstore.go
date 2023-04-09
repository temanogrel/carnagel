package minion

import "os"

type UpstoreClient interface {
	Upload(name string, file *os.File) (string, error)
}
