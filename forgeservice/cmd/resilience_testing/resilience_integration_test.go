package resilience_testing

import (
	"errors"
	"fmt"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	toxiproxy "github.com/Shopify/toxiproxy/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"testing"
)

var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker open, denying requests to save resources")
)

var application app.App

// connections
var SMELTERCONNECTION_NOCOLONS	string
var TOXIPROXY_CLIENTCONNECTION  string

var SMELTERCONNECTION_DEFAULT 	string


// ToxiProxy
var toxiClient *toxiproxy.Client
//var proxyForge *toxiproxy.Proxy
var proxySmelter *toxiproxy.Proxy
//var proxyMongoDB *toxiproxy.Proxy


func TestMain(m *testing.M) {
	// Reading configuration
	application = app.App{}
	config := configuration.ReadConfiguration(true)
	application.Configuration = &config

	// Set Connection URLS
	SMELTERCONNECTION_NOCOLONS = application.Configuration.SmelterConnectionNoColons

	// Set ToxiProxy Connection URLS
	TOXIPROXY_CLIENTCONNECTION = application.Configuration.ToxiProxy_ClientConnection

	// Remember Default Connections
	SMELTERCONNECTION_DEFAULT = application.Configuration.SmelterConnectionDefault

	// Initialize logging
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	application.Initialize()
	os.Exit(m.Run())
}

// currently not working, start .exe manually before resilience testing
func StartToxiProxyServer(){
	PathToToxiProxyExecutable := "/toxiproxy/toxiproxy-server-windows-amd64.exe"
	cmdToxiProxy := &exec.Cmd{
		Path:         PathToToxiProxyExecutable,
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}
	cmdToxiProxy.Run()
}

// Use this method everytime you want to test any service for resilient behaviour
// it creates a new client, removes all toxics and returns the client.
// In each method you can add the desired toxics with client.AddToxic(...)
// Make sure to remove the toxics and proxies with defer, even though they are deleted here anyway
func InitToxiProxy() (*toxiproxy.Client, error) {
	var err error
	toxiClient := toxiproxy.NewClient(TOXIPROXY_CLIENTCONNECTION)
	// Delete all proxies
	proxies, err := toxiClient.Proxies()
	if err != nil {
		panic("Couldnt connect toxiClient to ToxiProxy HTTP Server. Make sure it is running on Port 8474")
	}
	for _, v := range proxies {
		err = v.Delete()
		if err != nil {
			panic("Couldnt create Proxy on ToxiClient, make sure the Smelterservice and the toxiproxy server are running")
		}
	}
	proxySmelter, err = toxiClient.CreateProxy("proxySmelter", SMELTERCONNECTION_NOCOLONS, SMELTERCONNECTION_DEFAULT)
	if err != nil {
		panic("Couldnt create Proxy on ToxiClient, make sure the Smelterservice and the toxiproxy server are running")
	}
	return toxiClient, err
}

// helper to validate sword response
func Validate_Sword(sword entities.Sword) bool{
	return (sword.Type != "" && sword.Weight != 0)
}


// Integration Tests

// Single request GetIron from Smelterservice
// Success: Valid Iron Response
// Fails if enot valid Iron response
func Test_Integration_GetIronFromSmelterservice(t *testing.T){
	if testing.Short() {
		t.Skip("skipping Integration Test: TestGetIronFromSmelterservice")
	}
	toxiClient, _ = InitToxiProxy()
	iron, err := application.ForgeService.GetIronFromSmelterservice()
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, reflect.TypeOf(entities.Iron{}), reflect.TypeOf(iron))
}

// Single request GetSword
// Success: Valid Sword Response
// Fails if enot valid sword response
func Test_Integration_GetSword_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Integration Test: TestGetSwords")
	}
	toxiClient, _ = InitToxiProxy()
	sword, err := application.ForgeService.GetSword()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t,Validate_Sword(sword))
}

// Resilience Test

// Single request GetSword with timeout of the smelterservice
// Success: When timeout, error should be the Circuit Breaker Opening to save resources
// Fails if error is not ErrCircuitBreakerOpen
func Test_Integration_CircuitBreakerOpens_Timeout_GetSword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: Test_CircuitBreakerOpens_GetIron")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("timeout_smelter", "timeout", "", 1, toxiproxy.Attributes{})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}

	r, err :=  application.ForgeService.GetSword()
	defer proxySmelter.RemoveToxic("timeout_smelter")
	if err != nil && err.Error() == ErrCircuitBreakerOpen.Error() {
		logrus.Debug(err.Error(), r)
		return
	} else {
		t.Fail()
	}
}

// Single request GetSword with latency of the smelterservice
// Success: When latency is too high, error should be the Circuit Breaker Opening to save resources
// Fails if error is not ErrCircuitBreakerOpen
func Test_Integration_CircuitBreakerOpens_Latency_GetSword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: Test_CircuitBreakerOpens_GetIron")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("latency_downstream", "latency", "", 1, toxiproxy.Attributes{
		"latency": 3000,
	})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}

	r, err :=  application.ForgeService.GetSword()

	defer proxySmelter.RemoveToxic("latency_downstream")
	if err != nil && err.Error() == ErrCircuitBreakerOpen.Error() {
		logrus.Debug(err.Error(), r)
		return
	} else {
		t.Fail()
	}
}

