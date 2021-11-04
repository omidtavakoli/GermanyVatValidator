## German Vat Validator service

This project has been developed using Hexagonal architecture and repository pattern.


## How to run?

``docker-compose up -d``
Call this endpoint 

``curl localhost:8082/validator/vat/[VAT_ID]``

## Attention
This service uses external service for vat validation, so it might not build because of test cases has been in build time to check external resource.
In this case please try a bit later

## Future plan:
- [ ] add monitoring service
- [ ] add cache
- [ ] distribute to multi nodes