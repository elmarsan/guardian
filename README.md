# Guardian

Http file server protected by jwt written in Go.

## Requirements

### Environment Variables


| Variable Name | Required | Description | Example Value |
|----------------|-------------|---------------|----------|
| `JWT_KEY` | YES | Secret for signing jww tokens | `ef51c9fc4b73b74149f8dd0a0ee5e9aaf605a1cb` |
| `JWT_EXPIRATION_TIME` | NO | Time in minutes for expire jwt policy | `60` |