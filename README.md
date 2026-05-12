# Ledger system

```mermaid
C4Context
    Boundary(b0, "Identity", "Non financial data related to users and their accounts") {
        System(SystemA, "User", "Data about users")
        System(SystemB, "Account", "Accounts that users have")
    }
    
    Boundary(b1, "Ledger", "Financial data") {
        System(SystemC, "Ledger", "Source of truth of financial data")
    }

    Boundary(b2, "Transfer", "Metadata for financial data") {
        System(SystemD, "Transfer", "Attempts at changing financial data")
    }

    Rel(SystemA, SystemB, "Get accounts from user", "")
    Rel(SystemD, SystemA, "Fetch account data", "(TBD)")
    Rel(SystemD, SystemC, "Update financial data", "")
    Rel(SystemC, SystemB, "Convert id into something usable", "(TBD)")
```

## Kubernetes
Commands below will allow you to setup and run.
Use kubernetes over just deploying apps manually because if an app crashes or duplicates need to be spun up, it will require lots of manual intervention whereas kubernetes will rerun the app if it crashes and if the number of replicas need to change automatically, this can be done (so during peak periods, more replicas can be spun up).
```bash
kubectl apply -f ./k8s/
kubectl port-forward svc/go-app-service 8080:80
```

## Helm
Will need to wait until database is up and live before attempting to make any queries.
Helm is preferred over just using kubernetes because it can be used across environments instead of copying and pasting YAML across local, demo and live.
```bash
helm install go-app-service ./lambda-system-prac-workflow
kubectl port-forward svc/go-app-service-lambda-system-prac-workflow 8080:8080
```
