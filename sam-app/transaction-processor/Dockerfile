FROM public.ecr.aws/docker/library/golang:1.23 as build-image
WORKDIR /src
COPY . .
RUN go build -o lambda-handler
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build-image /src/lambda-handler .
# The transactions.csv file should be copied into the transaction-processor directory before building
# or provided at runtime through environment configuration
COPY transactions.csv .
ENTRYPOINT ./lambda-handler
