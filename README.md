# iksctl

`iksctl` is a command-line tool containing scripts used in the EKS Anywhere tooling.

## Quick Start

### Run Without Building
To use `iksctl` without building, navigate to the `/eks-anywhere-tooling/iksctl/` directory in your terminal and run:
```bash
./bin/iksctl help
```

### Build and Install

#### Prerequisites
We use Taskfile for workflows related to `iksctl`. Follow the [Taskfile installation instructions](https://taskfile.dev/installation/).

#### Build the iksctl Binary
Run the following command to build the `iksctl` binary, which will be output to the `bin` folder:
```bash
task build
```

#### Install the Binary
1. Remove any existing version of `iksctl`:
    ```bash
    sudo mv /usr/local/bin/iksctl /usr/local/bin/iksctl.bak
    ```
2. Copy the new binary to the appropriate location:
    ```bash
    sudo cp ~/projects/itm-iks/eks-anywhere-tooling/iksctl/bin/iksctl /usr/local/bin/iksctl
    ```
3. Ensure the binary is executable:
    ```bash
    sudo chmod +x /usr/local/bin/iksctl
    ```
4. Verify the installation:
    ```bash
    which iksctl
    iksctl help
    ```

### Docker

#### Build Docker Image
To build a Docker image for `iksctl`, run:
```bash
docker build -t iksctl .
```

#### Run with Docker
You can run `iksctl` commands using Docker by mounting your configuration files and examples:

```bash
# Example: Parse and validate a cluster manifest
docker run --rm -v "$(pwd)/examples:/app/examples" iks-harbor-base.systematicgroup.local/iks-library/iksctl:<iksctl-version> parse validate -f /app/examples/cluster-manifest.yaml --host="testprojectcloud.systematicgroup.local" --username='adm-jwn@systematicgroup.local' --password=your-password

# Example: Run airgap command with config file
docker run --rm -v "$(pwd)/configs:/app/configs" -e IKS_PWDB_APIKEY=your-api-key iks-harbor-base.systematicgroup.local/iks-library/iksctl:<iksctl-version> airgap --dry-run -c /app/configs/airgap/airgapconfig-cilium.yaml

# Example: Get help
docker run --rm iksctl help
```

#### Push to Registry
To push the Docker image to a registry:
```bash
# Tag the image
docker tag iksctl iks-harbor-base.systematicgroup.local/iks-library/iksctl:<iksctl-version>

# Push to registry
docker push iks-harbor-base.systematicgroup.local/iks-library/iksctl:<iksctl-version>
```

## Usage

### Overview
`iksctl` supports two primary use cases:
1. **airgap**: Airgap a Helm chart based on a configuration file to a local registry.
2. **pwdb**: Fetch credentials from a password database.
3. **ntnx**: Run validation tests for Nutanix clusters.

### Global Environment Variable
Set the `IKS_PWDB_APIKEY` environment variable to the desired API key:
```bash
export IKS_PWDB_APIKEY=your-pwdb-api-key-123456789
```

### PWDB Example
To fetch a password, use the following command:
```bash
iksctl pwdb password get notes 1234 --endpoint https://pwdb --base64-decode
```

#### Working Example
Using an API key from https://pwdb/plid=1655 ("Global Business Services/ITM/ITM Infrastructure K8S"), you can retrieve a username:

```bash
export IKS_PWDB_APIKEY=REDACTED
go run ./cmd/iksctl pwdb password get username 25667
```
Example output:
```
gbs-work01-dsm-iks-kubeconfig
```

### Airgap Example
#### Determine Helm Chart Version
Add the Helm repository and search for the desired chart:
```bash
helm repo add dell-csi https://dell.github.io/helm-charts
helm repo update
helm search repo dell-csi
```
Example output:
```
NAME                                            CHART VERSION   APP VERSION     DESCRIPTION
dell-csi/csi-powerstore                         2.13.0          2.13.0          PowerStore CSI (Container Storage Interface) driver
```

#### Create Configuration File
Specify the chart details in a configuration file:
```bash
cat <<EOF > airgap.yaml
apiVersion: iksctl.kubematic.io/v1alpha1
kind: AirgapConfig
ociRepo:
     host: iks-harbor-base.systematicgroup.local
     repository: iks-library
     pwdbCredentialID: 23114
charts:
     - repo: "https://dell.github.io/helm-charts"
        name: "csi-powerstore"
        version: "2.13.0"
EOF
```

Cilium example:
Check helm version https://artifacthub.io/packages/helm/cilium/cilium
- set the pwdb key
```bash
export IKS_PWDB_APIKEY=REDACTED (use the pwdb list API key from https://pwdb/plid=1655)
```
- switch the to eks-anywhere-tooling/iksctl/configs/airgap
```bash
cd /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/configs/airgap
```
- create the config file

```bash
cat <<EOF > airgapconfig-cilium.yaml
apiVersion: iksctl.kubematic.io/v1alpha1
kind: AirgapConfig
ociRepo:
  host: iks-harbor-base.systematicgroup.local
  repository: isovalent/cilium
  pwdbCredentialID: 23114
charts:
  - repo: "https://helm.cilium.io/"
    name: "cilium"
    version: "1.17.4"
EOF
```
- switch to eks-anywhere-tooling/iksctl dir and run the dry-run, and if successful, run without dry-run
```bash
cd /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl 
go run ./cmd/iksctl airgap --dry-run -c /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/configs/airgap/airgapconfig-cilium.yaml
go run ./cmd/iksctl airgap -c /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/configs/airgap/airgapconfig-cilium.yaml
```

#### Run Airgap in Dry-Run Mode
Execute the airgap process in dry-run mode to preview actions:
```bash
go run ./cmd/iksctl airgap --dry-run -c /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/examples/api/airgapconfig-csi.yaml
```

Example output:
```
INFO[0000] Using config file                             file=/home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/examples/api/airgapconfig-csi.yaml
INFO[0000] Dry run is enabled. Printing actions without doing them.
INFO[0000] Fetching oci credentials from PWDB
INFO[0000] Pulling chart                                 chart=csi-powerstore repo="https://dell.github.io/helm-charts" version=2.13.0
INFO[0003] helm push                                     chart=csi-powerstore repo="https://dell.github.io/helm-charts" tarball=/tmp/airgap-helm-charts3880473328/csi-powerstore-2.13.0.tgz targetRepo="oci://iks-harbor-base.systematicgroup.local/iks-library" version=2.13.0
INFO[0003] templating chart in order to look for images  chart=csi-powerstore repo="https://dell.github.io/helm-charts" version=2.13.0
INFO[0003] extracting images from chart                  chart=csi-powerstore repo="https://dell.github.io/helm-charts" version=2.13.0
INFO[0003] copying image                                 srcRef="registry.k8s.io/sig-storage/csi-resizer:v1.13.1" tgtRef="iks-harbor-base.systematicgroup.local/iks-library/csi-powerstore/csi-resizer:v1.13.1"
...
```

Once the dry-run is successful, execute the airgap process without dry-run:
```bash
go run ./cmd/iksctl airgap -c /home/mtp/projects/itm-iks/eks-anywhere-tooling/iksctl/examples/api/airgapconfig-csi.yaml
```

### Location for Non-Example Airgap Configuration Files

For non-example `airgapconfig.yaml` files, place them in the `/eks-anywhere-tooling/iksctl/configs/airgap/` 'directory within the project to keep them organized and separate from example files. 


### onboard currently used helmcharts
This will get the repo added, updates your local cache, and searches for charts in it.
```bash
helm repo add dell-csi https://dell.github.io/helm-charts
helm repo update
helm search repo dell-csi

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm search repo ingress-nginx

helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
helm repo update
helm search repo sealed-secrets

helm repo add cert-manager https://charts.jetstack.io
helm repo update
helm search repo cert-manager

helm repo add metallb https://metallb.github.io/metallb
helm repo update
helm search repo metallb

helm repo add projectsveltos https://projectsveltos.github.io/helm-charts
helm repo update
helm search repo projectsveltos

helm repo add velero https://vmware-tanzu.github.io/helm-charts
helm repo update
helm search repo velero

helm search repo bitnami/kubectl
```


#### Futher helm commands you may use
Inspect the chart
```bash
helm show values ingress-nginx/ingress-nginx > ingress-values.yaml
```
Preview the manifests without applying them
```bash
helm template ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx > ingress-manifests.yaml
```
Upgrade or reconfigure later
```bash
helm upgrade ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx -f new-values.yaml
```

Uninstall
```bash
helm uninstall ingress-nginx --namespace ingress-nginx
```

### Ntnx validations

The following sub-command will run the validations against individual nutanix resources.

```bash
iksctl ntnx
```

### Ntnx parse

Will validate a cluster configuration file and parse the nutanix resources from it.

```bash
iksctl parse validate -f /path/to/cluster-config.yaml
```
