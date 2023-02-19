docker run -it --restart always --name zookeeper -p 12181:2181 -d wurstmeister/zookeeper:latest

docker run -it --restart always --name kafka01 -p 19092:9092 -d -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=172.20.16.20:12181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://172.20.16.20:19092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 wurstmeister/kafka:latest

docker run -it --restart always --name kafka02 -p 19093:9092 -d -e KAFKA_BROKER_ID=1 -e KAFKA_ZOOKEEPER_CONNECT=172.20.16.20:12181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://172.20.16.20:19093 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 wurstmeister/kafka:latest

docker run -it --restart always --name kafka03 -p 19094:9092 -d -e KAFKA_BROKER_ID=2 -e KAFKA_ZOOKEEPER_CONNECT=172.20.16.20:12181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://172.20.16.20:19094 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 wurstmeister/kafka:latest
