The Idempotency-Key HTTP Header Field

API caller can provide an optional unique cache id header (Idempotency-Key) to cache the backend response (for a default of 5 minutes).

Cache id is recommended to be of the form {microservice}-{api-operation-id} for APIs that are site wide eg (lan-someoperation). If the API response is specific to tenants then affix the tenant id (eg lan-someoperation-tenantid).

Sample Api call below. The first time the backend encounters the cache id header it will save the (valid) response in an lru cache. Succeeding response would retrieve data from the cache and not hit the controllers.

curl -X GET "http://10.81.85.127:8765/lan/api/v1/tags?tenantId=836fc9a3-a440-4d15-ada1-4927e0be47d3&entityType=NETWORK&page=0&pageSize=100" -H "accept: application/json" -H "Authorization: Bearer $jwt" -H "Idempotency-Key: lan-listTags-836fc9a3-a440-4d15-ada1-4927e0be47d3" 

For more information about Idempotency-Key HTTP Header Field please visit https://greenbytes.de/tech/webdav/draft-ietf-httpapi-idempotency-key-header-latest.html