kind: Service
apiVersion: v1
metadata:
  name: {{.appname}}-service
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
    - name: 80-80-tcp
      port: 80
      targetPort: 8080
