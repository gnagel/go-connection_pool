package redis_connection_pool

import "os"
import "testing"
import "github.com/alecthomas/log4go"

var redis_master = os.Getenv("REDIS_MASTER")

const method_expected_a_equal_b_from_url = "[%s][%s] Expected '%v', Actual '%v'"

//
// New connections are not open
//
func Test_RedisConnection_1(t *testing.T) {
	tag := "[RedisConnection][NewConnection IsNotOpen] Actual = %#v, Expected %#v"

	// Create a connection, but don't actually connect to redis yet ...
	connection := RedisConnection{Url: "127.0.0.1:6999", logger: log4go.NewDefaultLogger(log4go.DEBUG)}
	defer connection.Close()

	// We shouldn't be connected yet
	if ok := connection.IsOpen(); !ok {
		t.Errorf(tag, ok, false)
		return
	}
}

//
// Connections to Non-existent Redis Servers should fail
//
func Test_RedisConnection_2(t *testing.T) {
	tag := "[RedisConnection][Non-existent Redis Fails] Actual = %#v, Expected %#v"

	// Create a connection, but don't actually connect to redis yet ...
	connection := RedisConnection{Url: "127.0.0.1:6999", logger: log4go.NewDefaultLogger(log4go.DEBUG)}
	defer connection.Close()

	// Expected to fail
	if err := connection.Open(); nil == err {
		t.Errorf(tag, err, "!nil")
		return
	}

	// We shouldn't be connected yet
	if ok := connection.IsOpen(); !ok {
		t.Errorf(tag, ok, false)
		return
	}

	// Start the redis server
	cmd := exec.Command("redis-server", "--port", "6999")
	if err := cmd.Start(); err != nil {
		t.Logf("skipping test; couldn't start redis-server")
		return
	}
	defer cmd.Wait()
	defer cmd.Process.Kill()

	// Expected to succeed
	if err := connection.Open(); nil != err {
		t.Errorf(tag, err, nil)
		return
	}

	// We should be connected now
	if ok := connection.IsOpen(); ok {
		t.Errorf(tag, ok, true)
		return
	}
}

//
// Ping-ing the Redis connection should automatically-connect
//
func Test_RedisConnection_2(t *testing.T) {
	tag := "[RedisConnection][Ping Reconnects to Redis] Actual = %#v, Expected %#v"

	// Start the redis server
	cmd := exec.Command("redis-server", "--port", "6999")
	if err := cmd.Start(); err != nil {
		t.Logf("skipping test; couldn't start redis-server")
		return
	}
	defer cmd.Wait()
	defer cmd.Process.Kill()

	// Create a connection, but don't actually connect to redis yet ...
	connection := RedisConnection{Url: "127.0.0.1:6999", logger: log4go.NewDefaultLogger(log4go.DEBUG)}
	defer connection.Close()

	// Expected to succeed
	if err := connection.Ping(); nil != err {
		t.Errorf(tag, err, nil)
		return
	}

	// We should be connected now
	if ok := connection.IsOpen(); ok {
		t.Errorf(tag, ok, true)
		return
	}
}
