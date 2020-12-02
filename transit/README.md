# MSX Transit Module

MSX transit is an implementation of transit encryption.  It allows swappable encryption implementations via the `Provider` interface.

## Usage

The primary mode of consumption for MSX Transit is within business Models.  To add transit encryption support to your model, add an anonymous `transit.WithSecureData` member to your model structure:

```go
type Organization struct {
	transit.WithSecureData
	OrganizationId string     `db:"organization_id"`
	TenantId       gocql.UUID `db:"tenant_id"`
}
```

This embedded struct will store its data in a field named `secure_data`, so a migration will need to add such a field if the table already exists:

```sql
ALTER TABLE organization ADD COLUMN secure_data TEXT;
```

To store and retrieve individual encrypted fields from your model, add accessors:

```go
const secureDataMerakiApiKey = "merakiApiKey"

func (o *Organization) MerakiApiKey(ctx context.Context) (string, error) {
	return o.SecureValue(ctx, organizationSecureFieldMerakiApiKey)
}

func (o *Organization) SetMerakiApiKey(ctx context.Context, value *string) error {
	return o.SetSecureValue(ctx, o.TenantId.Bytes(), secureDataMerakiApiKey, value)
}
```

You can then use these accessors in your converter and services to retrieve and store the values from your model.

### Per-Application Encryption

Sometimes you will want values to be encrypted, but non on a per-tenant basis.  In this case, define your key id within your domain package, and use it in place of the TenantId in your property setters:

```go
var appKeyId := types.MustParseUUID("3e246fc7-12d8-4626-a739-1fd22bbf47f0")

func (o *Organization) SetMerakiApiKey(ctx context.Context, value *string) error {
	return o.SetSecureValue(ctx, appKeyId.Bytes(), secureDataMerakiApiKey, value)
}
```

