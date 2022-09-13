package service

import (
	"context"
	"encoding/json"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/configuration"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/entities"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/app/domain/repositories"
	"git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/circuitbreaker"
	_ "git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/circuitbreaker"
	_ "git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// Resilience Constants
const MAX_LATENCY_MS  = 6000*time.Millisecond
const MAX_TIMEOUT_MS  = 1000*time.Millisecond
const MAX_RETRIES  = 10
const MAX_CONCURRENT_REQUESTS = 2
const CIRCUITBREAKER_THRESHOLD = 3

// Use Case Constants
const IRONBADGE_AMOUNT  = 5

var cb circuitbreaker.CircuitBreaker
var bh = make(chan int, MAX_CONCURRENT_REQUESTS)

type ForgeService interface {
	GetSword()(entities.Sword, error)
	GetIronFromSmelterservice()(entities.Iron, error)
	GetIron_CircuitBreaker() (entities.Iron, error)
}

type ForgeServiceImpl struct {
	swordRepository    repositories.SwordRepository
	configuration      *configuration.Configuration
}

func NewForgeService(
	swordRepository repositories.SwordRepository,
	configuration *configuration.Configuration) ForgeService {

	return &ForgeServiceImpl{
		swordRepository: swordRepository,
		configuration:   configuration,
	}
}


func (forgeServiceImpl *ForgeServiceImpl) GetSword()(entities.Sword, error) {
	logrus.Debug("GetSword invoked")
	return ForgeSword(forgeServiceImpl)
}

// API for requesting a sword. It will request iron from the smelterservice
// It will retry to request the iron until success
func ForgeSword(forgeServiceImpl *ForgeServiceImpl) (entities.Sword, error){
	logrus.Debug("ForgeSword invoked")
	var err error

	sword := entities.Sword{}

	var ironBadge [IRONBADGE_AMOUNT]entities.Iron
	ch_iron := make(chan entities.Iron)
	ch_err := make(chan error)
	for i := 0; i < IRONBADGE_AMOUNT; i++ {

		bh <- 1 // Acquire worker from bulkhead thread pool

		logrus.Debug("IronBadge Worker Number: "+strconv.Itoa(i+1))

		go func(){
			logrus.Debug("Starting GoRotuine for IronBadge Worker starts trying to request Iron")
			forgeServiceImpl.GetIron_Retry(ch_iron, ch_err)
			<-bh // release worker from bulkhead thread pool
		}()

		ironBadge[i] = <- ch_iron
		err = <- ch_err

		// Stop getting a full IronBadge, when Circuit Breaker Opens
		if err == cb.ErrCircuitBreakerOpen{
			logrus.Debug(err.Error())
			break
		}

	}
	close(ch_iron)
	close(ch_err)
	if Validate_IronBadge(ironBadge){
		sword = entities.Sword{
			CreatedAt: time.Now(),
			Weight:    10,
			Quality:   10,
			Sharpened: true,
			Type:      "Sword",
		}
	}
	return sword, err
}

// wraps retry around get iron
func (forgeServiceImpl *ForgeServiceImpl) GetIron_Retry(ch_iron chan entities.Iron, ch_err chan error) (entities.Iron, error) {
	var err error
	iron := entities.Iron{}
	cb, err = circuitbreaker.NewCircuitBreaker(CIRCUITBREAKER_THRESHOLD)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	for retryCounter := 0; retryCounter < MAX_RETRIES; retryCounter++ {
		iron, err = forgeServiceImpl.GetIron_CircuitBreaker()
		// if the Circuit Breaker is open, we stop retrying
		if err == nil || err == cb.ErrCircuitBreakerOpen {
			break
		}
		// if the Circuit Breaker is not open, we can nil the err and try again
		err = nil
	}
	ch_iron <- iron
	ch_err <- err
	return iron, err
}

// wraps CircuitBreaker Logic around request
func (forgeServiceImpl *ForgeServiceImpl) GetIron_CircuitBreaker() (entities.Iron, error) {
	logrus.Debug("GetIron_CircuitBreaker invoked")
	var err error
	var iron entities.Iron
	if cb.IsClosed() {
		iron, err = forgeServiceImpl.GetIronFromSmelterservice()

		if err != nil {
			cb.IncrementErrorCounter()
			logrus.Debug("GetIron Circuit Breaker Error Counter: "+strconv.Itoa(cb.GetErrorCount()))
			logrus.Debug("Error from GetIron: "+err.Error())
			if cb.GetErrorCount() >= cb.GetThreshold() {
				cb.Open()
				err = cb.ErrCircuitBreakerOpen
			}
		}

	} else {
		err = cb.ErrCircuitBreakerOpen
		return iron, err
	}
	return iron, err
}

// actual request to get iron from smelterservice
func (forgeServiceImpl *ForgeServiceImpl) GetIronFromSmelterservice() (entities.Iron, error) {
	logrus.Debug("GetIronFromSmelterservice invoked")
	var err error
	iron := entities.Iron{}

	req, err := http.NewRequest("GET", forgeServiceImpl.configuration.SmelterConnection+"/iron", nil)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), MAX_TIMEOUT_MS)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return iron, err
	}

	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&iron)
	logrus.Debug("GetIron Request: "+forgeServiceImpl.configuration.SmelterConnection+"/iron")
	logrus.Debug("-> with Response:"+"StatusCode: "+strconv.Itoa(res.StatusCode) + " "+res.Status)
	return iron, err
}

// validates the ironBadge of three iron to be all of type entities.Iron
func Validate_IronBadge(ironBadge [IRONBADGE_AMOUNT]entities.Iron) bool{
	for i := 0; i < 3; i++ {
		ironT := ironBadge[i]
		if reflect.TypeOf(ironT) != reflect.TypeOf(entities.Iron{}){
			return false
		}
	}
	return true
}
