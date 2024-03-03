Develop a rate limiting middleware for a high-traffic web service written in Go Lang. The middleware should restrict the number of requests processed per client IP address within a specified time window to prevent abuse and ensure fair resource allocation. The key requirements for the middleware are as follows:
1. Dynamic Configuration: Allow dynamic configuration of rate limits for different endpoints and client IP addresses.
2. Efficiency: Implement efficient data structures and algorithms to track request rates and enforce rate limits with minimal overhead.
3. Concurrency Safety: Ensure concurrency safety to handle concurrent requests from multiple clients simultaneously.
4. Expiry Mechanism: Implement a mechanism to expire outdated rate limit records and reclaim resources efficiently.

Your solution should include the following components:
- A Go Lang middleware function to intercept incoming requests and enforce rate limits.
- Data structures and algorithms to track request rates and enforce rate limits.
- Configuration options for specifying rate limits for different endpoints and client IP addresses.
- Implementation of concurrency-safe mechanisms for tracking request rates and enforcing rate limits.
- Thorough testing to validate the correctness and performance of the rate limiting middleware.
