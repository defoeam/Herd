# kvs

Golang Kev Value Storage

Enforce TDD

## Backlog
- [x] Initialize our go.mod and main.go files
- [x] Add linting actions for Go standardization(GitHub workflows)
- [x] Protect main branch and require pull requests to main
- [x] Create a Dockerfile for our REST service
- [x] [Initialize our kvs][]
- [ ] Add persistence logic to store on file(docker volume--with docker compose)
- [ ] Write our persistence layer as an interface
- [ ] Create Docker Compose file to reliably create local environment
- [ ] Add standard dictionary methods to our base key value store and API, see the common python dictionary methods.
- [ ] Write unit tests for key value system
- [ ] HTTP tests for REST API
- [ ] Write manual developer setup documentation
- [ ] Setup Microsoft docker dev container for easy developer onboarding

### Stretch Opportunities
- [ ] Add Docker build & test CI
- [ ] Add transaction logger logic to recover from unreliable interruption(see
  Louis)
- [ ] Integrate [Go Gin Swagger Docs][]

[Docker Compose]: https://docs.docker.com/compose/
[Go Gin Swagger Docs]: https://medium.com/@kumar16.pawan/integrating-swagger-with-gin-framework-in-go-f8d4883f4833
[Initialize our kvs]: (https://medium.com/@anshurai8991/building-a-simple-key-value-store-in-go-adfbd781f16e)
[interface]: https://gobyexample.com/interfaces
