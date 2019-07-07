package client

import (
	"testing"
	"time"

	"github.com/mymmsc/go-rocketmq-client/v2/client/rpc"
	"github.com/mymmsc/go-rocketmq-client/v2/log"
	"github.com/mymmsc/go-rocketmq-client/v2/remote"
	"github.com/stretchr/testify/assert"
)

type fakeRemoteClient struct {
	remote.FakeClient

	requestSyncErr error
	command        remote.Command
}

func (f *fakeRemoteClient) RequestSync(string, *remote.Command, time.Duration) (*remote.Command, error) {
	return &f.command, f.requestSyncErr
}

func fakeClient() *MQClient {
	c, err := newMQClient(&Config{}, "", log.Std)
	if err != nil {
		panic(err)
	}

	c.Client = &fakeRemoteClient{}
	return c
}

func TestSendMessageSync(t *testing.T) {
	c := fakeClient()

	// no broker
	resp, err := c.SendMessageSync("", []byte{}, &rpc.SendHeader{}, time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	c.brokerAddrs.put("", map[int32]string{0: "a"})
	resp, err = c.SendMessageSync("", []byte{}, &rpc.SendHeader{}, time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}
