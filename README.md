# Nyx


##### Nyx is api that is written in Golang where user need not to give any data and communicate securely 

# Endpoints

``` go
    mux.Handle("/app/", apiCfg.MiddleWareMetrics(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountHits)
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
```

##### ``/app/`` 

###### ``/app/`` is where the homepage of the Nyx is severed HTML, CSS and Javascript 

#####  ``/api/healthz`` 

###### ``/api/healthz`` GET endpoint which gives code 200 if server is fine

##### ``/admin/metrics`` 

###### ``/admin/metrics`` GET endpoint gives the number of times the home page has been hit, it resets on server restart