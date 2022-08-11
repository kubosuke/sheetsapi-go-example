package usecase

type NotifyUsecase interface {
	Notify(id int, text string) (bool, error)
}
