go-connection_pool
==================

Abstract Connection Pool Interface written in GO.


Install
=======

	go get -u "github.com/gnagel/go-connection_pool/connection_pool"


GO Usage
========

	import connection_pool "github.com/gnagel/go-connection_pool/connection_pool"
	import redis "github.com/fzzy/radix/redis"
	
	type MyPool struct {
		createConnection func()(interface {}, err) "Lambda to create new connections"
		
		closeConnection func(interface{}) "Lambda to close connections"
		
		myPool * connection_pool.ConnectionPoolWrapper "Handle to the actual pool"
	}
	
	func NewMyPool(createConnection func()(interface {}, err), closeConnection func(interface{}), size int) (*MyPool, error) {
		var err error
		
		// Create the pool wrapper
		output = &MyPool{}
		output.createConnection = createConnection
		output.closeConnection = closeConnection
		
		// Initialize the connection pool
		output.myPool, err = connection_pool.NewPool(size, createConnection)
		
		if nil != err {
			// Abort on errors
			return nil, err
		}
		
		// Return the pool
		return output, nil
	}
	
	func (p* MyPool) Pop() (*redis.Client, error) {
		// Pop a connection from the pool
		if c := p.myPool.GetConnection(); c != nil {
			return c.(*redis.Client), nil
		}
	
		// Create a new connection
		if c, e := createConnection(); e != nil {
			return nil, e
		} elses {
			return c.(*redis.Client), nil
		}
	}
	
	func (p* MyPool) Push(c *redis.Client, e error) {
		// Close the connection on errors
		if nil != e {
			closeConnection(c)
			c = nil
		}
		
		// Nil is a placeholder to tell "Pop" to re-create the connection
		c.myPool.ReleaseConnection(c)
	}


Authors:
========

Glenn Nagel <glenn@mercury-wireless.com>, <gnagel@rundsp.com>


Credits:
========

Ryan Day's original implementation that inspired this is here: [rday's gist](https://gist.github.com/rday/3504674)
