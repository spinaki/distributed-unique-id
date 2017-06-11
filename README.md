## Distributed 64 bit Long Globally Unique ID Generator
Generating globally unique integer ids is a non-trivial task. Integer/ Long ids are
 almost always better than string ids since they help create faster indexes.
Further, if we want the ids to be ordered -- it adds another level of complexity. 
Time is an important dimesion for sorting. If the ids geenrated by the system are (pseudo) sorted by time -- we get temporal sorting for free, 
since in most databases records are already sorted by the keys (in this case they will be by default sorted by time).

The goal of this project was to create a highly available, low latency system to generate 64 bit long globally unique ids.
I adapted Twitter's Snowflake to build a fast id generator -- where ids are sorted by time.
I used Kubernetes and Docker to build a highly available system -- so that it can be deployed in the cloud very quickly.
## Take it for a quick spin
Find the image at [DockerHub](https://hub.docker.com/r/exifguy/uniqueid/). Pull the image and launch the the container.
Launch the ID service without writing a single line of code :).
```
docker pull exifguy/uniqueid:v1
docker run --name puid -p 8080:8080 -d exifguy/uniqueid:v1
curl localhost:8080/status
curl localhost:8080/longids
```

#### What can you do with this library
* You can generate a sorted list of  64 bit unique ids.
* As long as they are generated from same machine / cluster -- they are guaranteed to be unique.
* You can also generate (pseudo)random string ids, which are likely to be unique.

#### HowTo Configure
* Currently you can set the start time of the long id generator via Settings. See example in `IdService.go`
#### REST API Endpoints
* `/longids`: return a sorted list of 64 bit long ids (length: 256)
* `/longidrange`: returns two 64 bit long ids, the first and the last in a sorted set of 256 ids.
* `/stringids`: returns a set of n random string ids. Input params:
  * `num`: num of ids (default 10).
  * `len`: length in bytes of the ids. The greater this value is -- higher is the randomization and lower chance of collision.
  The default value is 32 bytes or 256 bits.
#### HowTo Run Locally via Go Binary
```
go build -v
./uniqueidgenerator // starts the server listening on port 8080
curl localhost:8080/longids // invoke the api endpoint from another cli
```
#### Latency
* Its really, really fast! The latency while running locally was microseconds to 10 milliseconds.
<p align="center">
<img src="unique-id-time.png?raw=true" width="450"/>
</p>

#### Howto Create Your Own  Image
* Create the binary which whill be used by the docker image. Run the following from the main directory.
Make sure the gopath is set correctly.
```
env GOOS=linux GOARCH=amd64 go build -v
```
* Create image
```
docker build -f Dockerfile -t exifguy/uniqueid:v1 .
```
* Run The image in a container
```
docker run --name puid -p 8080:8080 -d exifguy/uniqueid:v1
```
* Test
```
curl localhost:8080/longids
```
* Push
```
docker commit -m "initial commit" puid exifguy/uniqueid:v1
docker push exifguy/uniqueid:v1
```
#### Deploy on Kubernetes (either AWS or GCP)
* By deploying on  Kubernetes you can achieve high availability -- since if a node
goes down the system will restart another one.
* Use the  [Deployment YML File](./unique-id-deployment.yaml) to create a new deployment on kubernetes.
* Make sure about the pod-ip environment variable is added in the [yml file](./unique-id-deployment.yaml). 
The library running in the pod will query this variable to find the ip address required for unique id.
#### Details on the long Unique Id Format
* From MSB to LSB: 39 bit Time, 16 bit machine ID, 8 bit sequence.
* Most Significant Bits are time so that IDs can be sorted based on time.
* At every query, 256 ids (in sequence) are generated.
##### References
* [Twitter Dep Code](https://github.com/twitter/snowflake)
* [Slides](https://www.slideshare.net/davegardnerisme/unique-id-generation-in-distributed-systems)
* [Sonyflake](https://github.com/sony/sonyflake)