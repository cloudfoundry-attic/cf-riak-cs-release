package test_support

import()


type FakeCfClient struct {
}

func(cf *FakeCfClient) GetSpaces() string {
	return `{
		"total_results": 2,
		"total_pages": 1,
		"prev_url": null,
		"next_url": null,
		"resources": [
			{
				"metadata": {
					"guid": "413c4df3-66b6-4a7e-a681-bd7f89cffcd9",
					"url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9",
					"created_at": "2014-07-14T21:36:40+00:00",
					"updated_at": null
				},
				"entity": {
					"name": "console",
					"organization_guid": "f29c953f-2c1b-4f62-af32-7543cdb4d01e",
					"organization_url": "/v2/organizations/f29c953f-2c1b-4f62-af32-7543cdb4d01e",
					"developers_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/developers",
					"managers_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/managers",
					"auditors_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/auditors",
					"apps_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/apps",
					"routes_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/routes",
					"domains_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/domains",
					"service_instances_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/service_instances",
					"app_events_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/app_events",
					"events_url": "/v2/spaces/413c4df3-66b6-4a7e-a681-bd7f89cffcd9/events"
				}
			},
			{
				"metadata": {
					"guid": "0d8ff79a-b9b1-4d84-9de1-015a6c884269",
					"url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269",
					"created_at": "2014-07-16T18:20:59+00:00",
					"updated_at": null
				},
				"entity": {
					"name": "riak-test",
					"organization_guid": "1488c115-382e-4e56-96dc-d422f92f108e",
					"organization_url": "/v2/organizations/1488c115-382e-4e56-96dc-d422f92f108e",
					"developers_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/developers",
					"managers_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/managers",
					"auditors_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/auditors",
					"apps_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/apps",
					"routes_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/routes",
					"domains_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/domains",
					"service_instances_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/service_instances",
					"app_events_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/app_events",
					"events_url": "/v2/spaces/0d8ff79a-b9b1-4d84-9de1-015a6c884269/events"
				}
			}
		]
	}`
}
