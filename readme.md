## Features  
1. api mock
2. reverse proxy
3. response transform
4. az cli mock

## TODO


## Install

### linux 
````
export https_proxy=http://127.0.0.1:9999
cat public.pem >> /etc/pki/tls/certs/ca-bundle.crt
update-ca-trust
````

### windows
```
set https_proxy=http://127.0.0.1:9999
Import-Certificate -FilePath public.pem -CertStoreLocation "Cert:\CurrentUser\Root"
```