# ecoctl

This is CLI tool for interaction with the eco-server.

## Usage

### Setting
Firstly set these ENVs 
```asciidoc
ECO_URL=url to eco-server
ECO_TOKEN=token to interact
```

### Build 

Build ecoctl with the following command:
```asciidoc
make build
```
After that there is ecoctl in bin dir

### Optional
You can move ecoctl in dir which is in the env PATH

### Create cluster

Before you need to specify cluster configuration in yaml format.

Then to create cluster use following command:
```asciidoc
ecoctl cluster create -p examples/cluster.yaml 
```

### Watch cluster

To wait until cluster(for example 1) is ready, use this command:
```asciidoc
ecoctl cluster watch 1
```

### Get cluster

To get cluster use:
```asciidoc
ecoctl cluster get 1
```

### Get cluster config

To get cluster config use:
```asciidoc
ecoctl cluster config 1 -p ~/.kube/config -o yaml
```

### Delete cluster

To delete cluster use:
```asciidoc
ecoctl cluster delete 1
```

### Create pool

To update pool use:
```asciidoc
ecoctl pool create -c 1 -p examples/pool.yaml
```

### Update pool

To update pool use:
```asciidoc
ecoctl pool update 1 --count 1 --max 2 --min 1...etc
```

For additional info use:
```asciidoc
ecoctl pool update -h
```

### Delete pool

```asciidoc
ecoctl pool delete 1
```