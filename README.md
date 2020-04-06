# VolantMq Plugin Auth Mongo
Volantmq Plugin for supporting Authentication using MongoDb. 


If you are using Mongodb in your application. You can easily connect your database and use it as backend for authentication in VolantMQ. 

## Expected config in config.yaml 
```yaml
- name: authmongo
        backend: mongo
        config:
          mongoURI: mongodb://localhost:27017/mqtt
          database: mqtt
          collection: mqtt_user
```

