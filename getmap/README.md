# getmap

Little utility to cache google staticmaps api. Just pass in the API_KEY environment variable which contains the Google Maps API Key for your account.

```
docker run -d -p 8080:8080 -e API_KEY=xxxxxxxxx d.wplr.rocks/getmap
```

Then visit `http://localhost:8080/getmap?address=1600 Amphitheatre Pkwy, Mountain View, CA 94043` to get a static map.