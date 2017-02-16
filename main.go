package main

import (
	"flag"
	"net/http"

	"github.com/olorin/nagiosplugin"

	pagespeedonline "google.golang.org/api/pagespeedonline/v2"
)

type ApiKey struct {
	Key string
}

func (a ApiKey) Get() (key, value string) {
	return "key", a.Key
}

func main() {
	check := nagiosplugin.NewCheck()
	defer check.Finish()

	check.AddResult(nagiosplugin.OK, "Speed is OK")

	api := flag.String("apikey", "", "API key")
	mw := flag.Int("mw", 0, "Mobile score warning")
	mc := flag.Int("mc", 0, "Mobile score critical")
	dw := flag.Int("dw", 0, "Desktop score warning")
	dc := flag.Int("dc", 0, "Desktop score critical")
	url := flag.String("url", "", "URL to check")
	thirdparty := flag.Bool("thirdparty", false, "Include third party resources")

	flag.Parse()

	if *api == "" {
		check.Unknownf("API key missing (-apikey=)")
	}
	apikey := ApiKey{Key: *api}

	if *mc > *mw {
		check.Unknownf("MC %d can't be greater than MW %d", *mc, *mw)
	}

	if *dc > *dw {
		check.Unknownf("DC %d can't be greater than DW %d", *dc, *dw)
	}

	if *url == "" {
		check.Unknownf("URL missing (-url=)")
	}

	svc, err := pagespeedonline.New(&http.Client{})
	if err != nil {
		check.Unknownf("Unable to create PageSpeedOnline service: %v", err)
	}

	pagespeed := svc.Pagespeedapi.Runpagespeed(*url)
	pagespeed.FilterThirdPartyResources(!*thirdparty)

	result, err := pagespeed.Do(apikey)
	if err != nil {
		check.Unknownf("PageSpeedOnline service: %v", err)
	}

	if result.ResponseCode != 200 {
		check.AddResultf(nagiosplugin.WARNING, "Status code is %d", result.ResponseCode)
	}

	check.AddPerfDatum("NumberResources", "", float64(result.PageStats.NumberResources))
	check.AddPerfDatum("NumberHosts", "", float64(result.PageStats.NumberHosts))
	check.AddPerfDatum("TotalRequestBytes", "B", float64(result.PageStats.TotalRequestBytes))
	check.AddPerfDatum("NumberStaticResources", "", float64(result.PageStats.NumberStaticResources))
	check.AddPerfDatum("HtmlResponseBytes", "B", float64(result.PageStats.HtmlResponseBytes))
	check.AddPerfDatum("TextResponseBytes", "B", float64(result.PageStats.TextResponseBytes))
	check.AddPerfDatum("CssResponseBytes", "B", float64(result.PageStats.CssResponseBytes))
	check.AddPerfDatum("ImageResponseBytes", "B", float64(result.PageStats.ImageResponseBytes))
	check.AddPerfDatum("JavascriptResponseBytes", "B", float64(result.PageStats.JavascriptResponseBytes))
	check.AddPerfDatum("FlashResponseBytes", "B", float64(result.PageStats.FlashResponseBytes))
	check.AddPerfDatum("OtherResponseBytes", "B", float64(result.PageStats.OtherResponseBytes))
	check.AddPerfDatum("NumberJsResources", "", float64(result.PageStats.NumberJsResources))
	check.AddPerfDatum("NumberCssResources", "", float64(result.PageStats.NumberCssResources))

	desktopScore := int(result.RuleGroups["SPEED"].Score)
	check.AddPerfDatum("DesktopScore", "%", float64(desktopScore))

	pagespeed.Strategy("mobile")
	result, err = pagespeed.Do(apikey)
	if err != nil {
		check.Unknownf("PageSpeedOnline service: %v", err)
	}

	mobileScore := int(result.RuleGroups["SPEED"].Score)
	check.AddPerfDatum("MobileScore", "%", float64(mobileScore))
	check.AddPerfDatum("MobileUsability", "%", float64(result.RuleGroups["USABILITY"].Score))

	r := nagiosplugin.OK
	if desktopScore < *dw {
		if desktopScore < *dc {
			r = nagiosplugin.CRITICAL
		} else {
			r = nagiosplugin.WARNING
		}
		check.AddResultf(r, "Desktop score is %d", desktopScore)
	}

	if mobileScore < *mw {
		if mobileScore < *mc {
			r = nagiosplugin.CRITICAL
		} else {
			r = nagiosplugin.WARNING
		}
		check.AddResultf(r, "Mobile score is %d", mobileScore)
	}

}
