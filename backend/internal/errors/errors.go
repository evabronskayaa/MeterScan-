package errors

type SimpleError string

func (e SimpleError) Error() string {
	return string(e)
}

const (
	ErrIncorrectRequest SimpleError = "Некорректный запрос"
	ErrNotFoundUser     SimpleError = "Пользователь не найден"
	ErrDuplicateEmail   SimpleError = "Данная почта уже используется"
	ErrCreatePassword   SimpleError = "Произошла ошибка при хэшировании пароля"
	ErrSaveUser         SimpleError = "Произошла ошибка при сохранении пользователя"
	ErrAlreadyVerified  SimpleError = "Ваш аккаунт уже подтвержден"
)
