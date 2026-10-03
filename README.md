# Price Board

Tracks crypto prices, stores history, shows a live board. Runs on Kubernetes; all infra in code.

> 🚧 Work in progress

## Stack

- **App:** Go (`api` + `fetcher`), Postgres
- **Infra:** kind, Terraform, Kustomize
- **Delivery:** GitHub Actions, Argo CD
- **Ops:** Prometheus, Grafana, Sealed Secrets

## License

[MIT](LICENSE)
