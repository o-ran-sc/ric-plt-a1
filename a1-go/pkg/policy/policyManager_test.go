package policy

import (
	"os"
	"strconv"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SdlMock struct {
	mock.Mock
}

var sdlInst *SdlMock
var pm *PolicyManager

func TestMain(m *testing.M) {
	sdlInst = new(SdlMock)
	a1.Init()
	pm = createPolicyManager(sdlInst)
	code := m.Run()
	os.Exit(code)
}
func TestSetPolicyInstance(t *testing.T) {
	var policyTypeId int
	policyTypeId = 20001
	var policyInstanceID int
	policyInstanceID = 123456
	var status string
	status = "OK"
	instancehandlerKey := a1HandlerPrefix + strconv.FormatInt(20001, 10) + "." + strconv.FormatInt(int64(policyInstanceID), 10)
	instancearr := []interface{}{instancehandlerKey, status}
	sdlInst.On("Set", "A1m_ns", instancehandlerKey, instancearr).Return(nil)
	errresp := pm.SetPolicyInstanceStatus(policyTypeId, policyInstanceID, status)
	assert.NoError(t, errresp)

	sdlInst.AssertExpectations(t)
}

func (s *SdlMock) Set(ns string, pairs ...interface{}) error {
	args := s.MethodCalled("Set", ns, pairs)
	return args.Error(0)
}

func (s *SdlMock) Get(ns string, keys []string) (map[string]interface{}, error) {
	a1.Logger.Error("Get Called ")
	return map[string]interface{}{}, nil
}

func (s *SdlMock) GetAll(ns string) ([]string, error) {
	args := s.MethodCalled("GetAll", ns)
	return args.Get(0).([]string), nil
}
