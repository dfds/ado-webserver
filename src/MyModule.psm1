$Global:APIVersion = "api-version=6.1-preview.6"

function Get-BasicAuthHeader {
    [CmdletBinding()]
    [OutputType([psobject])]
    param (
        [Parameter(Mandatory = $false, ValueFromPipelineByPropertyName = $true, Position = 0)]
        [string]$PersonalAccessToken
    )
    process {

        $PersonalAccessToken = ''
        $Token = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes(":$($PersonalAccessToken)"))
        $BasicAuthHeader = @{Authorization = "Basic $Token"}
        
        return $BasicAuthHeader
    }
}

function Get-ProjectBuilds {
    [CmdletBinding()]
    [OutputType([psobject])]
    param (
        [Parameter(Mandatory = $true, ValueFromPipelineByPropertyName = $true, Position = 0)]
        [string]$Organization,

        [Parameter(Mandatory = $true, ValueFromPipelineByPropertyName = $true, Position = 1)]
        [string]$Project,

        [Parameter(Mandatory = $true, ValueFromPipelineByPropertyName = $true, Position = 2)]
        [System.Collections.IDictionary]$BasicAuthHeader
    )
    process {

        $URL = "https://dev.azure.com/$Organization/$Project/_apis/build/builds?$Global:APIVersion"

        $Builds = (Invoke-RestMethod -Uri $URL -Method Get -ContentType "application/json" -Headers $BasicAuthHeader).value

        return $Builds
    }
}


function Get-LatestBuilds {
    [CmdletBinding()]
    [OutputType([psobject])]
    param (
        [Parameter(Mandatory = $true, Position = 0)]
        [Object[]]$Builds
    )
    process {

        $LatestBuilds = $Builds | Sort-Object -Property {$_.definition.name} -Unique | Sort-Object -Property queuetime -Descending
        
        return $LatestBuilds
    }
}
function Get-SortedBuilds {
    [CmdletBinding()]
    [OutputType([psobject])]
    param (
        [Parameter(Mandatory = $true, ValueFromPipeline = $true, Position = 0)]
        [Object[]]$Builds
    )
    process {

        $SortedBuild = $Builds | Select-Object status, result, buildNumber, queueTime, startTime, finishTime, sourceBranch, @{Name = "’pipelineName’"; Expression = {$_.definition.name}}, @{Name = "projectID"; Expression = {$_.project.id}}
        
        # A slightly hacky way to generate the direct link to the build run.
        $BuildPageLink = "https://dev.azure.com/dfds/$($SortedBuild.projectID)/_build/results?buildId=$($SortedBuild.buildNumber)&view=results"
        $SortedBuild | Add-Member -NotePropertyName BuildPageLink -NotePropertyValue $BuildPageLink
        
        return $SortedBuild
    }
}

