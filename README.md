# stop-checker.com

Hello World

go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate

models:
  Stop:
    model: "stop-checker.com/db/model.Stop"
  model.StopRoute:
    model: "stop-checker.com/db/model.model.StopRoute"
  Route:
    model: "stop-checker.com/db/model.Route"
  Trip:
    model: "stop-checker.com/db/model.Trip"
  Service:
    model: "stop-checker.com/db/model.Service"
  ServiceException:
    model: "stop-checker.com/db/model.ServiceException"
  Location:
    model: "stop-checker.com/db/model.Location"
