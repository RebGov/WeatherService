# Weather Service
Weather Service API in GO

API backend service that uses latidude and longitude to get the weather conditions for that location. 

## How To Use

Docker Container that hosts Weather Service API.

Enter your decimal (DD) latitude/longitude in the request body and receive back the weather condition in given location.

### RUN
``` WEATHER_ID={use-your-value} docker compose up --build```

NOTE: `WEATHER_ID` is the app_id required by Open Weather Map. (See [Dependencies](https://github.com/RebGov/WeatherService/blob/feature-service-create2/README.md#dependencies))

### USE
API takes in the following attributes and returns the weather condition in the given location.
- Latitude: float64,
- Longitude: float64,

#### Endpoints:
- http://localhost:8001/weather/get
   
#### JSON Request Body:
```.json
{
    "Latitude": 32.777981,
    "Longitude": -96.796211,
}
```
#### CURL Command
If you change the PORT be sure to upate port in following:
```
curl --location --request GET 'http://localhost:8001/weather/get' \ 
--header 'Content-Type: application/json' \
--data '{
    "Latitude":32.777981,
    "Longitude":-96.796211
}'
```
#### JSON Response Body:
```
{
    "Message": "Outside it is extremely hot with light breeze and few clouds.",
    "Temp": "extremely hot",
    "Condition": "few clouds",
    "Wind": "light breeze"
}
```

## Swagger
  - TBD: please see docs

## Helpful pages
 - Need to get a latitude/Longitude for your area in decimal (DD) format? visit https://www.latlong.net/.
 - Wind Speeds are defined by [Beaufort Wind Scale to describe wind speeds](https://www.weather.gov/mfl/beaufort)
 - Temperature is based on: National Weather Service (NWS) terminology and general guidelines.

 ## Dependencies
 - Weather Service currently utilizes Open Weather Map for obtaining the weather condition at given latitude/longitude.
    - This requires our service to utilize an app_id key. Please be sure to start the docker container utilizing your own app_id key.
    - If you do not have a key visit: `https://openweathermap.org/`
- Go version 1.22
- Docker


## Improvements/Thoughts/Future Changes
- I used the decimal format for the latitude and longitude as it is easy enough to find online and it is the same format used by the Open Weather Map.
- It would be easy enough to update to utilize string format and write a package to convert from either of the following formats to decimal. This would be something I can do to improve the service in the next iteration.
    - `32°46'59.02\"N, 96°48'24.01\"W`
    - `32.7767° N, 96.7970° W`

## Contact:
- `becci.govert@gmail.com`

## Last Updated:
- 2024-07-03




© Becci Govert 2024-07-03
