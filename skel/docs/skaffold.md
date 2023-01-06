# Skel Skaffold Integration

## Generate Project Skaffold Support

A `skaffold.yaml` file will be created in the root of any app or service pack project you create using `skel`, no extra action is required.

In addition, `deployments/kubernetes/minivms/${app.name}-deployment.yaml` and `deployments/kubernetes/msxlite/${app.name}-deployment.yaml` will be created.

If your project has already been generated, you can use the `skel generate-kubernetes` command inside the project folder to add the skaffold support files.

## Configure Kubernetes

1. Connect to your msx-lite instance (the kubernetes host, not the installer), and retrieve your kubernetes configuration:

```
kubectl config view --raw --minify
```

2. Apply this configuration as either the default kubernetes configuration, or as a custom configuration referred to by the `KUBECONFIG` environment variable:

```
mkdir -p $HOME/.kube
cat > $HOME/.kube/config <<EOF
<config contents from instance go here>
EOF
```

3. Update the server URL in the kubeconfig file you just saved to refer to the lab IP address:

```
#    server: https://127.0.0.1:6443 
     server: https://10.81.85.174:6443
```

4. If using a non-default config file, ensure `KUBECONFIG` is set in your bash profile to point to the new file:

```
export KUBECONFIG=$HOME/.kube/rtp-4-msx-lite-35/config
```

## Setup Skaffold Support in GoLand

To set up `skaffold` in GoLand:

1. Install Skaffold: Follow the instructions at https://skaffold.dev/docs/install/ to install Skaffold 2.x or higher on your system.

2. Install the Skaffold plugin for GoLand: In GoLand, go to `Goland | Settings | Plugins...`, search for "Cloud Code", and click the Install button

3. You *may* need to restart GoLand

## Create an MSX-Lite Run Configuration

1. When you open your generated project in GoLand you should now see a popup saying "Kubernetes with Cloud Code. Skaffold configuration detected" since there will be a skaffold.yaml in the root

2. Via the `Add Configuration` link therein, or via the light blue `Add Configuration` button top right, or via `Run | Edit Configurations | +` add a run config

3. Select the config type: "Cloud Code: Kubernetes"

4. Give the configuration a name

5. On the run tab, `Environment Variables` specify:  

    `SKAFFOLD_PROFILE=msxlite`  

    this tells skaffold to use the msxlite deployment found in the msxlite subdir 

6. Give the path to the `skaffold.yaml` file on the "Build | Deploy" tab (it should default correctly)

7. You probably want "All Modules and Dependencies" selected 

8. Now you can run that config to deploy using skaffold 
