# Nyx

![Unit Tests](https://github.com/prateekkhenedcodes/Nyx/actions/workflows/ci.yml/badge.svg)

Nyx is HTTP webserver that is written from scratch in Golang where user need not to give any data and communicate securely.

# Requirements

```go1.24.1``` or greater

# Installation and setup


* clone the project with git clone ```https://github.com/prateekkhenedcodes/Nyx.git```
* Install project dependencies with ```go mod tidy``` in the root of the project


# Endpoints

``` go
    	mux.Handle("/app/", apiCfg.MiddleWareMetrics(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountHits)
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("POST /api/register", apiCfg.Register)
	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/token/refresh", apiCfg.RefreshToken)
	mux.HandleFunc("POST /api/token/revoke", apiCfg.RevokeToken)
	mux.HandleFunc("POST /api/logout", apiCfg.RevokeToken)
	mux.HandleFunc("POST /api/nyx-servers", apiCfg.CreateNyxServer)
	mux.HandleFunc("/api/nyx-servers/join", apiCfg.JoinNyxServer)
	mux.HandleFunc("POST /api/nyx-servers/disconnect", apiCfg.DisconnectNyxServer)
```

##### ``/app/`` 

* ``/app/`` is where the homepage of the Nyx is severed HTML, CSS and Javascript 

#####  ``/api/healthz`` 

* ``/api/healthz`` GET endpoint which gives code 200 if server is fine

##### ``/admin/metrics`` 

* ``/admin/metrics`` GET endpoint gives the number of times the home page has been hit, it resets on server restart

##### ``/api/register``

* ``/api/register`` end-point where user can register to the service its just one click, doesn't take any of your personal data

##### ``/admin/reset``

* ``/admin/reset`` endpoint where only admin can reset the database and the number of hits main page to Zero

##### ``/api/login``

* ``/api/login`` endpoint where user can login by entering id and the nyx code 

##### ``/api/token/refresh``

* ``/api/token/refresh`` endpoint where client can refresh their access token using their refresh token

##### ``/api/token/revoke``

* ``/api/token/revoke`` is endpoint where client can revoke their refresh token 

##### ``/api/logout``

* ``/api/logout`` is self-explanatory

##### ``/api/nyx-servers`` 

* ``/api/nyx-servers`` is endpoint where user can create credintials for the nys-server (chat room) , using them you can enter the nyx-server.

##### ``/api/nyx-servers/join``

* ``/api/nyx-servers/join`` end-point used to join the nyx-server using credentials from the previous end-point.

##### ``/api/nyx-servers/disconnect``

* ``/api/nyx-servers/disconnect`` end-point use to disconnect from the particular server with specifing the serverId in the request body.