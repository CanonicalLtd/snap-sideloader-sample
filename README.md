# Snap Sideloader

An example web service that allows a downloaded snap to be installed.
The service listens on http://localhost:5000 and needs the restricted
`snapd-control` interface to be connected.

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
