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
- [] implement storing and caching logic
- [] use zap package to log as json and lamberjack to rotate logs
- [] implement jwt auth system
- [] auth middleware based on jwt
- [] website routes
- [] website handlers

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
