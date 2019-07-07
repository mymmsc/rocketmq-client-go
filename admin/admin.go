package admin

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mymmsc/go-rocketmq-client/v1"
	"github.com/mymmsc/go-rocketmq-client/v1/client"
	"github.com/mymmsc/go-rocketmq-client/v1/client/rpc"
	"github.com/mymmsc/go-rocketmq-client/v1/log"
	"github.com/mymmsc/go-rocketmq-client/v1/message"
	"github.com/mymmsc/go-rocketmq-client/v1/route"
)

// Admin admin operations
type Admin struct {
	rocketmq.Server
	rocketmq.Client

	client mqClient

	logger log.Logger
}

// New create admin operator
func New(group string, namesAdders []string, logger log.Logger) *Admin {
	a := &Admin{
		Client: rocketmq.Client{
			HeartbeatBrokerInterval:       30 * time.Second,
			PollNameServerInterval:        30 * time.Second,
			PersistConsumerOffsetInterval: 5 * time.Second,
			NameServerAddrs:               namesAdders,
			GroupName:                     group,
		},
	}
	a.StartFunc = a.start
	a.logger = logger
	return a
}

// Start admin work
func (a *Admin) start() (err error) {
	a.InstanceName = strconv.Itoa(os.Getpid())
	a.ClientIP, err = rocketmq.GetIPStr()
	if err != nil {
		a.logger.Errorf("no ip")
		return
	}
	mqClient, err := a.buildMQClient()
	if err != nil {
		return
	}
	a.client = mqClient

	err = a.client.RegisterAdmin(a)
	if err != nil {
		a.logger.Errorf("register producer error:%s", err.Error())
		return
	}

	err = mqClient.Start()
	if err != nil {
		a.logger.Errorf("start mq client error:%s", err)
		return
	}
	a.buildShutdowner(mqClient.Shutdown)

	return
}

func (a *Admin) buildMQClient() (*client.MQClient, error) {
	a.ClientID = client.BuildMQClientID(a.ClientIP, a.UnitName, a.InstanceName)
	c, err := client.New(
		&client.Config{
			HeartbeatBrokerInterval: a.HeartbeatBrokerInterval,
			PollNameServerInterval:  a.PollNameServerInterval,
			NameServerAddrs:         a.NameServerAddrs,
		}, a.ClientID, a.logger,
	)
	a.client = c
	return c, err
}

func (a *Admin) buildShutdowner(f func()) {
	shutdowner := &rocketmq.ShutdownCollection{}
	shutdowner.AddLastFuncs(
		func() {
			a.logger.Infof("shutdown admin:%s START", a.GroupName)
		},
		a.shutdown, f,
		func() {
			a.logger.Infof("shutdown admin:%s END", a.GroupName)
		},
	)

	a.Shutdowner = shutdowner
}

func (a *Admin) shutdown() {
	a.client.UnregisterAdmin(a.GroupName)
}

// Group returns the GroupName of the producer
func (a *Admin) Group() string {
	return a.GroupName
}

// CreateOrUpdateTopic create a new topic in the broker
func (a *Admin) CreateOrUpdateTopic(brokerAddr, topic string, perm, queueCount int32) error {
	header := &rpc.CreateOrUpdateTopicHeader{
		Topic:           topic,
		ReadQueueNums:   queueCount,
		WriteQueueNums:  queueCount,
		DefaultTopic:    rocketmq.DefaultTopic,
		Perm:            perm,
		TopicFilterType: SingleTag.String(),
	}

	if err := a.client.CreateOrUpdateTopic(brokerAddr, header, 3*time.Second); err != nil {
		return err
	}
	return nil
}

// DeleteTopicInBroker delete the topic in the broker
func (a *Admin) DeleteTopicInBroker(addr, topic string) (err error) {
	err = a.client.DeleteTopicInBroker(addr, topic, 3*time.Second)
	if err != nil {
		a.logger.Errorf("delete topic %s in broker:%s error:%s", topic, addr, err)
		return
	}

	a.logger.Debugf("DELETE topic %s suc at broker %s", topic, addr)
	return
}

