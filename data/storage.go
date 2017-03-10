package data

type Storage interface {
	AddUploadedFragment(fragment UploadedFragment) error
}
