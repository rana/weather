# Weather

An HTTP server that serves a weather forecast.

### Features
* One GET endpoint
* Accepts latitude and longitude as query parameters
* Returns a short forecast in an area for today (“Partly Cloudy” etc)
* Returns a temperature characterization as “hot”, “cold”, or “moderate”
* Uses the [National Weather Service API](https://www.weather.gov/documentation/services-web-api) for data

### Running the Service

Start the server with:
```
go run main.go
```

The server will start on port 8080.

### Testing in Browser

Open your web browser and navigate to:
```
http://localhost:8080/?lat=41.837&lon=-87.685
```

### Example run

```sh
> go run main.go
2025/05/12 17:35:20 External weather request (lat: 41.8370, lon: -87.6850)
```

![Local weather API call](/browser.png "Local weather API call")
