Use this proxy for bridging HTTP/1 - HTTP/2

docker run --rm -it -p 50052:50052 ghcr.io/mirkolenz/grpc-proxy:latest --proxy-port 50052 --backend-port 50051