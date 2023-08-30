# terraform-cloud-deploy

> Deploy a URL shortener application with Terraform on Cloud Run
>
## Code

To test the code you should have postgres instance and update dns url with the postgres instance

1. Install go packages :

```bash
cd url-shortener && go get

```

2. Run application locally🤞:

```bash
go run main.go
```

If application is running successfully you can create a shortened url with the request:

```bash
curl --location 'http://localhost:8090/create' \
--form 'url="https://www.google.com/"'
```

## Deploy

You will need to have a GCP project with billing enable to proceed from here.

Set up gcloud SDK locally and set the project created

1. Build docker image
Still in `url-shortener` folder in the terminal:
This step can be skipped if you have an image available in a docker image registry of your choice.

```bash
gcloud builds submit --region=us-central1 --tag us-central1-docker.pkg.dev/REPOSITORY/docker/url-shortener:latest
```

2. Change directory to infra folder and update `PROJECT_ID` with your project id from GCP

3. Initialize the Terraform configuration

```bash
terraform init
```

4. Preview Terraform changes

```bash
terraform plan
```

5. Apply changes

```bash
terraform apply
```

6. Destroy infrastructure

```bash
terraform destroy
```

## License

MIT 2023
