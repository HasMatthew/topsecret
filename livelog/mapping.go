PUT /indexes
//opens and events should be an array, i am not sure about that part
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
                        "type": "string"
                     },
                     "Location": {
                        "type": "string"
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
                              "type": "string"
                           },
                           "Location": {
                              "type": "string"
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
                        "type": "string"
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
                              "type": "string"
                           },
                           "Location": {
                              "type": "string"
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
                     "Location": {
                        "type": "string"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string"
                     },
                     "StatImpressionId": {
                        "type": "string"
                     },
                     "StatInstallId": {
                        "type": "string"
                     },
                     "StatOpenId": {
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
                        "type": "string"
                     },
                     "Location": {
                        "type": "string"
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
                        "type": "string"
                     },
                     "Location": {
                        "type": "string"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string"
                     },
                     "StatImpressionId": {
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
               "Opens": {
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
                              "type": "string"
                           },
                           "Location": {
                              "type": "string"
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
                        "type": "string"
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
                              "type": "string"
                           },
                           "Location": {
                              "type": "string"
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
                     "Location": {
                        "type": "string"
                     },
                     "PostalCode": {
                        "type": "long"
                     },
                     "RegionCode": {
                        "type": "string"
                     },
                     "StatClickId": {
                        "type": "string"
                     },
                     "StatImpressionId": {
                        "type": "string"
                     },
                     "StatInstallId": {
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
               }
            }
         }
      }
   
}