//
// Redis Connection Pool written in GO
//

package redis_pool

import "time"
import "github.com/fzzy/radix/redis"
import "github.com/alecthomas/log4go"
import cp "../connection_pool"

//
// Constants for connecting to Redis & Logging
//
const timeout = time.Duration(10) * time.Second
const log_client_premade_success = "[RedisConnection][Client] - Pre-made availble for %s"
const log_client_premade_failure = "[RedisConnection][Client] - No saved connection available for %s"
const log_open_failed = "[RedisConnection][Open] - Failed to connect to %s, error = %#v"
const log_open_success = "[RedisConnection][Open] - Opened new connection to %s"
const log_closed = "[RedisConnection][Close] - Closed connection to %s, error = %#v"

//
// Connection Wrapper for Redis
//
type RedisConnection struct {
	Url string "Redis URL this factory will connect to"

	logger log4go.Logger "Handle to the logger we are using"

	client *redis.Client "Connection to a Redis, may be nil"
}

//
// Get a connection to Redis
//
func (p *RedisConnection) Client() (*redis.Client, error) {
	// If the connection is valid, return it
	if nil != p.client {
		// Log the event
		p.logger.Trace(log_client_premade_success, p.Url)

		// Return the connection
		return p.client, nil
	}

	// Log the event
	p.logger.Warn(log_client_premade_failure, p.Url)

	// Open a new connection to redis
	if err := Open(); nil != err {
		// Errors are already logged in Open()
		return nil, err
	}

	// Return the new redis connection
	return p.client, nil
}

//
// Ping the server, opening the client connection if necessary
// Returns:
//   nil   --> Ping was successful!
//   error --> Ping was failure
//
func (p *RedisConnection) Ping() error {
	// Open the connection to Redis
	client, err := p.Client()
	if nil != err {
		return err
	}

	// Ping the server
	client.Append("ping")

	// Get the response
	reply := client.GetReply()

	// Connection error? Then tell the factory to invalidate the Redis connection
	if nil != reply.Err {
		// Close the connection
		p.Close(reply.Err)

		return reply.Err
	}

	// Return nil on Success!
	return nil
}

//
// Return true if the client connection exists
//
func (p *RedisConnection) IsOpen() bool {
	return p.client != nil
}

//
// Open a new connection to redis
//
func (p *RedisConnection) Open() (err error) {
	// Connect to Redis
	p.client, err = redis.DialTimeout("tcp", p.Url, timeout)

	// Error connecting?
	if nil != err {
		// Clear the connection pointer
		p.client = nil

		// Log the event
		log.Critical(log_open_failed, p.Url, err)
	} else {
		// Log the event
		p.logger.Info(log_open_success, p.Url)
	}

	return
}

//
// Close the connection to redis
//
func (p *RedisConnection) Close(err error) {
	// Log the event
	log.Warn(log_closed, p.Url, err)

	// Close the connection
	if nil != p.client {
		p.client.Close()
	}

	// Set the pointer to nil
	p.client = nil
}
