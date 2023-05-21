package usecase

type UseCase interface {
	Execute(...interface{}) (interface{}, error)
}
