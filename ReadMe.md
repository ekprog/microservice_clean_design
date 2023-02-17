## About

This project is example of clean design for gRPC microservices.

**We have 4 main layer:**
1) Deliver layer
2) Interactors (Use Cases)
3) Repositories
4) Domain

**1** -> **2** -> **3** -> **4**

*First layer can use any other, Second layer can use only layers 3-4 and so on.*


### Domain

Models, UseCase behaviour (interfaces). 
!! NO IMPLEMENTATION !!

### Repositories

Layer for getting data from any source. In example you can see postgres DB and test Mock implementation.

### Interactors

Here we implement main logic of our service. Whole implementation depends on domain layer behaviour.

### Delivery layer

After running of UseCases we should send result back to client. Here we can implement different clients types. \
For example, gRPC client, REST client, CLI client. \
So this layer doesn't depend on implementation of our core (interactors) and easy for testing.