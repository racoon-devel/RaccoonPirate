package updater

type MetaInfoStorage interface {
	GetVersion() (string, error)
	SetVersion(version string) error
}
