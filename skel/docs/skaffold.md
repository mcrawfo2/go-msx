# Skel Skaffold Integration

## Generating

A `skaffold.yaml` file will be created in the root of any app or service pack project you create using `skel`, no extra action is required.

In addition, `deployments/kubernetes/minivms/${app.name}-deployment.yaml` and `deployments/kubernetes/msxlite/${app.name}-deployment.yaml` will be created.

## Setup Skaffold Support in GoLand

To set up `skaffold` in GoLand:

1. Install Skaffold: Follow the instructions at https://skaffold.dev/docs/install/ to install Skaffold on your system.

2. Install the Skaffold plugin for GoLand: In GoLand, go to `Goland | Settings | Plugins...`, search for "Cloud Code", and click the Install button

3. You *may* need to restart GoLand

## Create an MSXLite Run Configuration

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