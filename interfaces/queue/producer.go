package queue

type Producer interface {
	Put([]byte) (string, error)
	ListTubes() (string, error)
	StatsTubes() (string, error)
	StatsTube(string) (string, error)

	SetTube(string)
	SetDelay(int)
	SetPriority(uint32)
	SetTTR(int)
}
