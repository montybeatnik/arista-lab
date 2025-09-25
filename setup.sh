#!/usr/bin/env bash
set -euo pipefail

### ───────────────────────────
### EDIT THESE FOR YOUR SETUP
### ───────────────────────────
VM_NAME="clab-evpn"
UBUNTU_RELEASE="jammy"          # or "22.04"
CPUS="6"
MEM="12G"
DISK="60G"

# Your existing repo (must contain lab.clab.yml and configs/)
REPO_DIR="/Users/chrishern/src/github.com/montybeatnik/arista-lab"

# Path to your **ARM64** cEOS-lab tarball on your Mac
CEOS_TARBALL="$REPO_DIR/cEOSarm-lab-4.34.2.1F.tar.xz"

# How you want the Docker image tagged inside the VM
CEOS_TAG="ceosimage:4.34.2.1f"

# Set to 1 to auto-deploy the lab at the end, 0 to skip
AUTO_DEPLOY=1
LAB_FILE="lab.clab.yml"         # path inside your REPO_DIR

### ───────────────────────────
### sanity checks
### ───────────────────────────
err() { echo "ERROR: $*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || err "Missing dependency '$1' on host."; }

need multipass

[[ -d "$REPO_DIR" ]] || err "REPO_DIR not found: $REPO_DIR"
[[ -f "$REPO_DIR/$LAB_FILE" ]] || err "Lab file not found at: $REPO_DIR/$LAB_FILE"
[[ -f "$CEOS_TARBALL" ]] || err "cEOS tarball not found: $CEOS_TARBALL"

echo ">>> Using repo:       $REPO_DIR"
echo ">>> Using cEOS image: $CEOS_TARBALL  ->  $CEOS_TAG"

### ───────────────────────────
### launch (or reuse) multipass VM
### ───────────────────────────
if multipass info "$VM_NAME" >/dev/null 2>&1; then
  echo ">>> VM '$VM_NAME' already exists. Reusing."
else
  echo ">>> Launching VM '$VM_NAME' ($UBUNTU_RELEASE, $CPUS vCPU, $MEM RAM, $DISK disk)..."
  multipass launch "$UBUNTU_RELEASE" --name "$VM_NAME" --cpus "$CPUS" --memory "$MEM" --disk "$DISK"
fi

echo ">>> Waiting for cloud-init to finish..."
multipass exec "$VM_NAME" -- cloud-init status --wait

### ───────────────────────────
### mount repo & transfer cEOS tarball
### ───────────────────────────
# Mount your repo into the VM at ~/lab
if multipass mount --help | grep -q "mount"; then
  echo ">>> Mounting repo -> $VM_NAME:~/lab"
  multipass mount "$REPO_DIR" "$VM_NAME":/home/ubuntu/lab
else
  err "Your multipass doesn't support 'mount'. Update Multipass or copy files manually."
fi

# Ensure images dir exists
multipass exec "$VM_NAME" -- bash -lc "mkdir -p ~/images"

# Transfer the cEOS tarball (multipass mount expects directories; use transfer for a file)
echo ">>> Transferring cEOS tarball into VM ~/images/"
multipass transfer "$CEOS_TARBALL" "$VM_NAME":/home/ubuntu/images/

### ───────────────────────────
### install docker, containerlab, and tools
### ───────────────────────────
echo ">>> Installing Docker, Containerlab, and utilities inside the VM..."
multipass exec "$VM_NAME" -- bash -lc '
  set -euo pipefail
  sudo apt-get update
  sudo apt-get install -y ca-certificates curl gnupg lsb-release jq iproute2 arping iperf3

  # Docker (convenience script)
  if ! command -v docker >/dev/null 2>&1; then
    curl -fsSL https://get.docker.com | sh
    sudo usermod -aG docker ubuntu
  fi

  # Containerlab
  if ! command -v containerlab >/dev/null 2>&1; then
    curl -sL https://get.containerlab.dev | bash
  fi

  # Helpful sysctls (not strictly required for clab, but useful)
  echo "net.ipv4.ip_forward=1" | sudo tee /etc/sysctl.d/99-clab.conf >/dev/null
  sudo sysctl --system >/dev/null
'

### ───────────────────────────
### import cEOS image, sanity check arch
### ───────────────────────────
echo ">>> Importing cEOS Docker image as $CEOS_TAG ..."
multipass exec "$VM_NAME" -- bash -lc "
  set -euo pipefail
  ls -lh ~/images
  # Import (idempotent-ish; will re-create if the tag doesn't exist)
  if ! docker image inspect '$CEOS_TAG' >/dev/null 2>&1; then
    docker import ~/images/$(basename "$CEOS_TARBALL") '$CEOS_TAG'
  else
    echo 'Docker image $CEOS_TAG already present, skipping import.'
  fi

  echo '>>> Verifying image arch...'
  IMG_ARCH=\$(docker inspect '$CEOS_TAG' --format '{{.Architecture}}')
  HOST_ARCH=\$(uname -m)
  echo \"Image arch: \$IMG_ARCH | Host arch: \$HOST_ARCH\"
  # Acceptable mappings: arm64 <-> aarch64
  if [[ \"\$IMG_ARCH\" != \"arm64\" && \"\$IMG_ARCH\" != \"aarch64\" ]]; then
    echo 'ERROR: Imported image is not ARM64. Please import an ARM64 cEOS tarball.' >&2
    exit 1
  fi
  if [[ \"\$HOST_ARCH\" != \"aarch64\" && \"\$HOST_ARCH\" != \"arm64\" ]]; then
    echo 'WARNING: VM is not arm64/aarch64; ensure image arch matches host.'
  fi
"

### ───────────────────────────
### print versions & where to go
### ───────────────────────────
multipass exec "$VM_NAME" -- bash -lc '
  echo ">>> Versions:"
  docker --version
  containerlab version || true
  echo
  echo ">>> Repo is mounted at: ~/lab"
  ls -la ~/lab
'

# Fix permissions issue
multipass exec "$VM_NAME" -- bash -lc '
  mkdir -p ~/.clab-runs
  if ! grep -q CLAB_LABDIR_BASE ~/.bashrc; then
    echo "export CLAB_LABDIR_BASE=\$HOME/.clab-runs" >> ~/.bashrc
  fi
'

### ───────────────────────────
### optional: deploy the lab
### ───────────────────────────
if [[ "$AUTO_DEPLOY" -eq 1 ]]; then
  echo ">>> Deploying lab from ~/lab/$LAB_FILE ..."
  multipass exec "$VM_NAME" -- bash -lc "
    cd ~/lab
    sudo -E containerlab deploy -t '$LAB_FILE' --reconfigure
    echo
    echo '>>> Lab status:'
    containerlab inspect -t '$LAB_FILE' | sed -n '1,120p'
  "
else
  cat <<EOF

>>> Skipping auto-deploy (AUTO_DEPLOY=0).
To deploy manually:

  multipass shell $VM_NAME
  cd ~/lab
  sudo containerlab deploy -t $LAB_FILE --reconfigure

EOF
fi

echo ">>> Done. SSH into the VM with: multipass shell $VM_NAME"
