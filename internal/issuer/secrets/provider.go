package secrets

type SecretProvider interface {
	Get(string) ([]byte, error)
}

type FileSecretProvider struct {
	secretLoader
}

func NewFileSecretProvider() (SecretProvider, error) {
	loader, err := NewFileSecretLoader()
	if err != nil {
		return nil, err
	}
	return &FileSecretProvider{
		secretLoader: loader,
	}, nil
}

func (fsp *FileSecretProvider) Get(id string) ([]byte, error) {
	// 调用 kms 查询中间证书放置的路径
	return nil, nil
}
