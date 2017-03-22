package data

type StubStorage struct{}

func NewStubStorage() Storage {
	return &StubStorage{}
}

func (s *StubStorage) AddUploadedFragment(fragment UploadedFragment) error {
	return nil
}
