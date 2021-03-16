
## openratecache is under development and not even ready for evaluation, let alone production! ##

# openratecache
Go implementation of the RateCache core.

This is an implementation based on data formats and concepts as of the commercially licensed RateCache. Openratecache efficiently stores rate and availability information which can be queried through a web service with extremely low response times.

OpenRateCache consists of two services, one for writing data to the cache and another one for querying this data. There are various reasons for splitting up services into Read and Write operations:

 - More options for fine-grained optimization of set-ups
 - While write operations use normal file descriptors, the file is memory-mapped for read. 

Protection of concurrent access is left to the underlying os, there is no mutex on ws level for performance reasons.

The cached ARI information is kept in one file on disk plus an index which is kept in memory (the index is protected by a mutex for concurrent access). Testing in other implementations have shown that reads of ari data directly from disk outperform database operations for previous filtering. However, if you really want to improve access speed you can keep the cache file on a ram disk.

## Integration of openratecache with existing systems ##

This implementation does not do any open searches. However, you may submit an almost unlimited amount of codes in a single request.

It is up to other systems to provide a list of accommodation codes (and possibly room rate codes) that match open search criteria. 

Most systems should be able to provide such a list. This type of filtering, e.g. by countries, regions, features, amenities, etc. is a classic use case for a relational database. If your system is able to produce a list of accommodation codes (aka hotel codes) or even accommodation codes + room rate codes but you are struggling with returning ARI (availability, rate and inventory) information, openratecache might offer you a solution.

## Setup ##

## Operation ##

Except for GET which uses a normal query string, all other operations use json for both request and response. Make sure to use utf-8, especially when you are on Windows.

### Rate import format ###

Data is posted to

http://your.url/import

The format looks as follows:
´´´
{
    "accommodationCode":"AAL00324",
    "roomRateCode":"DBLFRHB396",
    "Occupancy":[
        {
            "minAge":3,
            "maxAge":16,
            "count":2
        },
        {
            "minAge":17,
            "maxAge":100,
            "count":2
        }
    ],
    "rates":[
        {
            "firstCheckIn":"2021-03-11",
            "lastCheckIn":"2021-03-22",
            "lengthOfStay":1,
            "rate":31.02
        },
        {
            "firstCheckIn":"2021-03-23",
            "lastCheckIn":"2021-03-24",
            "lengthOfStay":1,
            "rate":26.27
        }
        ...
    ],
    "availabilities":[
        {
            "firstCheckIn":"2021-03-11",
            "lastCheckIn":"2021-03-18",
            "lengthOfStay":1,"available":7
        },
        {
            "firstCheckIn":"2021-03-19",
            "lastCheckIn":"2021-03-22",
            "lengthOfStay":1,
            "available":6
        },
        ...
    ]
}
´´´
#### Codes ###
The maximum length of the codes is configurable and can be up to 255 chars long. 
Important char == byte! Keep this in mind if you use utf-8 and characters not in 8bit ascii, as
these characters may occupy more than one byte! Apart from that it is entirely up to you how
you compose them. It is up to you how you to parse and make sense of these codes.
#### accommodationCode ####

Accommodation codes (aka hotel codes or contract codes) are strings that should uniquely identify the hotel contract
and may be composed of various codes from your system. You should include any codes that you need for booking.

#### roomRateCode ####

As the name suggests these codes should uniquely identify a room and the associated rate plan and meal plan. As with
accommodation, the roomRateCode should contain all codes that are required for booking.

#### occupancy ####

While the closed source version handles adults and children differently, openRateCache just uses age ranges
to define the occupancy. Make sure that age ranges do not overlap and you do not want to leave any gaps
between ages either.

Example:
```
{
            "minAge":3,
            "maxAge":16,
            "count":2
        },
        {
            "minAge":17,
            "maxAge":100,
            "count":2
        }
```
The above example specifies that the rate applies for an occupancy with two guests between 3 and 16 years
and two guests older than 17 (the maxAge is set to 100 because it needs to be set to something.)

#### rates ####

Rates (prices) apply to a check-in date and a length of stay. If the rate for a specific los does not
change for various consecutive check-in dates, conform one group that looks as follows:

```
{
    "firstCheckIn":"2021-03-23",
    "lastCheckIn":"2021-03-24",
    "lengthOfStay":5,
    "rate":250.00
}
```
While in the closed source version every room rate may have a different currency and digits of currencies
are taken into account, this implementation only accepts one currency for the whole cache. It also assumes
that the currency has two digits.

This open source version is just a dumb cache, there is currently no possibility to plug in a mark-up engine
as in the closed source version.

The maximum price that can be stored for a check-in/los combination is 2684354.56. Please take this into account
if you are planning to use the cache with currencies that require more digits.

#### available ####

Availability is limited to 15. This may not be not enough to represent the whole allotment in your inventory
system but it is more than ok for searches. Nobody will ever book more than three rooms over the internet
for a sinlge party! Interpretation is up to you but recommended is:

- If rate is 0 the room is closed for sale
- If available is 0 the room is open for sale but the number of availability is unknown.
- If available is 1 - 14, these numbers represent the number of available rooms.
- If available is 15, there are at least 15 rooms available.

### Rate and availability updates ###




### Request format ###

## FAQ ##

### Why not use a conventional RDBMS?
    
Relational databases are good at what they are designed for and some of them such as MariaDB or PostGreSQL have amazing support for geospatial data which opens a whole lot of possibilities for the implementation of open searches. However, all RDBMS struggle with the sheer amount of data when it comes to rates.

Depending on the number of rate and meal plans, possible occupancies, check-in dates and maximum length of stay you end up with 30.000 to 80.000 rates, **just for one room**! This can be greatly compressed by storing the rates with date ranges, i.e. you store start and end dates for the ranges within which the rate doesn't change. This, on the other hand creates a huge overhead in times with frequent rate changes and increases the complexity of search for the DBMS. Anybody who has ever tried to do this (or has actually done this) knows, that this does not work very well. The scenario is even worse for NoSQL dbs.

### Why not use key-value stores? 
    
Key-value stores such as redis, DynamoDB BerkleyDB, etc. will usually not perform very well if you use dates as keys and will require a huge amount of resources, i.e. Memory. Huge means really huge.



