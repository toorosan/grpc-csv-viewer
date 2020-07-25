Helpful configuration for UI development needs (before gRPC client will not start serving HTML):
- `cd grpc-csv-viewer/internal/app/roles/client/ui && yarn run dev`
- `cd grpc-csv-viewer/cmd/client && go run main.go`
- `cd grpc-csv-viewer/misc && docker-compose -f nginx-dev-docker-compose.yml up`

Then open http://127.0.0.1:8888/ to check how UI is working with back-end.