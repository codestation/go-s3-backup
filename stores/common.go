package stores

type Storer interface {
	Store(filepath string, key string) error
	Retrieve(s3path string) (string, error)
	RemoveOlderBackups(prefix string, keep int) error
	FindLatestBackup(prefix string) (string, error)
}
