# OpenAPI-Compatible DNSDB Proxy

## About

This is a proxy that converts the output of the DNSDB 2.0 API from the
Farsight Streaming Application Framework (SAF) to an `application/json`
object. You would want to use this in an application framework that is
does not need to stream input and is capable of buffering entire result
sets into memory. We have provided an [OpenAPI 3.0 Specification](api.yaml)
that you can use to query this service and generate application code.

## API Documentation

You may browse the
[API Documentation](https://app.swaggerhub.com/apis-docs/hstern/dnsdb/2.0.0)
and run queries on Swagger Hub, or download the [OpenAPI 3.0 Specification](api.yaml)
and use it directly.

The endpoints and parameters are functionally identical to those of
[DNSDB API 2.0](https://docs.dnsdb.info) with Flexible Search. The output
of each endpoint is a JSON object as described in the OpenAPI Specification.

If an error occurs upstream, the server will respond with a status code of 500.
If your result set hits a limit the server will send a `Limited: true` header. If
you are using a HTTP/2 client the server will send a `Success: true` trailer to
indicate that the query was successful.
