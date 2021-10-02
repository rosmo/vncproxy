# VNC proxy with reMarkable 2.10+ authentication support

Forked from: 

Usage:
  - Compile the proxy
  - Run the proxy with the `-reMarkable DEVICE_ID` flag

To extract the device ID:
  - Log into reMarkable via SSH
  - Extract the `devicetoken` string (exclude the `@ByteArray` wrapper) the string from `/etc/remarkable.conf` 
  - Decrypt the device ID by running:
```bash
pip3 install --user PyJWT
python3 -c 'import sys,jwt;t=jwt.decode(sys.argv[1],options={"verify_signature":False});print(t)' '(DEVICE TOKEN HERE)
````
  - In output, you should get a string starting with `auth0|`. The whole string is your device ID which should be passed
    to be `-reMarkable` flag.

Works with:
  - TightVNC

Doesn't work with:
  - RealVNC
  - NoVNC