package rpc

import "github.com/mymmsc/go-rocketmq-client/v1/remote"

// request code
const (
	sendMessage                      remote.Code = 10
	pullMessage                                  = 11
	queryMessage                                 = 12
	queryBrokerOffset                            = 13
	queryConsumerOffset                          = 14
	updateConsumerOffset                         = 15
	updateAndCreateTopic                         = 17
	getAllTopicConfig                            = 21
	getTopicConfigList                           = 22
	getTopicNameList                             = 23
	updateBrokerConfig                           = 25
	getBrokerConfig                              = 26
	triggerDeleteFiles                           = 27
	getBrokerRuntimeInfo                         = 28
	searchOffsetByTimestamp                      = 29
	getMaxOffset                                 = 30
	getMinOffset                                 = 31
	getEarliestMsgStoretime                      = 32
	viewMessageByID                              = 33
	heartBeat                                    = 34
	unregisterClient                             = 35
	consumerSendMsgBack                          = 36
	endTransaction                               = 37
	getConsumerListByGroup                       = 38
	CheckTransactionState                        = 39
	NotifyConsumerIdsChanged                     = 40
	lockBatchMQ                                  = 41
	unlockBatchMQ                                = 42
	getAllConsumerOffset                         = 43
	getAllDelayOffset                            = 45
	checkClientConfig                            = 46
	putKvConfig                                  = 100
	getKvConfig                                  = 101
	deleteKvConfig                               = 102
	registerBroker                               = 103
	unregisterBroker                             = 104
	getRouteintoByTopic                          = 105
	getBrokerClusterInfo                         = 106
	updateAndCreateSubscriptiongroup             = 200
	getAllSubscriptiongroupConfig                = 201
	getTopicStatsInfo                            = 202
	getConsumerConnectionList                    = 203
	getProducerConnectionList                    = 204
	wipeWritePermOfBroker                        = 205
	getAllTopicListFromNameserver                = 206
	deleteSubscriptiongroup                      = 207
	getConsumeStats                              = 208
	suspendConsumer                              = 209
	resumeConsumer                               = 210
	resetConsumerOffsetInConsumer                = 211
	resetConsumerOffsetInBroker                  = 212
	adjustConsumerThreadPool                     = 213
	whoConsumeTheMessage                         = 214
	deleteTopicInBroker                          = 215
	deleteTopicInNamesrv                         = 216
	getKvlistByNamespace                         = 219
	ResetConsumerClientOffset                    = 220
	GetConsumerStatusFromClient                  = 221
	invokeBrokerToResetOffset                    = 222
	invokeBrokerToGetConsumerStatus              = 223
	queryTopicConsumeByWho                       = 300
	getTopicsByCluster                           = 224
	registerFilterServer                         = 301
	registerMessageFilterClass                   = 302
	queryConsumeTimeSpan                         = 303
	getSystemTopicListFromNs                     = 304
	getSystemTopicListFromBroker                 = 305
	cleanExpiredConsumequeue                     = 306
	GetConsumerRunningInfo                       = 307
	queryCorrectionOffset                        = 308
	ConsumeMessageDirectlyCode                   = 309
	sendMessageV2                                = 310
	getUnitTopicList                             = 311
	getHasUnitSubTopicList                       = 312
	getHasUnitSubUnunitTopicList                 = 313
	cloneGroupOffset                             = 314
	viewBrokerStatsData                          = 315
	cleanUnusedTopic                             = 316
	getBrokerConsumeStats                        = 317
	updateNamesrvConfig                          = 318
	getNamesrvConfig                             = 319
	sendBatchMessage                             = 320
	queryConsumeQueue                            = 321
)

// response code
const (
	UnknowError                remote.Code = -1
	Success                                = 0
	SystemError                            = 1
	SystemBusy                             = 2
	RequestCodeNotSupported                = 3
	TransactionFailed                      = 4
	FlushDiskTimeout                       = 10
	SlaveNotAvailable                      = 11
	FlushSlaveTimeout                      = 12
	MessageIllegal                         = 13
	ServiceNotAvailable                    = 14
	VersionNotSupported                    = 15
	NoPermission                           = 16
	TopicNotExist                          = 17
	TopicExistAlready                      = 18
	PullNotFound                           = 19
	PullRetryImmediately                   = 20
	PullOffsetMoved                        = 21
	QueryNotFound                          = 22
	SubscriptionParseFailed                = 23
	SubscriptionNotExist                   = 24
	SubscriptionNotLatest                  = 25
	SubscriptionGroupNotExist              = 26
	FilterDataNotExist                     = 27
	FilterDataNotLatest                    = 28
	TransactionShouldCommit                = 200
	TransactionShouldRollback              = 201
	TransactionStateUnknow                 = 202
	TransactionStateGroupWrong             = 203
	NoBuyerID                              = 204
	NotInCurrentUnit                       = 205
	ConsumerNotOnline                      = 206
	ConsumeMsgTimeout                      = 207
	NoMessage                              = 208
	ConnectBrokerException                 = 10001
	AccessBrokerException                  = 10002
	BrokerNotExistException                = 10003
	NoNameServerException                  = 10004
	NotFoundTopicException                 = 10005
)
