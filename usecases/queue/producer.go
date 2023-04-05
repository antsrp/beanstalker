package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/antsrp/beanstalker/config"
	iqueue "github.com/antsrp/beanstalker/interfaces/queue"

	"github.com/beanstalkd/go-beanstalk"
)

type Settings struct {
	priority uint32
	delay    time.Duration
	ttr      time.Duration
}

type Producer struct {
	conn     *beanstalk.Conn
	tube     *beanstalk.Tube
	settings *Settings
}

func NewProducer(host string, port int, s *config.Settings) (*Producer, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	c, err := beanstalk.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	p := &Producer{
		conn:     c,
		tube:     &c.Tube,
		settings: &Settings{},
	}
	p.SetTube(s.Tube)
	p.SetPriority(s.Priority)
	p.SetTTR(s.TTR)
	p.SetDelay(s.Delay)

	return p, nil
}

func (p Producer) SetTube(tube string) {
	*p.tube = *beanstalk.NewTube(p.conn, tube)
}

func (p Producer) SetDelay(delay int) {
	p.settings.delay = time.Duration(delay) * time.Second
}

func (p Producer) SetPriority(priority uint32) {
	p.settings.priority = priority
}

func (p Producer) SetTTR(ttr int) {
	p.settings.ttr = time.Duration(ttr) * time.Second
}

func (p Producer) Put(data []byte) (string, error) {
	id, err := p.tube.Put(data, p.settings.priority, p.settings.delay, p.settings.ttr)
	if err != nil {
		return "", fmt.Errorf("ERROR: %w", err)
	}
	return fmt.Sprintf("INSERTED %v\n", id), nil
}

func (p Producer) ListTubes() (string, error) {
	strs, err := p.conn.ListTubes()
	if err != nil {
		return "", fmt.Errorf("ERROR: %w", err)
	}
	return "Tubes\r\n" + strings.Join(strs, "\r\n"), nil
}

func (p *Producer) Close() {
	p.conn.Close()
	fmt.Println("connection closed")
}

var _ iqueue.Producer = Producer{}
