
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

## FAQ ##

Why not use a conventional RDBMS?
    Relational databases are good at what they are designed for and some of them such as MariaDB or PostGreSQL have amazing support for geospatial data which opens a whole lot of possibilities for the implementation of open searches. However, all RDBMS struggle with the sheer amount of data when it comes to rates.

    Depending on the number of rate and meal plans, possible occupancies check-in dates and maximum length of stay you end up with 30.000 to 80.000 rates, **just for a single room**! This can be greatly compressed by storing the rates with date ranges, i.e. you store start and end dates for the ranges within which the rate doesn't change. This, on the other hand creates a huge overhead in times with frequent rate changes and increases the complexity of search for the DBMS. Anybody who has ever tried to do this (or has actually done this) knows, that this does not work very well. The scenario is even worse for NoSQL dbs.

Why not use key-value stores? 
    Key-value stores such as redis, DynamoDB BerkleyDB, etc. will usually not perform very well if you use dates as keys and will require a huge amount of resources, i.e. Memory. Huge means really huge.

