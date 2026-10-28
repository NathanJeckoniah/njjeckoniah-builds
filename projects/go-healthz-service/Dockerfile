 #names the foundation where the image will be build from, yaani files to be packaged
FROM golang:1.27-alpine AS build
WORKDIR /src

#copies the go.mod file from my compiter into src, the workshop
COPY go.mod ./
# RUN is supposed to download libraries that I will need, 
# the ones my progam uses but mine needs none other so it will run and finish instantly
RUN go mod download
#copies my code files
COPY *.go ./

#compiles the program that depends on nothing else, so it can run in the box we have created
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service .


# starts a second almost empty image and copies the program into it
# the workshop gets tossed, so the uploaded image is small. 
# Also if the program was huge it gives the attacker nothing to use
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/service /service


# notes that the user uses port 8080, and runsthe program as an ordinary user
USER nonroot:nonroot
EXPOSE 8080


# docker checks the program with its own health checks and declares it unhealthy in three consecutive fails
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 CMD ["/service", "-healthcheck"]
ENTRYPOINT ["/service"]
