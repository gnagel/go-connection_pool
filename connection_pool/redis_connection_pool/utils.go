package redis_pool

//
// Helper to iterate urls
//
func loopStrings(values []string) func() string {
	i := 0
	return func() string {
		value := values[i%len(values)]
		i++
		return value
	}
}

//
// Lazily make a Redis Connection
//
func makeLazyConnection(url string, Logger *log4go.Logger) *RedisConnection {
	// Create a new factory instance
	p := &RedisConnection{Url: nextUrl(), p.Logger}

	// Return the factory
	return p, nil
}

//
// Agressively make a Redis Connection
//
func makeAgressiveConnection(url string, Logger *log4go.Logger) (*RedisConnection, error) {
	// Create a new factory instance
	p := newLazyFactory(url, logger)

	// Ping the server
	if err := p.Ping(); nil != err {
		// Close the connection
		p.Close()

		// Return the error
		return nil, err
	}

	// Return the factory
	return p, nil
}
