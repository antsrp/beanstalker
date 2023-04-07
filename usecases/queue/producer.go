package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/antsrp/beanstalker/config"
	iqueue "github.com/antsrp/beanstalker/interfaces/queue"
	"github.com/antsrp/beanstalker/usecases/visual"

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
	return visual.Colorize(visual.ColorCyan, fmt.Sprintf("\r\nINSERTED %v\r\n", id)), nil
}

func (p Producer) ListTubes() (string, error) {
	strs, err := p.conn.ListTubes()
	if err != nil {
		return "", fmt.Errorf("ERROR: %w", err)
	}
	return visual.Colorize(visual.ColorCyan, "Tubes:\r\n") + visual.Colorize(visual.ColorGreen, strings.Join(strs, "\r\n")), nil
}

func (p Producer) getTubeStats(tube *beanstalk.Tube) (map[string]string, error) {
	stats, err := tube.Stats()
	if err != nil {
		return nil, fmt.Errorf("can't get stats of tube: %w", err)
	}
	return stats, nil
}

func (p Producer) selectTube(name string) *beanstalk.Tube {
	return beanstalk.NewTube(p.conn, name)
}

func (p Producer) gatherTubesStats(names ...string) (string, error) {
	var sb strings.Builder
	for _, name := range names {
		t := p.selectTube(name)
		stats, err := p.getTubeStats(t)
		if err != nil {
			return "", fmt.Errorf("can't get stats of tube: %w", err)
		}
		sb.WriteString(visual.Colorize(visual.ColorGreen, fmt.Sprintf("\r\ntube: %s\r\n", name)))
		for k, v := range stats {
			if k == "name" {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s: %s\n", visual.Colorize(visual.ColorCyan, k), visual.Colorize(visual.ColorRed, v)))
		}
	}
	return sb.String(), nil
}

func (p Producer) StatsTubes() (string, error) {
	names, err := p.conn.ListTubes()
	if err != nil {
		return "", err
	}
	return p.gatherTubesStats(names...)
}

func (p Producer) StatsTube(tube string) (string, error) {
	return p.gatherTubesStats(tube)
}

func (p *Producer) Close() {
	p.conn.Close()
	fmt.Println("connection closed")
}

var _ iqueue.Producer = Producer{}
