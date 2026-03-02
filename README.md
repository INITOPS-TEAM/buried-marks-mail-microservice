# Mail Microservice

Mail Service is a Go-based microservice responsible for generating daily secret access codes, temporarily storing them in memory (Valkey), and dispatching bulk emails to users via the Brevo API.

## Features

* Scheduler (Cron): A background job runs daily (or at a configured testing interval).
* Generation & Storage: The service generates a unique token and saves it to **Valkey** with a strict 24-hour Time-To-Live (TTL).
* User Base Fetching: It makes an internal HTTP request to the `auth-service`  to retrieve a list of active user email addresses.
* Asynchronous Dispatch: Leveraging Go's concurrency (Goroutines), the service sends emails in parallel to each user via the Brevo API.

## Environment variables (.env)
 Create a '.env' file in the root directory based on example.env