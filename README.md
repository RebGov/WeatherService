# Weather Service
Weather Service API in GO

API backend service that uses latidude and longitude to get the weather conditions for that location. 

## How To Use

Docker Container that hosts Weather Service API.

Enter your latitude/longitude in the request body and receive back the weather condition in given location.

### RUN
``` WEATHER_ID={use-your-value} docker compose up --build```

NOTE: `WEATHER_ID` is the app_id required by Open Weather Map. (See [Dependencies](https://github.com/RebGov/WeatherService/blob/feature-service-create2/README.md#dependencies))

### USE
API takes in the following attributes and returns the weather condition in the given location.
- Latitude: float32,
- Longitude: float32,

#### Endpoints:
- http://localhost:8001/weather
   
#### JSON Request Body:
```.json
{
    "Latitude": 32.777981,
    "Longitude": -96.796211,
}
```
#### CURL Command
```
curl --location --request GET 'http://localhost:8001/weather' \
--header 'Content-Type: application/json' \
--data '{
    "Latitude":32.777981,
    "Longitude":-96.796211
}'
```
## Swagger


## Helpful pages:
 - Need to get a latitude/Longitude for your area visit https://www.latlong.net/.

 ## Dependencies
 - Weather Service currently utilizes Open Weather Map for obtaining the weather condition at given latitude/longitude.
    - This requires our service to utilize an app_id key. Please be sure to start the docker container utilizing your own app_id key.
    - If you do not have a key visit: `https://openweathermap.org/`
- Go version 1.22
- Docker

## Contact
- `bg.luv2code@gmail.com`
## Last Updated
- 2024-07-03



© Becci Govert 2024-07-03