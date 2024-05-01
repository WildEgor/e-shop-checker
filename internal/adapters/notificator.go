package adapters

type INotificator interface {
	Send(msg string) error
}

type NotificatorAdapter struct {
}

func NewNotificatorAdapter() *NotificatorAdapter {
	return &NotificatorAdapter{}
}

func (na *NotificatorAdapter) Send(msg string) error {

	return nil
}
