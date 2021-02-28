# go-producer-consumer-v2

### Lab2 - Multiple Producer - Queue -  Multiple Consumer

![Alt text](go-prod-con-2.png?raw=true "Architecture")

### Producer
  1. Expand the existing producer service to 3 replication services (3 for now) by registering each producer to a definite endpoint. Start each of the producer on a different port.
For ex:
```
      localhost:8080 -> as Producer-1
      localhost:8081 -> as Producer-2
      localhost:8082 -> as Producer-3
```
  2. Instead of directly writing to a consumer client, now the producer service should write to a Queue service (look for the description of Queue service)
  3. Since there will be multiple producers writing to a single queue, either parition the queue (sharding method) into multiple smaller queues (for better writing and consumption)  or dump the records into a buffered channel as and when a request is made to Queue service by any producer service.
  4. Once the record is successfully written to the Queue service, an acknowledgement is sent back to Producer service.

### Queue (new service)
  1. Implement a Queue service which acts as a middleman for Producer and Consumer services.
  2. Queue service registers multiple producers and multiple consumers in the system before accepting any records from the producer.
  3. Implement a queue (FIFO data structure) using array of buffered channels (as per convenience)
  3. Writing to Queue should be a simple strategy.
  4. Reading from queue should be implemented with a locking strategy (mutex/semphore)
  5. Record is popped from the queue only when it is successfully consumed by the consumer service and is stored in the database.
  6. Record remains in the queue until it is consumed. (called record retention)

### Consumer
  1. Expand the existing consumer service to 3 replication services (3 for now) by registering each consumer to a definite endpoint. Start each of the consumer on a different port.
```
      localhost:9090 -> as Consumer-1
      localhost:9091 -> as Consumer-2
      localhost:9092 -> as Consumer-3
```
   2. Instead of directly consuming the records from the producer, consumer will now be reading from the queue. Consumer constantly polls the Queue service for records. 
   3. Implement a simple record consumption strategy by either partitioning the queue and allocating the consumers such that each consumer gets to read from a certain buffered channel or from a particular index in array of buffered channels.
   4. Implement a Consumer Assignment strategy in a round-robin fashion. 
        
        If we parition Queue into 3 parts, the consumer assignment is a follows
              
            | Q-1 |  -> consumer-1
            | Q-2 |  -> consumer-2
            | Q-3 |  -> consumer-3

        If we parition Queue into 4 parts, the consumer assignment is a follows

            | Q-1 |  -> consumer-1
            | Q-2 |  -> consumer-2
            | Q-3 |  -> consumer-3
            | Q-4 |  -> consumer-1

        If we parition Queue into 6 parts, the consumer assignment is a follows

            | Q-1 |  -> consumer-1
            | Q-2 |  -> consumer-2
            | Q-3 |  -> consumer-3
            | Q-4 |  -> consumer-1
            | Q-5 |  -> consumer-2
            | Q-6 |  -> consumer-3

  5.  Implement a locking strategy during a record consumption on the queue to avoid multiple consumers reading from the same queue/channel/partition.
  6.  Send an acknowlegement back to Queue service once the record is successfully consumed and written to the database.
