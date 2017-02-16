# check_google_pagespeed
Icinga/Nagios Google Pagespeed Insights

Usage:

```
./check_google_pagespeed -apikey=XXX -url=http://www.example.org/ -thirdparty -mw=50 -mc=40 -dw=70 -dc=60
```

Where

```
-apikey is Google API key with enabled pagespeed service
-thirdparty Includes third party resources
-mc -mw is mobile warning/critical score
-dc -dw is desktop warning/critical score
```
Sample output

```
check_google_pagespeed -apikey=XXX -url=https://www.cnn.com
OK: Speed is OK | NumberResources=183;;;; NumberHosts=75;;;; TotalRequestBytes=52648B;;;; NumberStaticResources=87;;;;
HtmlResponseBytes=246538B;;;; TextResponseBytes=580B;;;; CssResponseBytes=1541293B;;;; ImageResponseBytes=756803B;;;;
JavascriptResponseBytes=4520858B;;;; FlashResponseBytes=0B;;;; OtherResponseBytes=235475B;;;; NumberJsResources=71;;;;
NumberCssResources=2;;;; DesktopScore=70%;;;; MobileScore=62%;;;; MobileUsability=95%;;;;
```
