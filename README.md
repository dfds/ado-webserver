# :warning: Repository not maintained :warning:

Please note that this repository is currently archived, and is no longer being maintained.

- It may contain code, or reference dependencies, with known vulnerabilities
- It may contain out-dated advice, how-to's or other forms of documentation

The contents might still serve as a source of inspiration, but please review any contents before reusing elsewhere.


<!-- omit in toc -->
# PowerShell Webserver

[![Contributors][contributors-shield]][contributors-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![Build Status](https://dev.azure.com/dfds/YourAzureDevOpsProject/_apis/build/status/Name-Of-CI-Pipeline?branchName=master)](https://dev.azure.com/dfds/YourAzureDevOpsProject/_build/latest?definitionId=1378&branchName=master)

<!-- TABLE OF CONTENTS -->
<!-- omit in toc -->
## Table of Contents

- [About The Project](#about-the-project)
  - [Structure](#structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
- [Usage](#usage)
- [Deployment Flow](#deployment-flow)
  - [Infrastructure Deployment](#infrastructure-deployment)
  - [Collector Deployment](#collector-deployment)
- [Additional Resources](#additional-resources)
- [License](#license)

<!-- ABOUT THE PROJECT -->
## About The Project

Webserver that does stuff

### Structure

Notable project directories and files:

| Path                   | Usage                                                                 |
| ---------------------- | --------------------------------------------------------------------- |
| `/Dockerfile`          | Dockerfile used for building images with the `eks-{sourceref}` tag    |
| `/src/`                | PowerShell scripts that make up the webserver                         |
| `/k8s/`                | Manifest for deploying image to Kubernetes                            |
| `azure-pipelines.yaml` | Pipeline spec to build and deploy container image to kubernetes       |

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- [Powershell Core][powershell-core] (tested with 7.x)
- [AWS Tools for Powershell NetCore][aws-powershell] 4.0.0+

<!-- USAGE EXAMPLES -->
## Usage

The script usage is documented natively in PowerShell. To get the most up-to-date help, invoke:


The priority is to update the built-in help when relevant. The following table of the script arguments and their description is provided for convenience, but might not always we up-to-date:

| Argument            | Description                                                                                                                          |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| AwsProfile          | Name of the AWS profile to use for authentication. If not specified, the normal credential search order is used.                     |
| AwsRegion           | AWS region where the CloudWatch Logs and target S3 bucket reside, e.g. 'eu-west-1'.                                                  |


See also [Additional Resources](#additional-resources).

## Deployment Flow

### Webserver Deployment

1. Change code in `./src/` folder
2. Update in code triggers CI/CD Pipeline

## Additional Resources

- [AWS Tools for Powershell Docs: Using AWS Credentials][aws-docs-credentials]

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/dfds/cloudwatchlogs-collector?style=plastic
[contributors-url]: https://github.com/dfds/cloudwatchlogs-collector/graphs/contributors
[issues-shield]: https://img.shields.io/github/issues/dfds/cloudwatchlogs-collector?style=plastic
[issues-url]: https://github.com/dfds/cloudwatchlogs-collector/issues
[license-shield]: https://img.shields.io/github/license/dfds/cloudwachlogs-collector?style=plastic
[license-url]: https://github.com/dfds/cloudwatchlogs-collector/blob/master/LICENSE
[powershell-core]: https://github.com/PowerShell/PowerShell/releases
[aws-powershell]: https://docs.aws.amazon.com/powershell/latest/userguide/pstools-getting-set-up.html
[docker-tags]: https://hub.docker.com/r/dfdsdk/cloudwatchlogs-collector/tags
[aws-docs-credentials]: https://docs.aws.amazon.com/powershell/latest/userguide/specifying-your-aws-credentials.html
