# Set-up guide for OpenRateCache #

## Requirements ##
OpenRateCache does not require a lot of resources if all you want to do 
is play around for testing and integration. Any Linux common Linux distro
for x86-64 will do.

OpenRateCache contains two web services and does not come with rc.d 
or sysctl scripts to start them up. It is up to you to decide 
how to run these services. This guide assumes you are going
to use supervisord to run the OpenRateCache services.

You will need sudo rights or login as root to follow these instructions
although for ad-hoc testing purposes you can even run these services as
a normal user.

## Step-by-step proceedings ##

The instructions assume you are logged in as root or you have opened
a bash shell with root privileges:
```
sudo bash
```

### Create user and directories ###

create user "openratecache"
```
useradd -r openratecache
```

create install directory
```
mkdir /openratecache/
mkdir /openratecache/bin
mkdir /openratecache/conf
chown -R openratecache:openratecache openratecache
```

Untar the tarball in a temporary directory and copy 
the files into the install directory

```
tar xvfz openratecache.tgz
mv wswrite wssearch demodatagen /openratecache/bin
mv wsrite.conf wssearch.conf /openratecache/conf
chown -R openratecache:openratecache openratecache/*
chmod 744 /openratecache/bin/*
chmod 644 /openratecache/conf/*.conf
```

Create a directory for the cache file
```
mkdir /var/local/openratecache
chown openratecache:openratecache /var/local/openratecache
```
### edit configuration files ###
