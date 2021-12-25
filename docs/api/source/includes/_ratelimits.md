# Rate Limits

There are 3 tiers of rate-limiting:    
Light rate-limiting: 1 request per 4 seconds with a max burst of 2.    
Standard rate-limiting: 5 requests per second with a burst of 5.    
Heavy rate-limiting: 10 requests per second with a burst of 10.    

The rate limits work as a 'bucket' system where there is a max (burst) and a regenerative rate. More info can be found in the source file server_middleware.go.

Each endpoint will depict which rate limit applies.

For example:
On the medium rate-limiter, you start with 5 requests available to you instantly, and this replenishes 5 requests per second.
