package eureka

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/props"
	"strconv"
	"testing"
)

func TestNewClient(t *testing.T) {
	conf := props.NewEmptyMapConfigSource("eureka")
	name := "test_demo"
	conf.Set("application.name", name)
	port := "28088"
	portInt, _ := strconv.Atoi(port)
	conf.Set("application.port", port)
	leaseRenewalIntervalInSeconds := "10"
	leaseRenewalIntervalInSecondsInt, _ := strconv.Atoi(leaseRenewalIntervalInSeconds)
	conf.Set("eureka.instance.leaseRenewalIntervalInSeconds", leaseRenewalIntervalInSeconds)
	preferIpAddress := "true"
	conf.Set("eureka.instance.preferIpAddress", preferIpAddress)
	defaultZone := "http://172.16.1.248:8761/eureka/"
	conf.Set("eureka.client.serviceUrl.defaultZone", defaultZone)
	conf.Set("eureka.client.serviceUrl.ali-hz-1", defaultZone)

	c := NewClient(conf)

	Convey("Test New Eureka client config", t, func() {
		So(c, ShouldNotBeNil)
		ins := c.InstanceInfo
		Convey("InstanceInfo", func() {
			So(ins, ShouldNotBeNil)
			So(ins.LeaseInfo, ShouldNotBeNil)
			So(ins.LeaseInfo.RenewalIntervalInSecs, ShouldEqual, leaseRenewalIntervalInSecondsInt)
			So(ins.Port, ShouldNotBeNil)
			So(ins.Port.Enabled, ShouldEqual, true)
			So(ins.Port.Port, ShouldEqual, portInt)
			So(ins.AppName, ShouldEqual, name)

		})
		ic := c.InstanceConfig
		Convey("InstanceConfig", func() {
			So(ic, ShouldNotBeNil)
			So(ic.PreferIpAddress, ShouldEqual, true)
			So(ic.LeaseRenewalIntervalInSeconds, ShouldEqual, uint(leaseRenewalIntervalInSecondsInt))
			So(ic.Appname, ShouldEqual, name)

		})
		cc := c.ClientConfig
		Convey("ClientConfig", func() {
			So(cc, ShouldNotBeNil)
			So(cc.ServiceUrl, ShouldNotBeNil)
			So(cc.ServiceUrl["defaultZone"], ShouldNotBeNil)
			So(cc.ServiceUrl["defaultZone"], ShouldEqual, defaultZone)
			So(cc.ServiceUrl["ali-hz-1"], ShouldNotBeNil)
			So(cc.ServiceUrl["ali-hz-1"], ShouldEqual, defaultZone)

		})

	})
}
