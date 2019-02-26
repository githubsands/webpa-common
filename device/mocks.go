package device

import (
	"net/http"

	"github.com/Comcast/webpa-common/convey"
	"github.com/stretchr/testify/mock"
)

type MockConnector struct {
	mock.Mock
}

var _ Connector = (*MockConnector)(nil)

func (m *MockConnector) Connect(response http.ResponseWriter, request *http.Request, header http.Header) (Interface, error) {
	arguments := m.Called(response, request, header)
	first, _ := arguments.Get(0).(Interface)
	return first, arguments.Error(1)
}

func (m *MockConnector) Disconnect(id ID) bool {
	return m.Called(id).Bool(0)
}

func (m *MockConnector) DisconnectIf(predicate func(ID) bool) int {
	return m.Called(predicate).Int(0)
}

func (m *MockConnector) DisconnectAll() int {
	return m.Called().Int(0)
}

type MockRegistry struct {
	mock.Mock
}

var _ Registry = (*MockRegistry)(nil)

func (m *MockRegistry) Len() int {
	return m.Called().Int(0)
}

func (m *MockRegistry) Get(id ID) (Interface, bool) {
	arguments := m.Called(id)
	first, _ := arguments.Get(0).(Interface)
	return first, arguments.Bool(1)
}

func (m *MockRegistry) VisitAll(f func(Interface) bool) int {
	return m.Called(f).Int(0)
}

type MockDevice struct {
	mock.Mock
}

func (m *MockDevice) String() string {
	return m.Called().String(0)
}

func (m *MockDevice) MarshalJSON() ([]byte, error) {
	arguments := m.Called()
	return arguments.Get(0).([]byte), arguments.Error(1)
}

func (m *MockDevice) ID() ID {
	return m.Called().Get(0).(ID)
}

func (m *MockDevice) Pending() int {
	return m.Called().Int(0)
}

func (m *MockDevice) Close() error {
	return m.Called().Error(0)
}

func (m *MockDevice) Closed() bool {
	arguments := m.Called()
	return arguments.Bool(0)
}

func (m *MockDevice) Statistics() Statistics {
	arguments := m.Called()
	first, _ := arguments.Get(0).(Statistics)
	return first
}

func (m *MockDevice) Convey() convey.Interface {
	arguments := m.Called()
	first, _ := arguments.Get(0).(convey.Interface)
	return first
}

func (m *MockDevice) ConveyCompliance() convey.Compliance {
	arguments := m.Called()
	first, _ := arguments.Get(0).(convey.Compliance)
	return first
}

func (m *MockDevice) PartnerIDs() []string {
	arguments := m.Called()
	first, _ := arguments.Get(0).([]string)
	return first
}

func (m *MockDevice) SatClientID() string {
	arguments := m.Called()
	first, _ := arguments.Get(0).(string)
	return first
}

func (m *MockDevice) Trust() Trust {
	arguments := m.Called()
	first, _ := arguments.Get(0).(Trust)
	return first
}

func (m *MockDevice) Send(request *Request) (*Response, error) {
	arguments := m.Called(request)
	first, _ := arguments.Get(0).(*Response)
	return first, arguments.Error(1)
}
