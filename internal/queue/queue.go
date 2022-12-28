package queue

type Queue interface {
	MakeCIMsgChan() (<-chan []byte, error)
	Close()
}
