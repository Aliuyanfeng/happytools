# Get objects related to a user
## 来源
https://docs.virustotal.com/reference/users-relationships
## 代码参考
```golang
package main

import (
	"fmt"
	"net/http"
	"io"
)

func main() {

	url := "https://www.virustotal.com/api/v3/users/2fbdb6b98badd2ca351e76f371e593b4ef281ea58b2a4c63fdc7caecad28be10/api_quota_group?limit=10"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", "2fbdb6b98badd2ca351e76f371e593b4ef281ea58b2a4c63fdc7caecad28be10")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
```
## 返回值参考
```json
{
  "links": {
    "self": "https://www.virustotal.com/api/v3/users/yanminghao/api_quota_group"
  },
  "meta": {
    "count": 1
  },
  "data": {
    "id": "antiy",
    "type": "group",
    "links": {
      "self": "https://www.virustotal.com/api/v3/groups/antiy"
    },
    "attributes": {
      "agreement_company_country": "CN",
      "lock_users_api_quota_group": true,
      "display_flat_rate_message": false,
      "force_users_saml_login": false,
      "agreement_signatory_title": "Antiy",
      "group_token": "YW50aXl8fHYzfHwxNzQxODY3MjYzfHw4MTZjYWFhOGQwMTk5ZjkxMjViYjAwODBlNDI4OTUyYzVjZmIwYWZlYTRkYmZiZTRkMGEwNzk4YjJkYmY0MDI0",
      "agreement_company_address": "Haerbin",
      "country_iso": "CN",
      "preferences": {
        "collection": {
          "default_visibility": "public",
          "block_public_collections": false
        },
        "monitor_partner": {
          "engines": [
            "Antiy-AVL"
          ]
        }
      },
      "agreement_signed_date": 1516067760,
      "gti_gcp_project": "",
      "country": "China",
      "force_users_oidc_login": false,
      "agreement_company_name": "shenchangwei",
      "agreement_signatory": "shenchangwei",
      "industry": "Antivirus",
      "organization": "Antiy",
      "domain_name": "antiy.cn",
      "workforce_identity_pool_name": "",
      "workforce_identity_pool_provider_name": "",
      "organization_legal_name": "Harbin Antiy Technology Co.,Ltd",
      "privileges": {
        "activity-logs": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "big-files": {
          "granted": true
        },
        "click_to_accept": {
          "granted": true
        },
        "creditcards": {
          "granted": true
        },
        "dogfooder": {
          "granted": false
        },
        "domain-feed": {
          "granted": false
        },
        "file-behaviour-feed": {
          "granted": false
        },
        "distribution-download": {
          "granted": true
        },
        "distribution": {
          "granted": true
        },
        "downloads-tier-1": {
          "granted": false
        },
        "downloads-tier-2": {
          "granted": true
        },
        "file-feed": {
          "granted": true
        },
        "full-distribution": {
          "granted": true
        },
        "intelligence-search-tier-1": {
          "granted": false
        },
        "intelligence-search-tier-2": {
          "granted": false
        },
        "intelligence-search-tier-3": {
          "granted": false
        },
        "intelligence": {
          "granted": true
        },
        "ioc-stream": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "ip-feed": {
          "granted": false
        },
        "livehunt": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "monitor-partner": {
          "granted": true
        },
        "monitor": {
          "granted": false
        },
        "oem_click_to_accept": {
          "granted": false
        },
        "private": {
          "granted": true
        },
        "private-collections": {
          "granted": false
        },
        "retrohunt-tier-1": {
          "granted": false
        },
        "retrohunt-tier-2": {
          "granted": false
        },
        "retrohunt-tier-3": {
          "granted": false
        },
        "retrohunt": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "sales-staff": {
          "granted": false
        },
        "saved-search": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "staff": {
          "granted": false
        },
        "url-distribution": {
          "granted": true
        },
        "url-feed": {
          "granted": true
        },
        "vtalerts": {
          "granted": false
        },
        "vtdiff-api": {
          "granted": false
        },
        "vtdiff-ui": {
          "granted": true,
          "inherited_from": "intelligence",
          "inherited_via": "privilege"
        },
        "vtinsights": {
          "granted": false
        }
      },
      "agreement_signatory_email": "scw@antiy.cn",
      "quotas": {
        "api_requests_hourly": {
          "allowed": 600000,
          "used": 20185
        },
        "api_requests_daily": {
          "allowed": 2000000,
          "used": 165822
        },
        "api_requests_monthly": {
          "allowed": 60000000,
          "used": 15091996
        },
        "intelligence_downloads_monthly": {
          "allowed": 200000,
          "used": 476
        },
        "intelligence_searches_monthly": {
          "allowed": 200000,
          "used": 1548
        },
        "intelligence_retrohunt_jobs_monthly": {
          "allowed": 5,
          "used": 0
        },
        "intelligence_hunting_rules": {
          "allowed": 25,
          "used": 241
        },
        "intelligence_graphs_private": {
          "allowed": 0,
          "used": 0
        },
        "intelligence_vtdiff_creation_monthly": {
          "allowed": 10000,
          "used": 0
        },
        "monitor_storage_bytes": {
          "allowed": 0,
          "used": 0
        },
        "monitor_storage_files": {
          "allowed": 0,
          "used": 0
        },
        "monitor_uploaded_bytes": {
          "allowed": 0,
          "used": 0
        },
        "monitor_uploaded_files": {
          "allowed": 0,
          "used": 0
        },
        "collections_creation_monthly": {
          "allowed": 0,
          "used": 0
        },
        "private_scans_monthly": {
          "allowed": 0,
          "used": 0
        },
        "private_scans_per_minute": {
          "allowed": 30,
          "used": 0
        },
        "private_urlscans_monthly": {
          "allowed": 0,
          "used": 0
        },
        "private_urlscans_per_minute": {
          "allowed": 30,
          "used": 0
        }
      }
    },
    "context_attributes": {
      "role": "user"
    }
  }
}
```