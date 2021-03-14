
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

This impementation does not do any open searches. It is up to other systems to provide a list of accommodation codes (and possibly room rate codes) that match open search criteria. Most systems should be able to provide such a list. This type of filtering, e.g. by countries, regions, features, amenities is a classic use case for a relational database. If your system is able to produce a list of accommodation codes (aka hotel codes) or even accommodation codes + room rate codes but you are struggling with returning ARI (availability, rate and inventory) information, openratecache might offer you a solution.
