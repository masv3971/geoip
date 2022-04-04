# GeoIP

## API
### POST /login_event
#### Request
```json
	"data": {
		"eppn_hashed": <string>,
		"client_ip": <string>,
		"device_id": <string>,
		"user_agent": {
			"browser": {
				"family": <string>,
				"version": [
					<int>,
					<int>,
					<int>
				],
				"version_string": <string>
			},
			"os": {
				"family": <string>,
				"version": [
					<int>,
					<int>,
					<int>
				]
			},
			"sophisticated": {
				"is_mobile": <boolean>,
				"is_tablet": <boolean>,
				"is_pc": <boolean>,
				"is_touch_capable": <boolean>,
				"is_bot": <boolean>
			}
		}
	}
}
```


#### Reply -- good
```json
{
	"data": {
		"status": true
	},
	"error": null
}
```

### GET /stats/overview
#### Reply
```
{
	"data": {
		"overview": [
			{
				"eppn_hashed": <string>,
				"number_of_login_events": <int>,
				"number_of_countries": <int>,
				"number_of_unique_countries": <int>,
				"number_of_ips": <int>,
				"number_of_unique_ips": <int>
			},
            }
            .
            .
            .
            }
		]
	},
	"error": null
}
```


### GET /stats/eppn/{eppn}/long
#### Reply
```
{
	"data": {
		"LoginEvents": [
			{
				"uid": <string>,
				"eppn_hashed": <string>,
				"hash": <string>,
				"timestamp": <string>,
				"timestamp_ml": <int>,
				"ip": {
					"ip_addr": <string>,
					"ip_addr_ML": <string>,
					"asn": {
						"number": <int>,
						"organization": <string>
					},
					"isp": null,
					"anonymous_ip": null
				},
				"device_id_hashed": <string>,
				"user_agent": {
					"browser": {
						"family": <string>,
						"family_ml": <int>,
						"version": [
							<int>,
							<int>,
							<int>
						],
						"version_string": <string>
					},
					"os": {
						"family": <string>,
						"family_ml": <int>,
						"version": [
							<int>,
							<int>,
							<int>
						],
						"version_string": <string>
					},
					"device": {
						"family": <string>,
						"family_ml": <int>,
						"brand": <string>,
						"model": <string>
					},
					"sophisticated": {
						"is_mobile": <boolean>,
						"is_mobile_ml": <int>,
						"is_tablet": <boolean>,
						"is_tablet_ml": <int>,
						"is_pc": <boolean>,
						"is_pc_ml": <int>,
						"is_touch_capable": <boolean>,
						"is_touch_capable_ml": <int>,
						"is_bot": <boolean>,
						"is_bot_ml": <int>
					}
				},
				"phishing": null,
				"location": {
					"coordinates": {
						"latitude": <float64>,
						"longitude": <float64>
					},
					"city": <string>,
					"city_ml": <int>,
					"country": <string>,
					"country_ml": <int>,
					"continent": <string>,
					"continents_ml": <int>
				},
				"login_method": "",
				"fraudulent": <boolean>,
				"fraudulent_ml": <int>
			},
            {
            .
            .
            .
            }
		]
	},
	"error": null
}
```

### GET /stats/eppn/{eppn}/specific
#### Reply
```
{
	"data": {
		"StatsData": {
			"ip": {
				"number_of_elements": <int>,
				"entropy": <float64>,
				"standardDeviation": <float64>
			},
			"user_agent_device": {
				"number_of_elements": <int>,
				"entropy": <float64>,
				"standardDeviation": <float64>
			},
			"user_agent_os": {
				"number_of_elements": <int>,
				"entropy": <float64>,
				"standardDeviation": <float64>
			}
		}
	},
	"error": null
}
