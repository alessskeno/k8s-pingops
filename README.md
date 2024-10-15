 # 🖥️ k8s-pingops: Kubernetes Connectivity Tester

 This repository contains a simple telnet and DNS client tool, containerized for use in Kubernetes environments. It helps developers and operations teams test connectivity between their Kubernetes clusters and external resources like databases, third-party services, or internal APIs.

![image](https://github.com/user-attachments/assets/7751eda7-88d1-420a-8cd6-b78eff0b5c4f)


 ## Key Features
 - **Telnet and DNS Client**: Tools to test network connectivity and DNS resolution from within Kubernetes clusters.
 - **Kubernetes Ready**: Easily deployable as a Kubernetes service.
 - **Lightweight**: Minimal resource usage, allowing for testing without heavy footprint.
 - **Fully Containerized**: Available as a Docker container, making it easy to integrate with your CI/CD pipelines.

 ## 🔧 Setup & Installation



 ### Prerequisites
 1. **Kubernetes Cluster**: Make sure you have a running Kubernetes cluster.
 2. **kubectl CLI**: To interact with your Kubernetes cluster.
 3. **Terraform**: For those using the `terraform-kubernetes-provider` to manage the cluster resources.
 4. **Docker**: (Optional) For local testing before pushing the container to your cluster.

 ### Docker Hub
 The image is available on DockerHub:
```bash
docker pull aleskerov/k8s-pingops:latest
```

 ## 🚀 Usage & Deployment

### Deploying to Kubernetes

 1. **Kubernetes YAML Manifest**
    If you're not using Terraform, here’s a YAML version of the Kubernetes manifest:

<details>
<summary><strong>Click to Expand YAML Manifest</strong></summary>

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    k8s-balancer: "true"
  name: k8s-pingops
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: k8s-pingops
  name: k8s-pingops
  namespace: k8s-pingops
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: k8s-pingops
  template:
    metadata:
      labels:
        app.kubernetes.io/name: k8s-pingops
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app.kubernetes.io/name: k8s-pingops
                topologyKey: kubernetes.io/hostname
              weight: 1
      containers:
        - image: aleskerov/k8s-pingops:latest
          imagePullPolicy: IfNotPresent
          name: k8s-pingops
          ports:
            - containerPort: 8080
          resources:
            limits:
              cpu: 150m
              memory: 128Mi
            requests:
              cpu: 10m
              memory: 10Mi
---
apiVersion: v1
kind: Service
metadata:
  name: k8s-pingops
  namespace: k8s-pingops
spec:
  ports:
    - name: http
      port: 80
      targetPort: 8080
  selector:
    app.kubernetes.io/name: k8s-pingops
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8s-pingops
  namespace: k8s-pingops
spec:
  ingressClassName: nginx
  rules:
    - host: pingops.yourdomain.com
      http:
        paths:
          - backend:
              service:
                name: k8s-pingops
                port:
                  number: 80
            path: /
            pathType: Prefix
```
</details>



2. **Apply Kubernetes YAML**
   Once you have the manifest, apply it to your Kubernetes cluster:

   ```bash
   kubectl apply -f k8s-pingops.yaml
   ```

 3. **Access the Service**
    The service is now available at `http://pingops.yourdomain.com`. You can test connectivity using the telnet and DNS client tools.

 ## 📜 Terraform Kubernetes Provider (Optional)

 If you prefer to use Terraform to manage Kubernetes resources, below is the Terraform configuration for deploying `k8s-pingops` in your cluster.

 ### Terraform Resources

 The following resources will be created:
 - **Namespace**: `k8s-pingops`
 - **Deployment**: A 2-replica deployment for the ping client.
 - **Service**: ClusterIP service exposing the container on port 80.
 - **Ingress**: Ingress route (using NGINX ingress controller).

 You can copy the provided Terraform code into your own `.tf` file and apply it to your Kubernetes environment.

 ```terraform
 resource "kubernetes_namespace" "k8s_pingops" {
   metadata {
     name = "k8s-pingops"
   }
 }

 resource "kubernetes_deployment" "k8s_pingops" {
   metadata {
     name      = "k8s-pingops"
     namespace = kubernetes_namespace.k8s_pingops.metadata[0].name
   }
   spec {
     replicas = 2
     selector {
       match_labels = {
         "app.kubernetes.io/name" = "k8s-pingops"
       }
     }
     template {
       metadata {
         labels = {
           "app.kubernetes.io/name" = "k8s-pingops"
         }
       }
       spec {
         container {
           name  = "k8s-pingops"
           image = "aleskerov/k8s-pingops:latest"
         }
       }
     }
   }
 }

 resource "kubernetes_service" "k8s_pingops" {
   metadata {
     name = "k8s-pingops"
   }
   spec {
     selector = {
       "app.kubernetes.io/name" = "k8s-pingops"
     }
     port {
       port = 80
     }
   }
 }
 ```

 To deploy with Terraform:
 ```bash
 terraform init
 terraform apply
 ```

 ## 📊 Resource Limits

 - **CPU**: Requests `10m`, Limits `150m`
 - **Memory**: Requests `10Mi`, Limits `128Mi`

 Make sure these values align with your cluster's capacity.

 ## 🛠️ Configuration Options

 You can modify the configuration of this app by editing:
 - **Replicas**: Change the `replicas` field in the Kubernetes manifest or Terraform file.
 - **Resource Limits**: Modify CPU and memory requests/limits as needed.

 ## 👥 Contributions

 Feel free to open an issue or submit a pull request if you have any improvements or suggestions!

 ## 🧑‍💻 Author

 This project was developed by [Kanan Alasgarli](https://github.com/alessskeno). If you like the project, please ⭐ it!

 ## 📄 License

 This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
