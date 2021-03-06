## Capesize (ALPHA)

An internal tool for spawning Docker hosts on the cloud. Supports Amazon's EC2, with GCE support coming soon. Is currently ALPHA / used in testing / and will likely be replaced by [Docker Machine](https://github.com/docker/machine)

> *Capesize* ships are the largest cargo ships; ships which are too large to transit the Suez Canal (Suezmax limits) or Panama Canal (Panamax limits), and so have to pass either the Cape of Good Hope or Cape Horn to transverse between oceans. [Wikipedia](http://en.wikipedia.org/wiki/Capesize)

#### Commands
* capesize *create* <provider> <num hosts>
* capesize *list* <provider>

#### Getting started

```
go build
./capesize create amazon 5
./capesize list amazon
```

###### Example to create 5 docker hosts on Amazon:
```
AWS_ACCESS_KEY_ID=abc123 \
AWS_SECRET_ACCESS_KEY=def456 \
SECURITY_GROUP="Backend Servers" \
EC2_KEY_PAIR_NAME=jenkinskey \
IDENTITY_FILE=~/.ssh/jenkins_id_rsa" \
./capesize create amazon 5
```

###### Required ENV vars:

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY
* SECURITY_GROUP
* EC2_KEY_PAIR_NAME
* IDENTITY_FILE
* DEVELOPER_KEYS

###### Optional ENV vars (default):

* BUILD_IDENTIFIER (capesize)
* EC2_IMAGE_ID (ami-56b7eb04) - Amazon docker ready AMI
* EC2_HOST_USER (ec2-user)
* EC2_INSTANCE_TYPE (m3.medium)
* EC2_AVAILABILITY_ZONE (ap-southeast-1b)
* REMOTE_DIR_PATH (opt)
* DOCKER_OPTS (e.g. --insecure-registry xx.xx.xx.xx:5000)
* DOCKER_RUN_POST_INSTALL (e.g. docker run ...; docker run ...)
* STATIC_MACHINE_NAME (creates machines with a static name, instead of random generation per instance)

#### Todo
* Better error handling. i.e. refactor excessive & abusive use of `panic`
* Better logging
* Add GCE support & ensure design is flexible across providers

#### License
MIT

Copyright (c) 2014 Zumata Technologies Pte Ltd.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
