package eureka

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "sync/atomic"
    "github.com/tietang/go-utils/errs"
    "strings"
)

type Discovery struct {
    apps      *Applications
    AppNames  map[string]string
    eurekaUrl []string
    callbacks []func(*Applications)
    ct        uint64
}

func NewDiscovery(eurekaUrl []string) *Discovery {
    return &Discovery{eurekaUrl: eurekaUrl, callbacks: make([]func(*Applications), 0)}
}

func (d *Discovery) AddCallback(callback func(*Applications)) {
    d.callbacks = append(d.callbacks, callback)

}

func (d *Discovery) execCallbacks(apps *Applications) {
    if len(d.callbacks) > 0 {
        for _, c := range d.callbacks {
            go c(apps)
        }
    }
}

func (d *Discovery) ScheduleAtFixedRate(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *Discovery) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}

func (d *Discovery) run() {
    apps, err := d.GetApplications()
    if err == nil || apps != nil {
        d.apps = apps
        d.execCallbacks(apps)
    } else {
        fmt.Println(err)
    }
}

func (c *Discovery) GetApps() *Applications {
    if c.apps == nil {
        apps, err := c.GetApplications()
        if err == nil {
            return apps
        }
    }
    return c.apps
}

func (c *Discovery) GetApp(name string) *Application {
    if c.apps == nil {
        fmt.Println("Applications is nil")
        return nil
    }
    for _, app := range c.apps.Applications {
        if strings.ToLower(app.Name) == strings.ToLower(name) {
            return &app
        }
    }
    return nil
}

func (c *Discovery) GetInstances(name string) (*Application, error) {
    //url := c.eurekaUrl + "/apps"
    url := c.getEurekaServerUrl() + "/apps/" + name

    //	req, err := http.NewRequest("GET", url, nil)
    //	req.Header.Add("Accept", "application/json")
    //	res, err := c.client.Do(req)
    //	http.Client.Do(req)
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    //	fmt.Println(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        fmt.Println(err)
        return nil, err
    }
    app := new(Application)
    err = xml.Unmarshal(respBody, app)

    //	fmt.Println(string(respBody))
    //	fmt.Println(err, applications)
    return app, err
}

func (c *Discovery) GetApplications() (*Applications, error) {
    //url := c.eurekaUrl + "/apps"
    url := c.getEurekaServerUrl() + "/apps"

    //	req, err := http.NewRequest("GET", url, nil)
    //	req.Header.Add("Accept", "application/json")
    //	res, err := c.client.Do(req)
    //	http.Client.Do(req)
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    //	fmt.Println(res.StatusCode)
    respBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        fmt.Println(err)
        return nil, err
    }
    var applications *Applications = new(Applications)
    err = xml.Unmarshal(respBody, applications)

    //	fmt.Println(string(respBody))
    //	fmt.Println(err, applications)
    return applications, err
}

func (c *Discovery) getEurekaServerUrl() string {
    ct := atomic.AddUint64(&c.ct, 1)
    size := len(c.eurekaUrl)
    if size == 0 {
        panic(errs.NilPointError("eureka url is empty"))
    }
    index := int(ct) % size
    url := c.eurekaUrl[index]
    //if strings.LastIndex(url,"/")>-1{
    url = strings.TrimSuffix(url, "/")
    //}
    return url
}
