# Everyone can recycle! ♻️

Todo:

- [ ] add name, address, age and ... for driver table
- [x] refresh token
- [x] graceful shutdown
- [x] validation
- [x] dockerfiles
- [ ] create a handler for verifying user's email when he is already logged in(VerifyUserEmail handler)
- [ ] best practice for initializing dbs using a sql file when having `docker-entrypoint-initdb.d` in `volume` section of docker?
- [ ] common format for validation error responses
- [ ] use the zap logger
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