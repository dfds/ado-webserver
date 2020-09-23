# This file imports the Pode webserver module, imports functions from your own module and starts a webserver on the given port.
# Routes are added which makes the webserver react on requests for specific endpoints.
# Docs for routes: https://badgerati.github.io/Pode/Tutorials/Routes/Overview/

Import-Module Pode

# Sample request againsat the server:
# Invoke-WebRequest -Uri 'http://localhost:8080/builds' -Method Post -Body '{ "project": "DevelopmentExcellence" }' -ContentType 'application/json'

Start-PodeServer {
    # Enable Open API and Swagger at the endpoint /docs/swagger
    Enable-PodeOpenApi -Title 'My Webservice API' -Version 1.0.0
    Enable-PodeOpenApiViewer -Type Swagger -Path '/docs/swagger' -DarkMode

    # Import our own custom modules.
    Import-PodeModule -Path "MyModule.psm1"

    # Add endpoint the server will listen on.
    Add-PodeEndpoint -Address localhost -Port 8080 -Protocol Http

    # Add routes and define respones to the routes.
    Add-PodeRoute -Method Post -Path '/builds' -ScriptBlock {
        param($request)
        $AllProjectBuilds = Get-ProjectBuilds -Organization dfds -Project $request.Data.project -BasicAuthHeader (Get-BasicAuthHeader)
        $SortedBuilds = Get-LatestBuilds -Builds $AllProjectBuilds | Get-SortedBuilds

        Write-PodeJsonResponse -Value $SortedBuilds

    }
}