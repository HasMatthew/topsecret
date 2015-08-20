PUT /indexes
{
   
      "mappings": {
         "types": {
            "properties": {
               "Click": {
                  "properties": {
                     "CountryCode": {
                        "type": "string"
                     },
                     "Created": {
                        "type": "date",
                        "format": "dateOptionalTime"
                     },
                     "DeviceIp": {
                        "type": "string"
                     },
                     "Id": {
                        "type": "string",
 			 "index": "not_analyzed"
                     },
                     "Location": {
                         "type": "geo_point"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "WurflBrandName": {
                        "type": "string"
                     },
                     "WurflDeviceOs": {
                        "type": "string"
                     },
                     "WurflModelName": {
                        "type": "string"
                     }
                  }
               },
               "Common": {
                  "properties": {
                     "AdNetworkId": {
                        "type": "long"
                     },
                     "AdvertiserId": {
                        "type": "long"
                     },
                     "AgencyId": {
                        "type": "long"
                     },
                     "CampaignId": {
                        "type": "long"
                     },
                     "CurrencyCode": {
                        "type": "string"
                     },
                     "GoogleAid": {
                        "type": "string"
                     },
                     "IosIfa": {
                        "type": "string"
                     },
                     "Language": {
                        "type": "string"
                     },
                     "PackageName": {
                        "type": "string"
                     },
                     "PublisherId": {
                        "type": "long"
                     },
                     "PublisherUserId": {
                        "type": "string"
                     },
                     "SiteId": {
                        "type": "long"
                     },
                     "WindowsAid": {
                        "type": "string"
                     }
                  }
               },
               "Events": {
                  "properties": {
                     "CountryCode": {
                        "type": "string"
                     },
                     "Created": {
                        "type": "date",
                        "format": "dateOptionalTime"
                     },
                     "DeviceIp": {
                        "type": "string"
                     },
                     "Id": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "Location": {
                         "type": "geo_point"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatImpressionId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatInstallId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatOpenId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "WurflBrandName": {
                        "type": "string"
                     },
                     "WurflDeviceOs": {
                        "type": "string"
                     },
                     "WurflModelName": {
                        "type": "string"
                     }
                  }
               },
               "Impression": {
                  "properties": {
                     "CountryCode": {
                        "type": "string"
                     },
                     "Created": {
                        "type": "date",
                        "format": "dateOptionalTime"
                     },
                     "DeviceIp": {
                        "type": "string"
                     },
                     "Id": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "Location": {
                         "type": "geo_point"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "WurflBrandName": {
                        "type": "string"
                     },
                     "WurflDeviceOs": {
                        "type": "string"
                     },
                     "WurflModelName": {
                        "type": "string"
                     }
                  }
               },
               "Install": {
                  "properties": {
                     "CountryCode": {
                        "type": "string"
                     },
                     "Created": {
                        "type": "date",
                        "format": "dateOptionalTime"
                     },
                     "DeviceIp": {
                        "type": "string"
                     },
                     "Id": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "Location": {
                        "type": "geo_point"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatImpressionId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "WurflBrandName": {
                        "type": "string"
                     },
                     "WurflDeviceOs": {
                        "type": "string"
                     },
                     "WurflModelName": {
                        "type": "string"
                     }
                  }
               },
               "Opens": {
                  "properties": {
                     "CountryCode": {
                        "type": "string"
                     },
                     "Created": {
                        "type": "date",
                        "format": "dateOptionalTime"
                     },
                     "DeviceIp": {
                        "type": "string"
                     },
                     "Id": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "Location": {
                        "type": "geo_point"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatImpressionId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "StatInstallId": {
                        "type": "string",
  "index": "not_analyzed"
                     },
                     "WurflBrandName": {
                        "type": "string"
                     },
                     "WurflDeviceOs": {
                        "type": "string"
                     },
                     "WurflModelName": {
                        "type": "string"
                     }
                  }
               }
            }
         }
      }
   
}