# cn (ceph-nano)

![Ceph, the future of storage](ceph-nano-logo-vertical.jpg)

## The project

cn is a little program written in Go that helps you interact with the S3 API by providing a REST S3 compatible gateway. The target audience is developers building their applications on Amazon S3. It is also an exciting tool to showcase Ceph Rados Gateway S3 compatibility.
This is brought to you by the power of Ceph and Containers. Under the hood, cn runs a Ceph container and exposes a [Rados Gateway](http://docs.ceph.com/docs/master/radosgw/). For convenience, cn also comes with a set of commands to work with the S3 gateway. Before you ask "why not using s3cmd instead?", then you will be happy to read that internally cn uses `s3cmd` and act as a wrapper around the most commonly used commands.
Also, keep in mind that the CLI is just for convenience, and the primary use case is you developing your application directly on the S3 API.

## Build

You can build `cn` by using `make`.
Be sure `dep` is installed:

```
$ go get github.com/golang/dep/cmd/dep
```

Then, add `~/go/bin` to your `$PATH`:

```
$ export PATH=$PATH:~/go/bin
```

Build `cn`:

```
$ make
rm -f cn cn &>/dev/null || true
dep ensure
GOOS=linux GOARCH=amd64 go build -i -ldflags="-X main.version=cea247c-dirty -X main.tag=devel -X main.branch=guits-doc_build" -o cn-devel-cea247c-dirty-linux-amd64 main.go
ln -sf "cn-devel-cea247c-dirty-linux-amd64" cn
```

Once the build is done, you should have a symlink `cn` pointing to the binary that just got built:

```
$ ls -l
total 10692
-rw-rw-r--. 1 guits guits    15292 20 nov.  22:03 ceph-nano-logo-vertical.jpg
drwxrwxr-x. 2 guits guits     4096 20 nov.  22:03 cmd
lrwxrwxrwx. 1 guits guits       34 20 nov.  22:27 cn -> cn-devel-cea247c-dirty-linux-amd64
-rwxrwxr-x. 1 guits guits 10881196 20 nov.  22:27 cn-devel-cea247c-dirty-linux-amd64
```

## Installation

cn relies on Docker so it must be installed on your machine. If you're not running a Linux workstation you can install [Docker for Mac](https://docs.docker.com/docker-for-mac/).

Once Docker is installed you're ready to start.
Open your terminal and download the cn binary.

macOS:

```
curl -L https://github.com/ceph/cn/releases/download/v2.0.4/cn-v2.0.4-darwin-amd64 -o cn && chmod +x cn
```

Linux amd64:

```
curl -L https://github.com/ceph/cn/releases/download/v2.0.4/cn-v2.0.4-linux-amd64 -o cn && chmod +x cn
```

Linux arm64:

```
curl -L https://github.com/ceph/cn/releases/download/v2.0.4/cn-v2.0.4-linux-arm64 -o cn && chmod +x cn
```

Test it out

```
$ ./cn
Ceph Nano - One step S3 in container with Ceph.

                  *(((((((((((((
                (((((((((((((((((((
              ((((((((*     ,(((((((*
             ((((((             ((((((
            *((((,               ,((((/
            ((((,     ((((((/     *((((
            ((((     (((((((((     ((((
            /(((     (((((((((     ((((
             (((.     (((((((     /(((/
              (((                *((((
              .(((              (((((
         ,(((((((*             /(((
          .(((((  ((( (/  //   (((
                 /(((.  /(((((  /(((((
                        .((((/ (/

Usage:
  cn [command]

Available Commands:
  cluster      Interact with a particular Ceph cluster
  s3           Interact with a particular S3 object server
  image        Interact with cn's container image(s)
  version      Print the version of cn
  kube         Outputs cn kubernetes template (cn kube > kube-cn.yml)
  update-check Print cn current and latest version number
  flavors      Interact with flavors

Flags:
  -h, --help   help for cn

Use "cn [command] --help" for more information about a command.
```

## Get started

Start the program with a working directory `/tmp`, the initial start might take a few minutes since we need to download the container image:

```
$ ./cn cluster start -d /tmp my-first-cluster
Running ceph-nano...
The container image is not present, pulling it.
This operation can take a few minutes......................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................

HEALTH_OK is the Ceph status
Your working directory is: /tmp
S3 access key is: 4EO5O2IP4LSGSBC8CUG3
S3 secret key is: Ap8ue6ait0MqgnCIKT7UrCRJGLUoyBhtS2McbNxE
S3 object server address is: http://10.36.116.164:8001
Ceph Nano browser address is: http://10.36.116.164:5001
```

### Selecting the cluster flavor
Is it possible to select the cluster flavor by using the ``-f `` option on the command line.

```
$ ./cn cluster start mycluster -f huge
```

The full documentation of flavors can be found [here](CONFIGURATION.md)

## Your first S3 bucket

Create a bucket with `cn`:

```
$ ./cn s3 mb my-first-cluster my-buc
Bucket 's3://my-buc/' created

$ ./cn s3 put my-first-cluster /etc/passwd my-buc
upload: '/tmp/passwd' -> 's3://my-buc/passwd'  [1 of 1]
 5925 of 5925   100% in    1s     4.57 kB/s  done
 ```

## Multi-cluster support

`cn` can manage any number of clusters on your local machine:

```
$ ./cn cluster ls
+------+---------+-------------------------------------------------------------------------------------+----------------+--------------------------------+
| NAME | STATUS  | IMAGE                                                                               | IMAGE RELEASE  | IMAGE CREATION TIME            |
+------+---------+-------------------------------------------------------------------------------------+----------------+--------------------------------+
| d    | running | ceph/daemon:latest                                                                  | master-77e3d8d | 2018-04-05T15:01:40.323603472Z |
| b    | running | ceph/daemon@sha256:369867e450ccdea9bcea7f54e97ed8b2cb1a0437fbef658d2d01fce2b8a2c648 | master-5f44af9 | 2018-03-30T21:08:31.117367166Z |
+------+---------+-------------------------------------------------------------------------------------+----------------+--------------------------------+
```

## List Ceph container images available

`cn` can list the available Ceph container images, the default output shows the 100 first images:

```
$ ./cn image ls
ceph/daemon:latest-mimic
ceph/daemon:latest-luminous
ceph/daemon:latest-master
ceph/daemon:master-0b3eb04-mimic-centos-7
ceph/daemon:master-0b3eb04-luminous-centos-7
ceph/daemon:master-0b3eb04-luminous-opensuse-42.3-x86_64
ceph/daemon:master-0b3eb04-master-centos-7-x86_64
ceph/daemon:master-0b3eb04-luminous-centos-7-x86_64
ceph/daemon:master-0b3eb04-mimic-centos-7-x86_64
[...]
```

### Using images aliases
The image option (`-i`) support aliases to simply the command line.
It is possible to list the aliases by running the `image show-aliases` command as per below :

```
$ ./cn image show-aliases
+----------+--------------------------------------------------+
| ALIAS    | IMAGE_NAME                                       |
+----------+--------------------------------------------------+
| mimic    | ceph/daemon:latest-mimic                         |
| luminous | ceph/daemon:latest-luminous                      |
| redhat   | registry.access.redhat.com/rhceph/rhceph-3-rhel7 |
+----------+--------------------------------------------------+
```
Aliases can be use in place of the traditional image name as per the following example:

```
$ ./cn cluster start mycluster -i mimic
```

It is also possible to create new aliases as detailed [here](CONFIGURATION.md)
