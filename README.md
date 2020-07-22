# kubectl-emitevent

`kubectl-emitevent` is a [kubectl plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that emits event for requested object.

# Usage

<details><summary>start the minikube cluster (skip if you are using an existing cluster) </summary>
<p>

```bash
➜  kubectl-emitevent git:(master) minikube start
😄  minikube v1.10.1 on Darwin 10.15.5
✨  Using the hyperkit driver based on existing profile
👍  Starting control plane node minikube in cluster minikube
🔄  Restarting existing hyperkit VM for "minikube" ...
🎉  minikube 1.12.1 is available! Download it: https://github.com/kubernetes/minikube/releases/tag/v1.12.1
💡  To disable this notice, run: 'minikube config set WantUpdateNotification false'

🐳  Preparing Kubernetes v1.18.2 on Docker 19.03.8 ...
🌟  Enabled addons: default-storageclass, ingress, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube"
```

</p>
</details>


Run `kubectl-emitevent` daemonset/kube-proxy -n kube-system --reason "foo-reason" --message "bar-message"

```bash
  ## emit event
  ➜  kubectl-emitevent daemonset/kube-proxy -n kube-system --reason "foo-reason" --message "bar-message"

  ## verify event
  ➜  kubectl describe daemonset/kube-proxy -n kube-system

Name:           kube-proxy
Selector:       k8s-app=kube-proxy
Node-Selector:  beta.kubernetes.io/os=linux
Labels:         k8s-app=kube-proxy
Annotations:    deprecated.daemonset.template.generation: 1
.
.
.
.
.
Events:
  Type    Reason               Age    From               Message
  ----    ------               ----   ----               -------
  Normal  foo-reason           4s    kubectl-emitevent  bar-message

```

