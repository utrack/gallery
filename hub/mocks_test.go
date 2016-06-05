package hub

import (
	"encoding/json"
	"github.com/utrack/gallery/messages"
)

// lsiterMock mocks the storage.Lister.
type listerMock struct {
	retList messages.FilesInfo
	retErr  error
}

func (m *listerMock) GetList() (messages.FilesInfo, error) {
	return m.retList, m.retErr
}

type notifierMock struct {
	retChan chan messages.FileChangeNotification
}

func (m *notifierMock) GetNotificationChan() <-chan messages.FileChangeNotification {
	return m.retChan
}

func (m *notifierMock) Close() error {
	return nil
}

type connMock struct {
	rcvChan    chan json.RawMessage
	disconChan chan error
}

func (m *connMock) Send(msg json.RawMessage) {
	m.rcvChan <- msg
}

func (m *connMock) DisconChan() <-chan error {
	return m.disconChan
}

func (m *connMock) Disconnect() {

}
