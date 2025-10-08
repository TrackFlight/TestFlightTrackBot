package testflight

import (
	"errors"
	"sync"

	"github.com/Laky-64/http"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db"
	"github.com/TrackFlight/TestFlightTrackBot/internal/db/models"
	"github.com/TrackFlight/TestFlightTrackBot/internal/tor"
	"github.com/TrackFlight/TestFlightTrackBot/internal/utils"
)

func (c *Client) Check(links []db.GetUsedLinksLinkRow) (map[string]Result, error) {
	request, err := http.ExecuteRequest(AwesomeTestFlightURL)
	if err != nil {
		return nil, err
	}
	if request.StatusCode != 200 {
		return nil, errors.New("failed to fetch awesome testflight links")
	}
	awesomeList := RegexAwesomeTestFlight.FindAllStringSubmatch(request.String(), -1)
	awesomeAppNames := make(map[string]string)
	for _, item := range awesomeList {
		awesomeAppNames[item[3]] = item[2]
	}
	transaction, err := c.TorClient.NewTransaction(len(links))
	if err != nil {
		return nil, err
	}
	defer transaction.Close()
	results := make(map[string]Result, len(links))
	var wg sync.WaitGroup
	var mu sync.Mutex
	pool := utils.NewPool(tor.MaxRequestsPerInstance * c.TorClient.InstanceCount())
	for _, link := range links {
		wg.Add(1)
		pool.Enqueue(func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				res, err := transaction.ExecuteRequest(
					link.URL,
					pickRandomUserAgent(c.userAgents),
				)
				mu.Lock()
				if err != nil {
					if res != nil && res.StatusCode == 404 {
						results[link.URL] = Result{
							Status: models.LinkStatusEnumInvalid,
						}
					} else {
						results[link.URL] = Result{
							Error: err,
						}
					}
					mu.Unlock()
					return
				}
				bodyString := res.String()
				if len(bodyString) < 300 {
					mu.Unlock()
					continue
				}
				var appName, appIcon, description string
				var status models.LinkStatusEnum
				var isPublic bool
				rawAppName := RegexAppName.FindAllStringSubmatch(bodyString, -1)
				if len(rawAppName) == 0 {
					status = models.LinkStatusEnumClosed
				} else {
					rawStatus := RegexStatus.FindAllStringSubmatch(bodyString, -1)
					if rawStatus[0][1] == "This beta is full." {
						status = models.LinkStatusEnumFull
					} else {
						status = models.LinkStatusEnumAvailable
					}
					appName = rawAppName[0][1]
					appIcon = RegexAppIcon.FindStringSubmatch(bodyString)[1]
					description = RegexDescription.FindStringSubmatch(bodyString)[1]
				}
				if awesomeAppName, exists := awesomeAppNames[link.URL]; exists {
					if len(appName) == 0 {
						appName = awesomeAppName
					}
					isPublic = true
				}

				results[link.URL] = Result{
					AppName:     appName,
					IconURL:     appIcon,
					Description: description,
					Status:      status,
					IsPublic:    isPublic,
				}
				mu.Unlock()
				return
			}
			results[link.URL] = Result{
				Error: errors.New("failed to check link after 3 attempts"),
			}
		})
	}
	wg.Wait()
	return results, nil
}
