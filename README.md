# calendar-api

a very rudimentary calendar image generator API

## !! DISCLAIMER !! this service does not use any caching, authentication, sophisticated rate-limiting, etc. use at your own risk

# usage

- send a POST request to the /calendar endpoint with a list of event objects as the body.

```
[
  {
    "title": "test",
    "day": 0,
    "start_time": 12:00,
    "end_time": 15:00,
    "color": "green"
  }
]
```
