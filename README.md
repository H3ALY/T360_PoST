Transfer360 Proof Of Skill Test
---

## Prerequisites
- Go 1.24 or later.
#### Using Pub/Sub Emulator
**Requires Setup/Install of Python Dev Environment, JDK, Google Cloud CLI**
- Setup Emulator (https://cloud.google.com/pubsub/docs/emulator)
#### Using Google Cloud
- A Google project with Pub/Sub enabled.
- Service account with Pub/Sub Publisher Role and associated clientId/clientKey

## Setup

1. Clone the repository

```sh
git clone https://https://github.com/H3ALY/T360_PoST
```

2. Install dependencies
```sh
cd T360_PoST
```
```sh
go mod tidy
```
3. Set the required environment variables in Config.yaml
- <mark>server.port</mark> Default set is 8888 to avoid potentially already running services.
- <mark>google.usingCloud</mark> (Set True if using Google's Pub/Sub or False if using Google's Pub Sub Emulator)
- <mark>google.pubSubTopic</mark> (Pub/Sub Topic)
- <mark>local_emulator.projectId</mark> (Your Pub/Sub Project ID)
- <mark>local_emulator.pubSubTopic</mark> Pub/Sub Topic

3a. If using a Google's Cloud Pub/Sub enter the service account details in Config/service-account-key.json which is generated when creating your service account.
  
4. Run the Application
```sh
go run App\Main.go
```

## Local Emulator

Run the local emulator with the following .bat file
```bat
@echo off
rem Set environment variable for the emulator
set PUBSUB_EMULATOR_HOST=localhost:8085

rem Wait for the emulator to fully start (add a delay)
timeout /t 5 /nobreak

rem Keep the script running to allow interaction with the emulator
pause
```
- Open a new cmd prompt and run

```cmd
set CLOUDSDK_CORE_PROJECT=transfer360
```

```cmd
gcloud config get-value auth/disable_credentials
```
If it does not return True, run:
```cmd
gcloud config set auth/disable_credentials true
```
Create aSubscription aand Topic with -
```cmd
curl -X PUT http://localhost:8085/v1/projects/transfer360/subscriptions/my-sub -H "Content-Type: application/json" -d "{ \"topic\": \"projects/transfer360/topics/positive_searches\" }"
```
You can confirm creation with - 
```cmd
curl http://localhost:8085/v1/projects/transfer360/subscriptions
```
When you've published some messages you will be able to retrieve them with -
```cmd
curl -X POST http://localhost:8085/v1/projects/transfer360/subscriptions/my-sub:pull -H "Content-Type: application/json" -d "{ \"maxMessages\": 10 }"
```
## Principles
### SOLID

***S***RP (Single Responsibility Principle)

- Client.go handles API requests.
- Publisher.go Handles message publishing
- SearchService.go processes responses and integrates components

***O***CP (Open/Closed Principle)

- If we want to add a new endpoint, we don't need to change the SearchService logic.
- Allows for retry logic/batching within the Publisher

***L***SP (Liskov Substitution Principle)

- If we want to use another object storage instead of Pub/Sub such as RabbitMQ/MQTT only Publisher.go needs to be replaced.

***I***SP (Interface Segregation Principle)

- Separte interfaces for API and Pub/Sub logic making components independant

***D***IP (Dependency Inversion Principle)

- Instead of hardcoding within the service they're injected as interfaces
- HandleRequest works on an abstracted publisher

### RED

***R***esilient

- If an API Files or times out it does not crash the syste,
- Uses timeouts to ignore slow responses
- Single client for PubSub is used and allows for pooling

***E***lastic
- Uses goroutines for concurrent API calls
- Can easily scale to handle more APIs without modifying core logic

***D***urable
- Uses Google Pub/sub to ensure messages are persisted and retired if processing files
- Retries failed messages automatically.
- Pub/Sub's batching and connection pooling improve performance