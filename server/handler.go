package main

import "sync"

// To store the key-value
var SETs = map[string]string{}

// Creating a mutex to allow different go routines to access the shared map
// Can Read the following article to understand better: https://medium.com/bootdotdev/golang-mutexes-what-is-rwmutex-for-5360ab082626
var SETsMu = sync.RWMutex{}

// Hash sets
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

// PING Command
func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

// Function to return the value from the map
func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}

}

// Function to set value in the Map
func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// Function to set the value in hash set\
func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	// checking if the hash is already present or not
	// if not present, then create a map for the given hash
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}

	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// Function to get the value from the hash set
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' function"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}

}
