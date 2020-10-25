#!/usr/bin/env bash
set -eux -o pipefail

# Set the launch context of the script
export CONTEXT_DIR=${CONTEXT_DIR:-"/vagrant"}

# Load the standard environment
source /home/vagrant/.hashicorprc

# Make scratch directory
work_dir=$(mktemp -d)

# Make install target for build requirements
mkdir -p ${work_dir}/gopath

# Make our install targets
job_dir=/etc/hashicorp/hcp.d/nomad
sudo mkdir -p ${job_dir}
sudo chown -R vagrant:vagrant ${job_dir}

# Download, extract, and install Go
cat << EOF > ${work_dir}/golang_checksum
85007dec7ca582e262dba97c24261e99ca387ed2500e86999d5170aad70d39abe27f270f61d00de4a6727b8009900e2bee20c8086a7dfeb5fe484b65758002a9  go1.15.3.linux-amd64.tar.gz
EOF

pushd ${work_dir}
  wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz

  sha512sum -c golang_checksum
  tar -C ${work_dir} -xzf go1.15.3.linux-amd64.tar.gz
popd

# Add build time requirements to PATH
export PATH=${work_dir}/go/bin:${PATH}

# Set Go environment variables
export GOPATH=${work_dir}/gopath
export GOROOT=${work_dir}/go

# Check the Go version
go version

# Build the service
pushd ${CONTEXT_DIR}/cmd/figlet-cadence-worker
  git rev-parse HEAD > ${work_dir}/git_sha
  go build

  # TODO no longer place this in such a global location after creating better chroots with Nomad
  sudo mv figlet-cadence-worker /usr/bin/figlet-cadence-worker
popd

# Render Nomad job, and deploy
git_sha=$(cat ${work_dir}/git_sha)
cat << FIGLET_NOMAD_JOB > ${job_dir}/figlet-service.nomad
job "figlet-service" {
  datacenters = ["dc1"]
  group "figlet-service" {
    volume "certs" {
      type      = "host"
      source    = "ca-certificates"
      read_only = true
    }

    task "figlet-cadence-worker" {
      driver = "exec"
      config {
        command = "/usr/bin/figlet-cadence-worker"
        args = []
      }
      volume_mount {
        volume      = "certs"
        destination = "/etc/ssl/certs"
      }

      resources {
        memory = 300
        network {
          mode = "host"
          port "api" { 
            static = 8091
            to = 8091
          }
        }
      }

      service {
        name = "figelt-service"
        port = "api"
      }

      env {
        GIT_SHA = "${git_sha}"
      }
    }
  }
}
FIGLET_NOMAD_JOB

nomad job run ${job_dir}/figlet-service.nomad

# Check deployment
sleep 10s
nomad status figlet-service

curl --connect-timeout 5 --max-time 10 --retry 10 --retry-max-time 60 --retry-connrefused http://127.0.0.1:8091/healthz

# Cleanup
rm -rf ${work_dir}