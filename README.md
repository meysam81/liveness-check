# liveness-check

<div align="center">

<!-- Project Status & Quality -->

[![CI/CD](https://github.com/meysam81/liveness-check/actions/workflows/ci.yml/badge.svg)](https://github.com/meysam81/liveness-check/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/meysam81/liveness-check)](https://goreportcard.com/report/github.com/meysam81/liveness-check)
[![Vulnerability Scan](https://img.shields.io/badge/🛡️_Zero_Vulnerabilities-Kubescape_Verified-brightgreen?style=flat-square)](https://github.com/meysam81/liveness-check/actions)

<!-- Release & Distribution -->

[![Latest Release](https://img.shields.io/github/v/release/meysam81/liveness-check?style=flat-square&logo=github&color=blue)](https://github.com/meysam81/liveness-check/releases/latest)
[![Docker Image](https://img.shields.io/badge/docker-meysam81%2Fliveness--check-blue?style=flat-square&logo=docker)](https://hub.docker.com/r/meysam81/liveness-check)
[![Docker Pulls](https://img.shields.io/docker/pulls/meysam81/liveness-check?style=flat-square&logo=docker)](https://hub.docker.com/r/meysam81/liveness-check)
[![Go Version](https://img.shields.io/github/go-mod/go-version/meysam81/liveness-check?style=flat-square&logo=go)](go.mod)

<!-- License & Community -->

[![License](https://img.shields.io/badge/License-Apache--2.0-green.svg?style=flat-square)](LICENSE)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/meysam81?style=flat-square&logo=github&color=pink)](https://github.com/sponsors/meysam81)

<!-- Technical Features -->

[![Single Binary](https://img.shields.io/badge/🚀_Single-Binary-blueviolet?style=flat-square)](https://golang.org/)
[![Cross Platform](https://img.shields.io/badge/🌐_Cross-Platform-orange?style=flat-square)](https://golang.org/)
[![Container Native](https://img.shields.io/badge/📦_Container-Native-2496ED?style=flat-square&logo=docker)](https://kubernetes.io/)

<!-- DevOps & Monitoring Features -->

[![Health Checks](https://img.shields.io/badge/💓_Health-Checks-FF6B6B?style=flat-square)](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)
[![Preview Envs](https://img.shields.io/badge/🔍_Preview-Environments-9C27B0?style=flat-square)](https://kubernetes.io/)
[![Production Ready](https://img.shields.io/badge/🏭_Production-Ready-darkgreen?style=flat-square)](https://sre.google/)
[![Air Gap Ready](https://img.shields.io/badge/🔒_Air--Gap-Compatible-darkred?style=flat-square)](#install)

<!-- Hackery & Performance -->

[![Jitter Logic](https://img.shields.io/badge/🎯_Smart-Jitter-purple?style=flat-square)](#usage)
[![Infinite Retries](https://img.shields.io/badge/♾️_Infinite-Retries-teal?style=flat-square)](#options)
[![Millisecond Precision](https://img.shields.io/badge/⏱️_ms-Precision-indigo?style=flat-square)](https://golang.org/pkg/time/)

</div>

A lightweight CLI tool for HTTP health checks with configurable retries and timeout.

This is best used in containerized environments and more specifically, to check
the healthcheck of a recently deployed pod for preview environment.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Install](#install)
- [Usage](#usage)
  - [Environment Variables](#environment-variables)
  - [Docker](#docker)
  - [Kubernetes Example](#kubernetes-example)
- [Options](#options)
- [Build](#build)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Install

```bash
go install github.com/meysam81/liveness-check@latest
```

Or use Docker:

```bash
docker pull ghcr.io/meysam81/liveness-check:latest
```

## Usage

```bash
# Basic health check
liveness-check check --http-target http://localhost:8080/health

# With custom timeout and retries
liveness-check check \
  --http-target http://api.example.com/health \
  --timeout 10 \
  --retries 5 \
  --status-code 200
```

### Environment Variables

```bash
export HTTP_TARGET=http://localhost:8080/health
export TIMEOUT=10
export RETRIES=3
export STATUS_CODE=200
export LOG_LEVEL=info

liveness-check check
```

### Docker

```bash
docker run --rm ghcr.io/meysam81/liveness-check:latest \
  check --http-target http://host.docker.internal:8080/health
```

### Kubernetes Example

```yaml
---
apiVersion: batch/v1
kind: Job
metadata:
  name: liveness-check
spec:
  template:
    spec:
      containers:
        - args:
            - echo
            - all good
          image: busybox
          name: busybox
      initContainers:
        - args:
            - check
            - "--http-target"
            - http://my-service.default.svc.cluster.local/health
          image: ghcr.io/meysam81/liveness-check
          name: liveness-check
          resources:
            limits:
              cpu: 10m
              memory: 10Mi
            requests:
              cpu: 10m
              memory: 10Mi
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
            readOnlyRootFilesystem: true
            runAsGroup: 65534
            runAsNonRoot: true
            runAsUser: 65534
          terminationMessagePolicy: FallbackToLogsOnError
      restartPolicy: OnFailure
```

## Options

- `--http-target, -u`: Target URL to check (required)
- `--timeout, -t`: Request timeout in seconds (default: 5)
- `--retries, -r`: Number of retries, 0 for infinite (default: 0)
- `--status-code, -c`: Expected HTTP status code (default: 200)
- `--log-level, -l`: Log verbosity: debug, info, warn, error, critical (default: info)

## Build

```bash
go build -o liveness-check
```
