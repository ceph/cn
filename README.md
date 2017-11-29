# cn (ceph-nano)

## The project

cn is a little program written in Go that helps you interact with the S3 API by providing a REST S3 compatible gateway. The target audience is developers building their applications on Amazon S3. It is also an exciting tool to showcase Ceph Rados Gateway S3 compatibility.
This is brought to you by the power of Ceph and Containers. Under the hood, cn runs a Ceph container and exposes a [Rados Gateway](http://docs.ceph.com/docs/master/radosgw/). For convenience, cn also comes with a set of commands to work with the S3 gateway. Before you ask "why not using s3cmd instead?", then you will be happy to read that internally cn uses `s3cmd` and act as a wrapper around the most commonly used commands.
Also, keep in mind that the CLI is just for convenience, and the primary use case is you developing your application directly on the S3 API.

Available calls are:

```
$  ./cn s3 -h
Interact with S3 object server

Usage:
  cn s3 [command]

Available Commands:
  mb          Make bucket
  rb          Remove bucket
  ls          List objects or buckets
  la          List all object in all buckets
  put         Put file into bucket
  get         Get file into bucket
  del         Delete bucket
  du          Disk usage by buckets
  info        Get various information about Buckets or Files
  cp          Copy object
  mv          Move object
  sync        Synchronize a directory tree to S3

Flags:
  -h, --help   help for s3

Use "cn s3 [command] --help" for more information about a command.
```

## Installation

cn relies on Docker so it must be installed on your machine. If you're not running a Linux workstation you can install [Docker for Mac](https://docs.docker.com/docker-for-mac/) or [Windows](https://docs.docker.com/docker-for-windows/).

Once Docker is installed you're ready to start.
Open your terminal and download the cn binary.

macOS
``` 
$ curl -L https://github.com/ceph/cn/releases/download/v1.0.0/cn-v1.0.0-darwin-amd64 -o cn && chmod +x cn
```

linux
```
$ curl -L https://github.com/ceph/cn/releases/download/v1.0.0/cn-v1.0.0-linux-x86 -o cn && chmod +x cn
```

Test it out
```
$ ./cn
Ceph Nano - One step S3 in container with Ceph.

Usage:
  cn [command]

Available Commands:
  start       Start object storage server
  stop        Stop object storage server
  status      Stat object storage server
  purge       Purge object storage server. DANGEROUS!
  logs        Print object storage server logs
  restart     Restart object storage server
  s3          Interact with S3 object server
  update      Update the container image
  version     Print the version number of Ceph Nano
  help        Help about any command

Flags:
  -h, --help   help for cn

Use "cn [command] --help" for more information about a command.
```

## Get started

Start the program with a working directory `/tmp`, the initial start might take a few minutes since we need to download the container image:

```
$ ./cn start -d /tmp
Running ceph-nano...
The container image is not present, pulling it.
This operation can take a few minutes......................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................

HEALTH_OK is the Ceph status
S3 object server address is: http://192.168.0.10:8000
S3 user is: nano
S3 access key is: E7MCNPWED4BIV21J1BZG
S3 secret key is: JG5wmxw2bGOxXLc7v6NQ2yg50atPCu3Nxe4XXvEf
Your working directory is: /tmp
```

##  Work with this CLI

Your first S3 bucket!

```
$ ./cn s3 mb my-buc
Bucket 's3://my-buc/' created

$ ./cn s3 put /etc/passwd my-buc
upload: '/tmp/passwd' -> 's3://my-buc/passwd'  [1 of 1]
 5925 of 5925   100% in    1s     4.57 kB/s  done
 ```
