Central API
===========

This is the application run on the central server and is incharge of handling scraping the model sites and handling
online performers, as well as dispatching downloading requests.

# Expected consul kv keys

The following keys must exist in consul 

- central/rabbitmq/user
- central/rabbitmq/pass
- central/rabbitmq/vhost
- central/aphrodite/api
- central/aphrodite/token
- central/hermes/api
- central/hermes/token
- central/camgirlgallery/api
- central/camgirlgallery/token
- central/ultron/api
- central/ultron/token
