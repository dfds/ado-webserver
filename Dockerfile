FROM mcr.microsoft.com/powershell:lts-alpine-3.10

# Install dependencies
RUN pwsh -c "Install-Module Pode -Force; Install-Module AWSPowerShell.NetCore -Force"

RUN mkdir /webserver

COPY src/ /webserver/

WORKDIR /webserver

EXPOSE 8080

ENTRYPOINT ["pwsh", "-File", "./Webserver.ps1"]