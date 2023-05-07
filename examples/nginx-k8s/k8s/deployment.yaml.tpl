apiVersion: apps/v1
kind: Deployment
metadata:
  name: '{{.appname}}'
  namespace: default
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
        created-by: '{{.author}}'
    spec:
      containers:
        - name: nginx-containers
          image: nginx
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: 100m
              memory: 128Mi
            requests:
              cpu: 100m
              memory: 128Mi

---
apiVersion: v1
kind: Namespace
metadata:
  name:  httpd
---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpd-deployment
  namespace: httpd
spec:
  selector:
    matchLabels:
      app: httpd
  replicas: {{.replicas}}
  template:
    metadata:
      labels:
        app: httpd
    spec:
      containers:
        - name: httpd-containers
          image: httpd
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: 100m
              memory: 128Mi
            requests:
              cpu: 100m
              memory: 128Mi
