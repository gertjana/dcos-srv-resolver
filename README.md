Small microserver that tests a docker containers ability to access SRV records on the DNS server in a marathon/mesos environment. it will do so by returning the SRV information for a given service

# build/deploy

use `make` to compile and build the docker image, `make dev` to run it locally and `make deploy` to push the image (make sure you change the organization)

Deploy the image into your mesos/marathon environment for instance by using the provided `marathon.json`

# usage

goto `http://service.domain/service/{name of service}`

you should see your service and on which host/port it's instances are located


