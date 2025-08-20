package app

type ReturnMessage struct {
	Text string
	Err  error
}

const defaultFailMsg = "Что-то пошло не так"
