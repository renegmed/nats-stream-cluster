## Learning NATS Streaming ##


This application showcases the use of NATS cluster streaming.

1. Create and run the containers - cluster nats streams, publisher api and subscribers

    $ make up

2. Monitor publisher and subscribers. Use separate terminal on each item.
    

    $ make tail-api

    $ make tail-curious

    $ make tail-patient

    $ make tail-binge

3. Publish a json data

    $ make publish



