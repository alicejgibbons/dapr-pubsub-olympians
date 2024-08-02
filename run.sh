# Run subscriber
dapr run --app-id olympians-sub \
         --app-protocol http \
         --app-port 8080 \
         --dapr-http-port 3500 \
         --resources-path ../config \
         go run sub.go

# Run publisher
dapr run --app-id olympians-pub \
         --resources-path ../config \
         go run pub.go