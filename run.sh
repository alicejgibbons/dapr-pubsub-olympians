# Run subscriber

dapr run --app-id sub \
         --app-protocol http \
         --app-port 8080 \
         --dapr-http-port 3500 \
         --log-level debug \
         --resources-path ./config \
         go run sub/sub.go

# Run publisher

dapr run --app-id pub \
         --log-level debug \
         --resources-path ./config \
         go run pub/pub.go