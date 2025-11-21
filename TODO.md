# Project Check Map

## Backend

---

- [x] graceful shutdown
- [x] loading config using viper
- [x] self assign TLS
- [x] server run on HTTPS/2
- [x] init router
- [x] connect to MongoDB
- [x] connect to DragonflyDB
- [x] implement HTTP/2 pusher (push exists in gin)
- [x] implement storing and caching logic
- [x] use zap package to log as json and lamberjack to rotate logs
- [x] implement jwt auth system
- [x] auth middleware based on jwt
- [] add admin only middleware
- [] website routes
- [] student route
- [] teacher route
- [] admin route
- [] website handlers
- [] toggle isActive

## Frontend

---

- [] design the Frontend
- [] implement the design
- [] implement AJAX via HTMX
- [] connect the frontend with the backend
- [] strip the code using templ & compile
- [] connect grafana to HTMX dashboard via HTMX WS extention
- [] local storage

## Deploy

---

- [] build the dockerfile
- [] make the nix config
- [] implement grafana and prometheus
