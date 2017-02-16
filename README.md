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
