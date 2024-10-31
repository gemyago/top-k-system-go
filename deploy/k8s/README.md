# k8s deployment tools

This folder includes k8s deployment tools:
- kubectl
- helm
- [helmfile](https://github.com/helmfile/helmfile)

## Usage for local development

Assuming [direnv](https://github.com/direnv/direnv) is installed and properly configured.

Install tools:
```bash
make tools
```

Deploying locally:
```bash
# Check the diff before applying
make diff

# Apply the changes
make apply
```

To deploy to a specific environment please initialize `kubectl` first.
```bash
# Show available contexts
kubectl config get-contexts

# If you already have the context, switch to it
kubectl config use-context "<context-name>"

# Otherwise you can get the credentials for specific cluster and initialize the context
# This is provider specific.

# Run cluster-info to check if you're pointing to the right cluster
# After running it make sure control plane is pointing to the clusters public IP address
kubectl cluster-info
```
