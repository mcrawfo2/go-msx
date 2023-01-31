# The Idempotency-Key HTTP Header Field

API caller can provide an optional unique cache id header (Idempotency-Key) to cache the backend response (for a default of 5 minutes).

Cache id is recommended to be of the form {microservice}-{api-operation-id} for APIs that are site wide eg (lan-someoperation). If the API response is specific to tenants then affix the tenant id (eg lan-someoperation-tenantid).

Sample Api call below. The first time the backend encounters the cache id header it will save the (valid) response in an lru cache. Succeeding response would retrieve data from the cache and not hit the controllers.

curl -X GET "http://10.81.85.127:8765/lan/api/v1/tags?tenantId=836fc9a3-a440-4d15-ada1-4927e0be47d3&entityType=NETWORK&page=0&pageSize=100" -H "accept: application/json" -H "Authorization: Bearer $jwt" -H "Idempotency-Key: lan-listTags-836fc9a3-a440-4d15-ada1-4927e0be47d3" 

For more information about Idempotency-Key HTTP Header Field please visit https://greenbytes.de/tech/webdav/draft-ietf-httpapi-idempotency-key-header-latest.html


## Requiring or recommending idempotency from the APIs 

```
	return svc.GET("/").
		Operation("listTags").
		Doc("List Tags").
		Notes("API to get tags by tenant and entity type").
		Param(paramQueryTenantId).
		Param(paramQueryEntityType).
		Param(paramQuerySearchPattern).
		Param(paramQueryPage).
		Param(paramQueryPageSize).
		Do(webservice.PopulateParams(new(params))).
		Do(webservice.ValidateParams(paramsValidator)).

		// Do(idempotency.ValidateIdempotency(idempotency.ValidateIdempotencyRecommend)).
		Do(idempotency.ValidateIdempotency(idempotency.ValidateIdempotencyRequire)).

		Do(webservice.StandardList).
		Do(viewPermission).
		To(webservice.Controller(
			func(req *restful.Request) (body interface{}, err error) {
				converted, ok := webservice.Params(req).(*params)
				if !ok {
					return nil, errors.New("failed to parse param")
				}
				var pagingRequest = paging.NewRequestFromQuery(converted.Page, converted.PageSize)

				res, err := c.service.Search(req.Request.Context(), converted.TenantId, converted.EntityType, converted.SearchPattern, pagingRequest)
				if err != nil {
					return nil, err
				}

				return res, nil
			}))
}
```

## Sample idempotency recommended or required response



curl -v -X GET "http://10.81.85.127:8765/lan/api/v1/tags?tenantId=$tenantId&entityType=NETWORK&page=0&pageSize=100" -H "accept: application/json" -H "Authorization: Bearer $jwt"

### Regular non idempotency

```
{
"command": "listTags",
"httpStatus": "OK",
"message": "",
"params": {},
"responseObject": {
"page": 0,
"pageSize": 100,
"totalItems": 0,
"hasNext": false,
"hasPrevious": false,
"contents": null
},
"success": true
}
```

### Recommend idempotency 
Do(idempotency.ValidateIdempotency(idempotency.ValidateIdempotencyRecommend))

```
< Go-Msx-Idempotency: This operation is idempotent and it requires correct usage of Idempotency Key.
< Link: https://cto-github.cisco.com/NFV-BU/go-msx/blob/main/idempotency.md
<
{
"command": "listTags",
"httpStatus": "OK",
"message": "",
"params": {},
"responseObject": {
"page": 0,
"pageSize": 100,
"totalItems": 0,
"hasNext": false,
"hasPrevious": false,
"contents": null
},
"success": true
}
```

### Require idempotency 
Do(idempotency.ValidateIdempotency(idempotency.ValidateIdempotencyRequire))

```
< HTTP/1.1 400 Bad Request
<
{
"type": "https://cto-github.cisco.com/NFV-BU/go-msx/blob/main/idempotency.md",
"title": "Idempotency-Key is missing",
"detail": "This operation is idempotent and it requires correct usage of Idempotency Key."
}
```