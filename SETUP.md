# Set-up guide for OpenRateCache #

## Requirements ##
OpenRateCache does not require a lot of resources if all you want to do 
is play around for testing and integration. Any common Linux distro
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

Open /opt/conf/wswrite.conf. After editing make sure the file
is still valid json. You might use an online json validator
for this purpose.

- cacheDir: points to the directory where the cache file lives. 
  Set this to `/var/local/ratecache`
- indexDir: points to the directory for the index file. 
  Set this to `/var/local/ratecache`
- cacheFilename: You can use any valid filename. You migth set this
  to the purpose of the cache, e.g. "demo", "test", "int", "prod", etc.
  or if you cache data from different suppliers in different cache 
  instances you might use the supplier code such as "HOB", "WBEDS", etc.
- supplier: Set this to a suitable supplier code, if you only cache
  your own data or a mixed bag you can set this to your own code. Currently
  the code is only used when initializing a new cache file and is set in 
  the header. This might change in future versions.
- currency: Set this to the 3-digit ISO currency code. This information will
  only be used when initializing a new cache file. Once a cache file exists
  the currency is taken from the cache file header.
- "decimalPlaces: set this to the number of decimal places you want to use.
  If you use a currency that has decimal places such as GBP, EUR, USD, etc.
  you can still decide to set this to 0 which means that no decimal places
  are stored or returned upon request.
- maxLos: Set this to the maximum length of stay. If you try to import rates 
  for lengths of stay > then maxLos these are dicarded. Be aware that setting
  this to a value that is much bigger than what you are actually going to 
  import into the cache will unnecessarily blow up the cache size which in 
  turn will have an impact on performance.
- "days": the number of check-in dates for which you are going to import data.
  If you want to import rates for the next two seasons you might set this to
  365. However it is a good idea to add an extre couple of weeks or so to this
  value, i.e. in the abobe example you might want to set this to 380.
- accoCodeLength: The accommodation code (aka hotel code, contract code...)
- can be up to 255 bytes long. IMPORTANT: bytes != characters, utf-8 code
  points are of variable width and can take up more than one byte if you 
  decide to use non-ascii chars. You can compose this code from various of 
  your codes. In any case this code should contain everything required on hotel
  level for later booking. Adjust this length to the maximum length of the 
  code on hotel level you are going to use.
- roomRateCodeLength: Everything mentioned for accoCodeLength goes for 
  the roomRateCode, but on room level. This code commonly contains information
  on toom type, room category, rate plan and meal plan.
- initialRateBlockCapacity: when the cache file is created for the first time
  a number of empty rate blocks are already created. You can set this to 
  an approximate number, i.e. if you plan to import 1.000 hotels and you 
  count with an average of 20 room rates per hotel you might set this
  to 20.000. If the number is smaller thant the real number of room rates
  the cache file will simply start to grow.
  addIndexUrsl: Whenever a new block of rates for a room rate / occupance
  combination is added, this data will be written not only to the cache
  but also added to the index and appended to the index on disk. In order
  to use the imported data the web services for search need to add this
  information to their own index (or re-load the idxe from disk) This parameter
  expects a list of urls to which the new index information is sent.
  notify: you can switch off the notification of new index entries by setting
  this to false. 
  