# Snap Sideloader

An example web service that allows a downloaded snap to be installed.
The service listens on http://localhost:5000 and needs the restricted
`snapd-control` interface to be connected.

### Sideload a snap
Download a snap and copy the file and assertion to `/mnt`. Then run:
```
curl http://localhost:5000/snap-name/snap-revision
```

For example:
```
snap download chuck-norris-webserver
sudo cp chuck-norris-webserver_16.* /mnt/
curl -v http://localhost:5000/chuck-norris-webserver/16
```

Snapd installations are asynchronous, so installation will not occur immediately.

### List installed snaps
List the installed snaps and return the JSON response showing details of the installed
snaps.
```
curl -v http://localhost:5000/list
```

Return a JSON response similar to this:
```
{
    "type": "sync",
    "status-code": 200,
    "status": "OK",
    "result": [
        {
            "id": "J60k4JY0HppjwOjW8dZdYc8obXKxujRu",
            "title": "LXD",
            "summary": "System container manager and API",
            ...
        },
        ...
    ]
}
```