// Bulk requests with a temporary timeout of the smelterservice
// Success: When timeout, error should be the Circuit Breaker Opening to save resources
// Fails if error is not ErrCircuitBreakerOpen
func Test_Integration_BulkGetSword_TemporarySmelterTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: TestSmelterConnectionSlow")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("timeout_smelter", "timeout", "", 0, toxiproxy.Attributes{})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}

	// simulate a temporary timeout
	// make a bulk of requests to to forgeservice including internal requests to smelterservice which has occasional timeouts,
	//first 19 success, 20 to 179 failure, rest success
	for i := 0; i < 10; i++ {

		if i > 2 && i < 7 {
			if i == 3 {
				proxySmelter.UpdateToxic("timeout_smelter", 1, toxiproxy.Attributes{})
			}
			r, err :=  application.ForgeService.GetSword()
			if err != nil {
				logrus.Debug("Service Timeout HTTP Status: "+err.Error())
				if err.Error() == ErrCircuitBreakerOpen.Error() {
					logrus.Debug(err.Error(), r)
					return
				} else {
					t.Fail()
				}
			}
		}
		 if i < 3 || i > 6 {
			if i == 7{
				proxySmelter.UpdateToxic("timeout_smelter", 0, toxiproxy.Attributes{})
			}
			 r, err := application.ForgeService.GetSword()
			 if err != nil {
				 t.Fatal(err)
			 }
			 logrus.Debug("Response from Smelterconnection: ",r)
		}

	}
	// clean up toxic later
	defer proxySmelter.RemoveToxic("timeout_smelter")

}

// Bulk requests with a temporary latency of the smelterservice
// Success: When latency is too high, error should be the Circuit Breaker Opening to save resources
// Fails if error is not ErrCircuitBreakerOpen
func Test_Integration_BulkGetSword_TemporarySmelterLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: TestSmelterConnectionSlow")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("latency_downstream", "latency", "", 0, toxiproxy.Attributes{
		"latency": 3000,
	})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}

	// simulate a temporary latency
	// make a bulk of requests to to forgeservice including internal requests to smelterservice which has occasional timeouts,
	//first 19 success, 20 to 179 failure, rest success
	for i := 0; i < 5; i++ {

		if i == 4 {
			logrus.Debug("NOW TOXIC ACTIVE: latency_downstream")
			proxySmelter.UpdateToxic("latency_downstream", 1, toxiproxy.Attributes{})
			r, err :=  application.ForgeService.GetSword()
			if err != nil {
				logrus.Debug("Service Timeout HTTP Status: "+err.Error())
			}
			logrus.Debug("Response from timed out service "+"http://localhost:8080/sword",r)
		} else {
			r, err :=  application.ForgeService.GetSword()
			if err != nil {
				t.Fail()
			}
			logrus.Debug("Response from smelterservice: ",r)
		}

	}
	// clean up toxic later
	defer proxySmelter.RemoveToxic("latency_downstream")

}

// Bulk requests with a temporary latency and jittering of the smelterservice connection
// Success: When latency is too high, error should be the Circuit Breaker Opening to save resources
// Fails if error is not ErrCircuitBreakerOpen
func Test_Integration_BulkGetSword_SmelterLatencyJittering(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: TestSmelterConnectionSlow")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("latency_downstream", "latency", "", 1, toxiproxy.Attributes{
		"latency": 2500,
		"jittering": 1000,
	})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}
	logrus.Debug("NOW TOXIC ACTIVE: latency_downstream")

	// simulate a temporary latency with high jittering
	for i := 0; i < 5; i++ {
		r, err :=  application.ForgeService.GetSword()
		if err != nil {
			logrus.Debug("Service Timeout HTTP Status: "+err.Error())
		}
		logrus.Debug("Response from timed out service "+"http://localhost:8080/sword",r)
	}
	// clean up toxic later
	defer proxySmelter.RemoveToxic("latency_downstream")

}

func Test_Integration_BulkGetSword_SmelterBandwidth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: TestSmelterConnectionSlow")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("low_bandwidth_smelter", "bandwidth", "", 1, toxiproxy.Attributes{
		"rate": 0,
	})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}
	logrus.Debug("NOW TOXIC ACTIVE: low_bandwidth_smelter")

	// simulate a temporary latency with high jittering
	for i := 0; i < 2; i++ {
		r, err :=  application.ForgeService.GetSword()
		if err != nil {
			logrus.Debug("Service Timeout HTTP Status: "+err.Error())
		}
		logrus.Debug("Response from slow service "+"http://localhost:8080/sword",r)
	}
	// clean up toxic later
	defer proxySmelter.RemoveToxic("low_bandwidth_smelter")

}

func Test_Integration_BulkGetSword_Smelter_limit_data(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: TestSmelterConnectionSlow")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}
	_, err = proxySmelter.AddToxic("limit_data_smelter", "limit_data", "", 1, toxiproxy.Attributes{
		"bytes": 100,
	})
	if err != nil{
		logrus.Error("Could not add toxic to smelter proxy")
		t.Fatal(err)
	}
	logrus.Debug("NOW TOXIC ACTIVE: limit_data_smelter")

	for i := 0; i < 2; i++ {
		r, err :=  application.ForgeService.GetSword()
		if err != nil {
			logrus.Debug("Service Timeout HTTP Status: "+err.Error())
		}
		logrus.Debug("Response from limited service "+"http://localhost:8080/sword",r)
	}
	// clean up toxic later
	defer proxySmelter.RemoveToxic("limit_data_smelter")

}

func workerGetIron(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	sword, err := application.ForgeService.GetSword()
	logrus.Debugln(sword, err)
	fmt.Printf("Worker %d done\n", id)
}

// makes a bunch of concurrent requests to GetIron that wait for each other to finish
// Success: No errors
func Test_Integration_ConcurrentRequests_GetIron(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Resilience Test: Test_Integration_ConcurrentRequests_GetIron")
	}
	var err error
	toxiClient, err = InitToxiProxy()
	if err != nil {
		logrus.Error("Could not create proxy listen on smelter connection")
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go workerGetIron(i, &wg)
	}
	wg.Wait()
}


