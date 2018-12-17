# Using cluster flavors
Ceph nano offers a set of pre-defined flavors to tune your cluster:

|Flavor |RAM (MB)  |CPU    |
|-------|------|------|
|default| 512  | 1   |
|medium |768 | 1 |
|large   |1024   |1   |
|huge   |4096   |2   |

## Listing available flavors
The flavor list can be obtained by running the `flavors ls` command:

```
$ cn flavors ls
+---------+-------------+-----------+
| NAME    | MEMORY_SIZE | CPU_COUNT |
+---------+-------------+-----------+
| default | 512MB       | 1         |
| medium  | 768MB       | 1         |
| large   | 1GB         | 1         |
| huge    | 4GB         | 2         |
+---------+-------------+-----------+
```

## Getting details of a flavor
A flavor can also define some other tuning parameters. They can be listed by using the `flavors show` command:

```
$ cn flavors show default

Details of the flavor default:
{
  "cpu_count": 1,
  "data": "",
  "memory_size": "512MB",
  "privileged": false,
  "size": "",
  "use_default": true,
  "work_directory": "/usr/share/ceph-nano"
}
```

## Using a flavors
Flavors can be selected with the `-f` flags when starting a new cluster.
```
$ cn cluster start mycluster -f huge
```

## Overriding existing flavors
By using the configuration file, it is possible to override existing flavors by overriding the existing values.

The following example override the amount of memory for the large flavor by editing `/etc/cn/cn.toml`

```
[flavors.large]
   memory_size="2GB"
```

## Creating new flavors
It is possible to create new flavors in the configuration file by creating a `flavors.<something>` entry into it.

The following example creates a new flavors called `superfat` by editing `/etc/cn/cn.toml`

```
[flavors.superfat]
    memory_size="8GB"
    cpu_count=4
```


## Configuration items of a flavor
A flavor can be tuned with various built-in items as defined below:


|Item | Role | Default value| Cli flag |
|-----|------|--------------|----------|
| cpu_count  |   Set the amount of processors | 1   | none|
|  memory_size | Set the amount of memory   | 512MB  | none |
|  work_directory |  Set the working directory   | /usr/share/ceph-nano  | -d  or --work-dir |
|  data | Set the underlying storage with a specific directory or physical block device  |  none | -b or --data  |
| size  |  Set the underlying storage size when using a specific directory | none   |  -s or --size |
|privileged   | Defines if the container runs in privileged mode  |   false | none  |
| use_default   | Defines if this flavor inherit from the `default` flavor  | true  | none  |

If a flavor defines a `ceph.conf` sub entry, this one will be used as items for the ceph.conf configuration as per bellow:

```
[flavors.huge]
   memory_size=2G
   [flavors.huge.ceph.conf]
      osd_memory_target = 536870912
      osd_memory_base = 268435456
```

# Images aliases
To ease the usage of ceph nano, it is possible to use aliases instead of regular image names.

## Default Aliases
The following table list the default builtin aliases

| Alias    | Image name |
|----------|--------------------------------------------------|
| mimic    | ceph/daemon:latest-mimic                         |
| luminous | ceph/daemon:latest-luminous                      |
| redhat   | registry.access.redhat.com/rhceph/rhceph-3-rhel7 |

## Managing image aliases
It is possible to add new image aliases by adding a custom entry in the `[images]` section of the configuration file.

| Item | Role |Default value|
|------|------|-------------|
|image_name | Defines the associated image name for an alias| ceph/daemon|


The following example creates an alias name `sharktopus` in the `/etc/cn/cn.toml` configuration file.
```
[images.sharktopus]
  image_name="ceph/daemon:latest-sharktopus"
```

# Configuration file
Ceph nano can read its configuration from 3 different locations, they are search in the following order:
- /etc/cn/cn.toml
- ~/.cn/cn.toml
- ./cn.toml

If no configuration file is found, the builtin flavors and image aliases are used.

If a configuration file is found, additional flavors and images aliases are loaded from it.

By default, `use_default` feature is enabled meaning that all user-defined but also built-in flavors inherit for every feature of the `default` flavor unless it's defined in the new flavor.

If `use_default` is set to false, the user-defined flavor will have not inheritance from `default` meaning all the definition of this flavor is done in the flavor section.


In the following example, the `test1` flavor will only have `cpu_count` set to 2 while the `test2` flavor will have both `cpu_count` set to 2 and `memory_size` set to 640MB.

```
[flavors.default]
  memory_size=640MB

[flavors.test1]
  use_default=false
  cpu_count=2

|flavors.test2]
  cpu_count=2
```
