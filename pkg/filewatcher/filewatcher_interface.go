package filewatcher

type FileEventCallbackFn func(path string)

type FileWatcher interface {
	Start(path string, created, modified, deleted FileEventCallbackFn) error
	Stop()
}
