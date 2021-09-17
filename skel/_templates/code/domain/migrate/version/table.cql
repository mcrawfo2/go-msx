CREATE TABLE lower_snake_singular (
    lower_snake_singular_id uuid PRIMARY KEY,
--#if TENANT_DOMAIN
    tenant_id uuid,
--#endif TENANT_DOMAIN
    data text
);
--#if TENANT_DOMAIN
CREATE INDEX ON lower_snake_singular(tenant_id);
--#endif TENANT_DOMAIN
