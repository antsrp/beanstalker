package queue

type Producer interface {
	Put([]byte) (string, error)
	ListTubes() (string, error)

	SetTube(string)
	SetDelay(int)
	SetPriority(uint32)
	SetTTR(int)
}
