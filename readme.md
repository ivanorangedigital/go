# Dev Server

## Requirements

1. watchexec
2. nc

install them with your favorite package manager and enjoy

## Run server.js
run `npm start` inside `dev-ws-server` to start ws server for dev purposes

## Run make dev
this command will start go server with watcher for go and tmpl files (tmpl included because templates was cached)

## Run make ws
this command will start watcher for static files (css,js,tmpl) and send ping to reload page

## WS client
WS client is inside the file `ui/static/js/dev-watch.js` this file is included in layout base (remove it on production)

### Structure
cmd folder contain applications  
pkg folder contain code non-specific application (reusable in other applications)
