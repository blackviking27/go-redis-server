import redis

# connecting to redis server
r = redis.Redis("127.0.0.1", "6379", decode_responses=True)

try:
    # get a value
    print("Sending value to server")
    print(r.ping()) # pinging tge 
    # print("response: ", r.get("foo"))

except Exception as e:
    print("Error:", e)