// DeleteTopicInAllNamesrv delete the topic in the namesrv
func (a *Admin) DeleteTopicInAllNamesrv(topic string) (err error) {
	for _, addr := range a.NameServerAddrs {
		err = a.client.DeleteTopicInNamesrv(addr, topic, 3*time.Second)
		if err != nil {
			a.logger.Errorf("delete topic %s in namesrv:%s error:%s", topic, addr, err)
			continue
		}
		a.logger.Debugf("DELETE topic %s suc at namesrv %s", topic, addr)
	}
	return
}

// GetBrokerClusterInfo get broker cluster info
func (a *Admin) GetBrokerClusterInfo() (info *route.ClusterInfo, err error) {
	l := len(a.NameServerAddrs)
	for i, c := rand.Intn(l), l; c > 0; i, c = i+1, c-1 {
		addr := a.NameServerAddrs[i%l]
		info, err = a.client.GetBrokerClusterInfo(addr, 3*time.Second)
		if err == nil {
			return
		}

		a.logger.Errorf("request broker cluster info from %s, error:%s", addr, err)
	}
	return
}

// QueryMessageByID querys the message by message id
// offsetID is the offset id generated by the broker server
func (a *Admin) QueryMessageByID(offsetID string) (*message.Ext, error) {
	addr, offset, err := message.ParseMessageID(offsetID)
	if err != nil {
		return nil, err
	}

	return a.client.QueryMessageByOffset(addr.String(), offset, 3*time.Second)
}

// MaxOffset fetches the max offset of the consume queue
func (a *Admin) MaxOffset(q *message.Queue) (int64, error) {
	addr, err := a.client.FindBrokerAddr(q.BrokerName, rocketmq.MasterID, false)
	if err != nil {
		err = a.client.UpdateTopicRouterInfoFromNamesrv(q.Topic)
		if err != nil {
			return -1, err
		}

		addr, err = a.client.FindBrokerAddr(q.BrokerName, rocketmq.MasterID, false)
		if err != nil {
			return -1, err
		}
	}

	return a.client.MaxOffset(addr.Addr, q.Topic, uint8(q.QueueID), 3*time.Second)
}

// GetConsumerIDs get the consumer ids from the broker
func (a *Admin) GetConsumerIDs(addr, group string) ([]string, error) {
	return a.client.GetConsumerIDs(addr, group, time.Second*3)
}

// ResetConsumeOffset requests the broker to reset the offsets of the specified topic, the offsets' owner
// is specified by the group
func (a *Admin) ResetConsumeOffset(
	broker, topic, group string, timestamp time.Time, isForce bool,
) (
	map[message.Queue]int64, error,
) {
	addr, err := a.client.FindBrokerAddr(broker, rocketmq.MasterID, true)
	if err != nil {
		err = a.client.UpdateTopicRouterInfoFromNamesrv(topic)
		if err != nil {
			return nil, err
		}

		addr, err = a.client.FindBrokerAddr(broker, rocketmq.MasterID, true)
		if err != nil {
			return nil, err
		}
	}

	return a.client.ResetConsumeOffset(addr.Addr, topic, group, timestamp, isForce, 3*time.Second)
}

// TopicFilter details
type TopicFilter int8

func (f TopicFilter) String() string {
	switch f {
	case SingleTag:
		return "SINGLE_TAG"
	case MultiTag:
		return "MULTI_TAG"
	default:
		panic("BUG:unknow topic filter:" + strconv.Itoa(int(f)))
	}
}

// TopicFilter defination
const (
	SingleTag TopicFilter = iota
	MultiTag
)

// ConsumeMessageDirectly sends request to the broker to push the message specified by the id
// to the client in the specified group
func (a *Admin) ConsumeMessageDirectly(group, clientID, offsetID string) (ret client.ConsumeMessageDirectlyResult, err error) {
	addr, _, err := message.ParseMessageID(offsetID)
	if err != nil {
		return
	}
	return a.client.ConsumeMessageDirectly(addr.String(), group, clientID, offsetID)
}
