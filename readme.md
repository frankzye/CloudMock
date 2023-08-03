tf mock is a great tool help run terraform codes without cloud support.

## Features  
1. terraform provider requests mock
2. az cli mock
3. python modules integrate with tftest

## Future 


## Install
```pip
pip install tfmock
```

don't forget trust the certificate if you are in windows or mac, find the certificate inside the tfmock-xx.egg-info bin folder

### windows
```
Import-Certificate -FilePath public.pem -CertStoreLocation "Cert:\CurrentUser\Root"
```

### macos
```mac
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain public.pem
```