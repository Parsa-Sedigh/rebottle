# Everyone can recycle! ♻️

Todo:

- [x] refresh token
- [x] graceful shutdown
- [x] validation
- [x] dockerfiles
- [ ] more secure static fileserver
- [ ] create a handler for verifying user's email when he is already logged in(VerifyUserEmail handler)
- [ ] best practice for initializing dbs using a sql file when having `docker-entrypoint-initdb.d` in `volume` section of docker?
- [x] common format for validation error responses
- [x] use the zap logger
- [ ] separate package for http handlers
- [ ] production-grade Makefile
- [x] when user signs up but doesn't verify with OTP and come back later, we should generate an OTP for him to verify 
- [ ] how avoid being known as spammer when sending email and SMS? 
- [ ] otp through sms using kavehnegar
- [ ] which service for mailing?
- [ ] microservices
- [ ] queue(rabbit or NATS?)
- [ ] github actions
- [ ] observability (Prometheus)
- [ ] fully document the api using openapi
- [ ] deployment using K8S

---

## architecture

### DTO and datastruct(model)
Transfer req payload into DTO struct, put it in service layer, in service transform dto into datasturct and pass it to repository.

We use DTO to transfer data from handlers to services and we use datastruct to work with repository level.

So repository level works only with datastructs and service and service and handlers work with DTOs.

- application level => handlers
- service level => business logic
- repository level => we work with database
