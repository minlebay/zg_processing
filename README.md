--- 
# Processing Service

The Processing Service is a part of the Message Generator and Router Project. It processes messages received from the router and integrates with Kafka for message brokering.

## Components

### Processing Service (`zg_processing`)
This component processes messages routed from the router.

#### Docker Compose Configuration
```yaml
version: '2'

networks:
  local-net:
    external: true

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:6.2.0
    hostname: zookeeper
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - local-net

  kafka:
    image: confluentinc/cp-kafka:6.2.0
    hostname: kafka
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "29092:29092"
      - "9092:9092"
      - "9101:9101"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      KAFKA_JMX_PORT: 9101
      KAFKA_JMX_HOSTNAME: localhost
    networks:
      - local-net

  zg_processing:
    build:
      context: .
      dockerfile: ./Dockerfile
    container_name: zg_processing
    env_file:
      - .env-docker
    networks:
      - local-net
    volumes:
      - .:/app
    depends_on:
      - kafka
    restart: unless-stopped

  kafdrop:
    image: obsidiandynamics/kafdrop
    container_name: kafdrop
    ports:
      - "9000:9000"
    environment:
      KAFKA_BROKERCONNECT: kafka:29092
      SERVER_SERVLET_CONTEXTPATH: "/"
    depends_on:
      - kafka
    networks:
      - local-net
```

#### Configuration File (`config.yaml`)
```yaml
grpc_server:
  listen_address: ${GRPC_SERVER_LISTEN_ADDRESS}

kafka:
  address: ${ZG_KAFKA_ADDRESS}
  group_id: group_id
  user: guest
  password: guest
  topic: processing_1

logstash:
  url: ${LOGSTASH_URL}
```

#### .env-docker File
```env
GRPC_SERVER_LISTEN_ADDRESS=zg_processing:50052
ZG_KAFKA_ADDRESS=kafka:29092
LOGSTASH_URL=http://logstash:5000
```

## Other Components

- **Message Generator**: Generates messages and sends them to the router.
- **Router**: Receives messages from the generator and routes them to processing servers.
- **Prometheus**: Monitors the application and collects metrics.
- **ELK Stack**: Collects and analyzes logs.
- **Grafana**: Visualizes the metrics collected by Prometheus.
- **Kafka**: A message broker that integrates with the backend.
- **Databases**: Includes MongoDB, MySQL, Redis for caching and indexing, and SQL/NoSQL repositories.

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Running the Processing Service
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/message-generator.git
   cd message-generator/processing
   ```
2. Build and run the Docker containers:
   ```bash
   docker-compose up --build
   ```

### Environment Variables
Ensure to set the following environment variables in the `.env-docker` file:
- `GRPC_SERVER_LISTEN_ADDRESS`: Address of the gRPC server (e.g., `zg_processing:50052`).
- `ZG_KAFKA_ADDRESS`: Address of the Kafka server (e.g., `kafka:29092`).
- `LOGSTASH_URL`: URL of the Logstash server (e.g., `http://logstash:5000`).

## License
This project is licensed under the MIT License.

---
